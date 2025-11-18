// 标签领域模型：用于标注文章主题与属性
package tag

// Tag 表示一个文章标签，通常与文章多对多关联
type Tag struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex;size:64"` // 标签名称
    Slug string `gorm:"uniqueIndex;size:64"` // URL 友好短标识
}