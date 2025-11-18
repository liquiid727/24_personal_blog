// File: internal/service/file_service.go
// Purpose: File domain service to persist uploads and manage records.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Saves file metadata, supports retrieval and deletion with filesystem cleanup.
package service

import (
	"os"
	"path/filepath"
	"time"

	f "go_blog/internal/domain/file"
	"go_blog/internal/repository"
)

type FileService interface {
	Save(uploaderID uint, path, mime string, size int64) (*f.File, error)
	Get(id uint) (*f.File, error)
	Delete(id uint) error
}

type fileService struct{ repo repository.FileRepository }

// NewFileService 创建文件服务
// Params: r 文件仓储实现
// Returns: FileService
func NewFileService(r repository.FileRepository) FileService { return &fileService{repo: r} }

// Save 保存上传文件记录
// Params: uploaderID 上传者ID；path 文件路径；mime MIME类型；size 文件大小
// Behavior: 写入记录并返回
func (s *fileService) Save(uploaderID uint, path, mime string, size int64) (*f.File, error) {
	fi := &f.File{UploaderID: uploaderID, Path: path, Mime: mime, Size: size, CreatedAt: time.Now()}
	if err := s.repo.Create(fi); err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *fileService) Get(id uint) (*f.File, error) { return s.repo.FindByID(id) }

// Delete 删除文件记录
// Behavior: 先尝试删除物理文件，再删除数据库记录；忽略文件不存在错误
func (s *fileService) Delete(id uint) error {
	fi, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	_ = os.Remove(filepath.Clean(fi.Path))
	return s.repo.Delete(id)
}
