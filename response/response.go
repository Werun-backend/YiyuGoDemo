// response.go
package response

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
}

// NewResponse 创建一个新的 ApiResponse 实例
func NewResponse(code int, data interface{}, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

// WriteJSON 将 ApiResponse 写入 HTTP 响应
func WriteJSON(c *gin.Context, response ApiResponse) {
	c.JSON(http.StatusOK, response)
}

// DiaryResponse 用于API响应的日记信息
type DiaryResponse struct {
	ID        uint           `json:"id"`
	CreatedAt string         `json:"createdAt,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	UserID    uint           `json:"userID"`
	Content   string         `json:"content"`
	//User      UserResponse   `json:"user,omitempty"`
	Tags []TagResponse `json:"tags,omitempty"`
}

// UserResponse 用于API响应的用户信息
type UserResponse struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Password    string `json:"password,omitempty"`
}

// TagResponse 用于API响应的标签信息
type TagResponse struct {
	ID        uint           `json:"id"`
	CreatedAt string         `json:"createdAt,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name      string         `json:"name"`
}
