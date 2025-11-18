// 分类领域模型：用于对文章进行分组分类
package category

// Category 表示一个文章分类，包含名称与唯一的短标识
type Category struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex;size:64"` // 分类名称
    Slug string `gorm:"uniqueIndex;size:64"` // URL 友好短标识
}