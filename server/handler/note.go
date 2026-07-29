package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notepad/model"
)

type createNoteRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

type updateNoteRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

type renameTagRequest struct {
	OldName string `json:"old_name" binding:"required,max=80"`
	NewName string `json:"new_name" binding:"required,max=80"`
}

type deleteTagRequest struct {
	Name string `json:"name" binding:"required,max=80"`
}

func ListNotes(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("q")
	tag := c.Query("tag")
	sortBy := c.DefaultQuery("sort_by", "updated_at")

	result, err := model.ListNotes(userID.(int64), page, pageSize, search, tag, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func CreateNote(c *gin.Context) {
	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	userID, _ := c.Get("userID")
	note, err := model.CreateNote(userID.(int64), req.Title, req.Content, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func GetNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	note, err := model.GetNote(id, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func UpdateNote(c *gin.Context) {
	var req updateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	userID, _ := c.Get("userID")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	if err := model.UpdateNote(id, userID.(int64), req.Title, req.Content, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	if err := model.DeleteNote(id, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllTags(c *gin.Context) {
	userID, _ := c.Get("userID")
	stats, err := model.GetTagStats(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}
	if stats == nil {
		stats = []model.TagStat{}
	}

	tags := make([]string, 0, len(stats))
	for _, stat := range stats {
		tags = append(tags, stat.Name)
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "items": stats})
}

func RenameTag(c *gin.Context) {
	var req renameTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供原标签名和新标签名"})
		return
	}

	userID, _ := c.Get("userID")
	affected, err := model.RenameTag(userID.(int64), req.OldName, req.NewName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "标签已重命名", "affected_notes": affected})
}

func DeleteTag(c *gin.Context) {
	var req deleteTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供标签名"})
		return
	}

	userID, _ := c.Get("userID")
	affected, err := model.DeleteTag(userID.(int64), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "标签已删除", "affected_notes": affected})
}
