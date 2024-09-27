package models

import (
	"gorm.io/gorm"
)

// Diary 代表日记的数据库模型
type Diary struct {
	gorm.Model
	UserID  uint   ` gorm:"not null"`            // 用户ID，外键
	Content string `gorm:"type:text"`            // 日记内容
	User    User   `gorm:"foreignKey:ID"`        // 关联用户
	Tags    []Tag  `gorm:"many2many:diary_tags"` // 多对多关联标签
}
