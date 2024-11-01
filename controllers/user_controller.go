package controllers

import (
	"awesomeProject1/pkg/jwt"
	"awesomeProject1/pkg/models"
	"awesomeProject1/request"
	"awesomeProject1/response"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

// UserController 用户控制器
type UserController struct {
	DB *gorm.DB
}

// Routes 确保在 UserController 的 Routes 方法中注册 Login 路由
func (ctrl *UserController) Routes(r *gin.Engine) {
	r.POST("/diary/login", ctrl.Login)
}

// NewUserController 创建 UserController 的新实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Login 用户注册或登录
func (ctrl *UserController) Login(c *gin.Context) {
	var loginDto request.LoginDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		response.WriteJSON(c, response.UserErrorResponse(c, "参数绑定失败"))
		return
	}

	// 检查用户是否存在
	var user models.User
	result := ctrl.DB.Where("phone_number = ?", loginDto.PhoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 用户不存在，执行注册逻辑
			hashPassword, err := bcrypt.GenerateFromPassword([]byte(loginDto.Password), bcrypt.DefaultCost)
			if err != nil {
				response.WriteJSON(c, response.UserErrorResponse(c, "密码加密失败"))
				return
			}
			user.PhoneNumber = loginDto.PhoneNumber
			user.Password = string(hashPassword)
			// 保存新用户
			createResult := ctrl.DB.Create(&user)
			if createResult.Error != nil {

				response.WriteJSON(c, response.UserErrorResponse(c, "用户注册失败"))
				return
			}
		} else {

			response.WriteJSON(c, response.UserErrorResponse(c, "数据库查询失败"))
			return
		}
	} else {
		// 用户存在，执行登录逻辑
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
		if err != nil {
			response.WriteJSON(c, response.UserErrorResponse(c, "密码错误"))
			return
		}
	}

	// 登录成功，生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.PhoneNumber)
	if err != nil {
		response.WriteJSON(c, response.UserErrorResponse(c, "令牌生成失败"))
	} else {
		// 登录成功，返回统一的响应格式
		loginVo := response.LoginVo{
			UserId:      strconv.Itoa(int(user.ID)),
			PhoneNumber: user.PhoneNumber,
			Token:       token,
		}
		response.WriteJSON(c, response.SuccessResponse(loginVo))
	}

}
