package server

import (
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"

    h "go_blog/internal/transport/http"
    "go_blog/internal/service"
)

type FileHandler struct{ svc service.FileService }

func NewFileHandler(s service.FileService) *FileHandler { return &FileHandler{svc: s} }

func (fh *FileHandler) Upload(c *gin.Context) {
    f, err := c.FormFile("file")
    if err != nil { h.Error(c, http.StatusBadRequest, 1, "file required"); return }
    if f.Size <= 0 { h.Error(c, http.StatusBadRequest, 2, "invalid size"); return }
    // 仅允许常见图片类型，可按需扩展
    mime := f.Header.Get("Content-Type")
    if !strings.HasPrefix(mime, "image/") { h.Error(c, http.StatusBadRequest, 3, "only image allowed"); return }
    _ = os.MkdirAll("uploads", 0755)
    filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(f.Filename)
    path := filepath.Join("uploads", filename)
    if err := c.SaveUploadedFile(f, path); err != nil { h.Error(c, http.StatusInternalServerError, 4, "save failed"); return }
    uid := c.GetUint("uid")
    rec, err := fh.svc.Save(uid, path, mime, f.Size)
    if err != nil { h.Error(c, http.StatusInternalServerError, 5, "persist failed"); return }
    h.OK(c, gin.H{"id": rec.ID, "url": "/uploads/" + filename, "mime": rec.Mime, "size": rec.Size})
}

func (fh *FileHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    // 基于 path 参数转换
    var nid uint
    if _, err := fmt.Sscanf(id, "%d", &nid); err != nil { h.Error(c, http.StatusBadRequest, 1, "invalid id"); return }
    // 仅管理员或本人可删（此处简单示例：管理员）
    if c.GetString("role") != "admin" { h.Error(c, http.StatusForbidden, 2, "forbidden"); return }
    if err := fh.svc.Delete(nid); err != nil { h.Error(c, http.StatusInternalServerError, 3, "delete failed"); return }
    h.OK(c, gin.H{"id": nid})
}

func (fh *FileHandler) Get(c *gin.Context) {
    id := c.Param("id")
    var nid uint
    if _, err := fmt.Sscanf(id, "%d", &nid); err != nil { h.Error(c, http.StatusBadRequest, 1, "invalid id"); return }
    rec, err := fh.svc.Get(nid)
    if err != nil { h.Error(c, http.StatusNotFound, 2, "not found"); return }
    h.OK(c, rec)
}