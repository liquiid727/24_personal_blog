// 评论领域模型：支持文章评论与多层回复
package comment

import "time"

// Comment 表示一条评论，支持通过 ParentID 形成树结构
type Comment struct {
    ID        uint      `gorm:"primaryKey"`
    PostID    uint      `gorm:"index"`      // 所属文章 ID
    UserID    uint      `gorm:"index"`      // 发表用户 ID
    ParentID  *uint     `gorm:"index"`      // 父评论 ID（为空表示顶层）
    Content   string    `gorm:"type:text"`  // 评论内容
    Status    string    `gorm:"size:16"`    // 审核状态（pending/approved/rejected）
    CreatedAt time.Time
    UpdatedAt time.Time
}