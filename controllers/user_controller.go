package controllers

import (
	"awesomeProject1/pkg/jwt"
	"awesomeProject1/pkg/models"
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

// NewUserController 创建 UserController 的新实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Routes 注册路由
func (ctrl *UserController) Routes(r *gin.Engine) {
	r.POST("/diary/login", ctrl.Login)
}

// Login 用户注册或登录
func (ctrl *UserController) Login(c *gin.Context) {
	var credentials models.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的输入"))
		return
	}

	// 检查用户是否存在
	var user models.User
	result := ctrl.DB.Where("phone_number = ?", credentials.PhoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 用户不存在，执行注册逻辑
			hashPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
			if err != nil {
				response.WriteJSON(c, response.NewResponse(2, nil, "密码加密失败"))
				return
			}
			user.PhoneNumber = credentials.PhoneNumber
			user.Password = string(hashPassword)
			// 保存新用户
			createResult := ctrl.DB.Create(&user)
			if createResult.Error != nil {
				response.WriteJSON(c, response.NewResponse(2, nil, "用户注册失败"))
				return
			}
		} else {
			response.WriteJSON(c, response.NewResponse(2, nil, "数据库错误"))
			return
		}
	} else {
		// 用户存在，执行登录逻辑
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "手机号或密码不正确"))
			return
		}
	}

	// 登录成功，生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.PhoneNumber)
	if err != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "令牌生成失败"))
	} else {
		// 登录成功，返回统一的响应格式
		data := map[string]string{
			"userId":      strconv.Itoa(int(user.ID)),
			"phoneNumber": user.PhoneNumber,
			"token":       token,
		}
		response.WriteJSON(c, response.NewResponse(0, data, "一键注册登录成功！"))
	}
}
