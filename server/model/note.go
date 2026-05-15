package model

import (
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
	rows, err := database.DB.Query("SELECT DISTINCT tags FROM notes WHERE user_id = ? AND tags != ''", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsStr string
		if err := rows.Scan(&tagsStr); err != nil {
			return nil, err
		}
		for _, tag := range splitTags(tagsStr) {
			if tag != "" {
				tagSet[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags, nil
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
