package request

import "awesomeProject1/pkg/models"

// DiaryDto 代表日记的请求模型
// Dto Data Transform Object
// TODO: Dto 以及数据转换
// DiaryDto 代表日记的请求模型
type DiaryDto struct {
	Content string `json:"content"` // 日记内容
	TagIds  []uint `json:"tagIds"`
}

// TagDto 代表标签的请求模型
type TagDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UserDto 代表用户的请求模型
type UserDto struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
}

// DiaryToDto 将Diary模型转换为DiaryDto
func DiaryToDto(diary models.Diary) DiaryDto {
	var tagIds []uint
	for _, tag := range diary.Tags {
		tagIds = append(tagIds, tag.ID)
	}
	return DiaryDto{
		Content: diary.Content,
		TagIds:  tagIds,
	}
}

// TagToDto 将Tag模型转换为TagDto
func TagToDto(tag models.Tag) TagDto {
	return TagDto{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

// UserToDto 将User模型转换为UserDto
func UserToDto(user models.User) UserDto {
	return UserDto{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
	}
}
