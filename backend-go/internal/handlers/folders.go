package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/instaagrammeta/crm-real/backend-go/internal/uploads"
	"gorm.io/gorm"
)

type FoldersHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

// ====================== Folders ======================

func (h *FoldersHandler) List(c *gin.Context) {
	parent := QueryInt(c, "parent_id", 0)
	var rows []models.Folder
	if err := h.DB.Where("parent_id = ?", parent).Order("name ASC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *FoldersHandler) Create(c *gin.Context) {
	var f models.Folder
	if err := c.ShouldBindJSON(&f); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	f.AuthorID = middleware.CurrentUserID(c)
	f.AuthorName = middleware.CurrentUserName(c)
	f.CreatedAt = time.Now()
	if err := h.DB.Create(&f).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": f.ID, "folder": f})
}

func (h *FoldersHandler) Rename(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.Folder{}).Where("id = ?", id).
		Update("name", body.Name).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *FoldersHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	// recursive delete: collect descendants
	allIDs := []uint{id}
	current := []uint{id}
	for len(current) > 0 {
		var children []models.Folder
		tx.Where("parent_id IN ?", current).Find(&children)
		current = current[:0]
		for _, ch := range children {
			allIDs = append(allIDs, ch.ID)
			current = append(current, ch.ID)
		}
	}
	tx.Where("folder_id IN ?", allIDs).Delete(&models.FolderFile{})
	if err := tx.Where("id IN ?", allIDs).Delete(&models.Folder{}).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ====================== Files ======================

func (h *FoldersHandler) ListFiles(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var rows []models.FolderFile
	if err := h.DB.Where("folder_id = ?", id).Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *FoldersHandler) UploadFile(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		JSONError(c, http.StatusBadRequest, "no file")
		return
	}
	path, err := uploads.SaveUpload(h.Cfg.UploadDir, "files", file)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	rec := models.FolderFile{
		FolderID:     id,
		Filename:     filepath.Base(path),
		OriginalName: file.Filename,
		Filepath:     path,
		Filetype:     file.Header.Get("Content-Type"),
		Filesize:     file.Size,
		AuthorID:     middleware.CurrentUserID(c),
		AuthorName:   middleware.CurrentUserName(c),
		CreatedAt:    time.Now(),
	}
	if err := h.DB.Create(&rec).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "file": rec})
}

func (h *FoldersHandler) DeleteFile(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var f models.FolderFile
	if err := h.DB.First(&f, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "file not found")
		return
	}
	// best-effort удаление физического файла
	if f.Filepath != "" {
		full := filepath.Join(h.Cfg.UploadDir, filepath.FromSlash(stripUploadsPrefix(f.Filepath)))
		_ = os.Remove(full)
	}
	if err := h.DB.Delete(&f).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Download — GET /api/download/:id
func (h *FoldersHandler) Download(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var f models.FolderFile
	if err := h.DB.First(&f, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "file not found")
		return
	}
	full := filepath.Join(h.Cfg.UploadDir, filepath.FromSlash(stripUploadsPrefix(f.Filepath)))
	c.FileAttachment(full, f.OriginalName)
}

func stripUploadsPrefix(s string) string {
	const prefix = "/uploads/"
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
