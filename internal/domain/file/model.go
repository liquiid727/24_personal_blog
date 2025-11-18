// 文件领域模型：用于管理上传文件的元数据
package file

import "time"

// File 表示一次文件上传记录
type File struct {
    ID         uint      `gorm:"primaryKey"`
    UploaderID uint      `gorm:"index"`     // 上传者用户 ID
    Path       string    `gorm:"size:512"`  // 存储路径或 URL
    Mime       string    `gorm:"size:128"`  // MIME 类型
    Size       int64                          // 文件大小（字节）
    CreatedAt  time.Time
}