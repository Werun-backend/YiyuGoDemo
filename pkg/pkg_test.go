package pkg

import (
	"awesomeProject1/pkg/models"
	"awesomeProject1/response"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	diary := models.Diary{
		Model: gorm.Model{
			ID:        1222,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		UserID:  1,
		Content: "123123123",
		Tags:    nil,
	}

	var diaryVo response.DiaryVo

	err := SimpleCopyProperties(&diaryVo, diary)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("\n"+
		"diary: %v\ndiaryVo: %v\n", diary, diaryVo)
}
