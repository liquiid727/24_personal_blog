package repository

import (
    "gorm.io/gorm"
    f "go_blog/internal/domain/file"
)

type FileRepository interface {
    Create(file *f.File) error
    FindByID(id uint) (*f.File, error)
    Delete(id uint) error
}

type fileRepository struct{ db *gorm.DB }

func NewFileRepository(db *gorm.DB) FileRepository { return &fileRepository{db: db} }

func (r *fileRepository) Create(file *f.File) error { return r.db.Create(file).Error }

func (r *fileRepository) FindByID(id uint) (*f.File, error) {
    var res f.File
    if err := r.db.First(&res, id).Error; err != nil { return nil, err }
    return &res, nil
}

func (r *fileRepository) Delete(id uint) error { return r.db.Delete(&f.File{}, id).Error }