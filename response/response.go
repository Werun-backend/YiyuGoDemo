package response

import (
	"github.com/gin-gonic/gin"
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

// TODO: 统一返回类组件

// SuccessResponse 创建一个表示成功的响应。
func SuccessResponse(data interface{}) ApiResponse {
	return NewResponse(200, data, "success")
}

// UserErrorResponse 创建一个表示用户错误的响应。
func UserErrorResponse(data interface{}, message string) ApiResponse {
	return NewResponse(400, data, message)
}

// UserErrorNoMsgResponse 创建一个表示用户错误的响应，不包含额外的数据。
func UserErrorNoMsgResponse(message string) ApiResponse {
	return NewResponse(400, nil, message)
}

// InternalErrorResponse 创建一个表示服务器内部错误的响应。
func InternalErrorResponse(data interface{}, message string) ApiResponse {
	return NewResponse(500, data, message)
}

// WriteJSON 将 ApiResponse 写入 HTTP 响应
func WriteJSON(c *gin.Context, response ApiResponse) {
	c.JSON(response.Code, response)
}
