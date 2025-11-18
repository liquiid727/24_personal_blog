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

func NewFileService(r repository.FileRepository) FileService { return &fileService{repo: r} }

func (s *fileService) Save(uploaderID uint, path, mime string, size int64) (*f.File, error) {
    fi := &f.File{UploaderID: uploaderID, Path: path, Mime: mime, Size: size, CreatedAt: time.Now()}
    if err := s.repo.Create(fi); err != nil { return nil, err }
    return fi, nil
}

func (s *fileService) Get(id uint) (*f.File, error) { return s.repo.FindByID(id) }

func (s *fileService) Delete(id uint) error {
    fi, err := s.repo.FindByID(id)
    if err != nil { return err }
    _ = os.Remove(filepath.Clean(fi.Path))
    return s.repo.Delete(id)
}