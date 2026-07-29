package model

import (
	"errors"
	"sort"
	"strings"
	"time"

	"notepad/database"
)

type Note struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteListResponse struct {
	Notes []Note `json:"notes"`
	Total int    `json:"total"`
}

type TagStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func CreateNote(userID int64, title, content, tags string) (*Note, error) {
	result, err := database.DB.Exec(
		"INSERT INTO notes (user_id, title, content, tags) VALUES (?, ?, ?, ?)",
		userID, title, content, tags,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &Note{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func GetNote(id, userID int64) (*Note, error) {
	note := &Note{}
	err := database.DB.QueryRow(
		"SELECT id, user_id, title, content, tags, created_at, updated_at FROM notes WHERE id = ? AND user_id = ?",
		id, userID,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.Tags, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func ListNotes(userID int64, page, pageSize int, search, tag, sortBy string) (*NoteListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int
	var err error
	var conditions string
	var args []interface{}

	conditions = "WHERE user_id = ?"
	args = append(args, userID)

	if search != "" {
		conditions += " AND (title LIKE ? OR content LIKE ?)"
		likeSearch := "%" + search + "%"
		args = append(args, likeSearch, likeSearch)
	}
	if tag != "" {
		conditions += " AND tags LIKE ?"
		args = append(args, "%"+tag+"%")
	}

	err = database.DB.QueryRow("SELECT COUNT(*) FROM notes "+conditions, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	queryArgs := append([]interface{}{}, args...)

	orderClause := "ORDER BY updated_at DESC"
	if sortBy == "title" {
		orderClause = "ORDER BY title ASC"
	} else if sortBy == "created_at" {
		orderClause = "ORDER BY created_at DESC"
	}

	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := database.DB.Query(
		"SELECT id, user_id, title, content, tags, created_at, updated_at FROM notes "+conditions+" "+orderClause+" LIMIT ? OFFSET ?",
		queryArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if notes == nil {
		notes = []Note{}
	}

	return &NoteListResponse{Notes: notes, Total: total}, nil
}

func UpdateNote(id, userID int64, title, content, tags string) error {
	_, err := database.DB.Exec(
		"UPDATE notes SET title = ?, content = ?, tags = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?",
		title, content, tags, id, userID,
	)
	return err
}

func DeleteNote(id, userID int64) error {
	_, err := database.DB.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func GetAllTags(userID int64) ([]string, error) {
	stats, err := GetTagStats(userID)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(stats))
	for _, stat := range stats {
		tags = append(tags, stat.Name)
	}
	return tags, nil
}

func GetTagStats(userID int64) ([]TagStat, error) {
	rows, err := database.DB.Query("SELECT tags FROM notes WHERE user_id = ? AND tags != ''", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tagsStr string
		if err := rows.Scan(&tagsStr); err != nil {
			return nil, err
		}
		seen := make(map[string]bool)
		for _, tag := range splitTags(tagsStr) {
			if tag != "" && !seen[tag] {
				tagCounts[tag]++
				seen[tag] = true
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(tagCounts))
	for name := range tagCounts {
		names = append(names, name)
	}
	sort.Strings(names)

	stats := make([]TagStat, 0, len(names))
	for _, name := range names {
		stats = append(stats, TagStat{Name: name, Count: tagCounts[name]})
	}
	return stats, nil
}

func RenameTag(userID int64, oldName, newName string) (int, error) {
	oldName = strings.TrimSpace(oldName)
	newName = strings.TrimSpace(newName)
	if oldName == "" || newName == "" {
		return 0, errors.New("标签名不能为空")
	}
	if oldName == newName {
		return 0, nil
	}

	return rewriteTag(userID, oldName, func(tag string) (string, bool) {
		if tag == oldName {
			return newName, true
		}
		return tag, true
	})
}

func DeleteTag(userID int64, name string) (int, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("标签名不能为空")
	}

	return rewriteTag(userID, name, func(tag string) (string, bool) {
		return tag, tag != name
	})
}

func rewriteTag(userID int64, target string, transform func(string) (string, bool)) (int, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, tags FROM notes WHERE user_id = ? AND tags != ''", userID)
	if err != nil {
		return 0, err
	}

	type tagUpdate struct {
		id   int64
		tags string
	}
	var updates []tagUpdate

	for rows.Next() {
		var id int64
		var tagsStr string
		if err := rows.Scan(&id, &tagsStr); err != nil {
			rows.Close()
			return 0, err
		}

		tags := splitTags(tagsStr)
		found := false
		seen := make(map[string]bool)
		rewritten := make([]string, 0, len(tags))
		for _, tag := range tags {
			if tag == target {
				found = true
			}
			nextTag, keep := transform(tag)
			nextTag = strings.TrimSpace(nextTag)
			if !keep || nextTag == "" || seen[nextTag] {
				continue
			}
			seen[nextTag] = true
			rewritten = append(rewritten, nextTag)
		}

		if found {
			updates = append(updates, tagUpdate{id: id, tags: strings.Join(rewritten, ", ")})
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, err
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}

	for _, update := range updates {
		if _, err := tx.Exec(
			"UPDATE notes SET tags = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?",
			update.tags, update.id, userID,
		); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(updates), nil
}

func splitTags(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range strings.Split(s, ",") {
		t := strings.TrimSpace(part)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
