package response

import (
	"awesomeProject1/pkg"
	"awesomeProject1/pkg/models"
)

// DiaryVo 用于API响应的日记信息
// Vo View Object
// TODO: 出参Vo 以及数据转换
type DiaryVo struct {
	ID        uint    `json:"id"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	UserID    uint    `json:"userID"`
	Content   string  `json:"content"`
	Tags      []TagVo `json:"tags,omitempty"`
}

// Copy 赋值
func (diaryVo *DiaryVo) Copy(diary models.Diary) {
	diaryVo.ID = diary.ID
	diaryVo.CreatedAt = diary.CreatedAt.Format("2006-01-02 15:04:05")
	diaryVo.UpdatedAt = diary.UpdatedAt.Format("2006-01-02 15:04:05")
	diaryVo.UserID = diary.UserID
	diaryVo.Content = diary.Content
	diaryVo.Tags = make([]TagVo, len(diary.Tags))
	for i, tag := range diary.Tags {
		diaryVo.Tags[i] = TagVo{
			ID:        tag.ID,
			CreatedAt: tag.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:      tag.Name,
		}
	}
}

type PaginatedDiaryVo struct {
	Data           []DiaryVo `json:"data"`
	PaginationData `json:"pagination"`
}

type PaginationData struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int   `json:"totalPages"`
}

func (pdv *PaginatedDiaryVo) Copy(diaryVos []DiaryVo, pagination PaginationData) {
	pdv.Data = diaryVos
	pdv.PaginationData = pagination
}

// UserVo 用于API响应的用户信息
type UserVo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

func (userVo *UserVo) Copy(user models.User) error {
	err := pkg.SimpleCopyProperties(userVo, user)
	if err != nil {
		return err
	}
	return nil
}

// TagVo 用于API响应的标签信息
type TagVo struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Name      string `json:"name"`
}

// Copy 赋值
func (tagVo *TagVo) Copy(tag models.Tag) {
	tagVo.ID = tag.ID
	tagVo.CreatedAt = tag.CreatedAt.Format("2006-01-02 15:04:05")
	tagVo.UpdatedAt = tag.UpdatedAt.Format("2006-01-02 15:04:05")
	tagVo.Name = tag.Name
}
