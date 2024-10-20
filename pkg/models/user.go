package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	PhoneNumber string `json:"phoneNumber" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(inputPassword string) bool {
	// 这里应该是密码比较逻辑，示例中假设密码以明文存储
	return u.Password == inputPassword
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 在这里添加创建用户之前的逻辑
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	// 在这里添加创建用户之后的逻辑
	return nil
}
