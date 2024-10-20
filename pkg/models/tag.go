package models

import (
	"gorm.io/gorm"
)

// Tag 代表标签的数据库模型
type Tag struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"` // 标签名称，需要唯一
	// 多对多关系，不返回空值
	Diaries []Diary `json:"diaries" gorm:"many2many:diary_tags" json:"diaries,omitempty"` // 与其他日记的关联
}
