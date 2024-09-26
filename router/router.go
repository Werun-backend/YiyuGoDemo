package router

import (
	"awesomeProject1/controllers"
	"awesomeProject1/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppRouter 结构体，包含 Gin 路由器和数据库连接
type AppRouter struct {
	DB     *gorm.DB
	Engine *gin.Engine
}

// NewAppRouter 初始化并返回 AppRouter 实例
func NewAppRouter(db *gorm.DB) *AppRouter {
	r := gin.Default() // 创建 Gin 路由器实例

	r.Use(middleware.CORSMiddleware())

	// 初始化控制器并设置路由
	controllers.NewDiaryController(db).Routes(r)
	controllers.NewTagController(db).Routes(r)

	r.POST("/diary/login", controllers.NewUserController(db).Login)

	return &AppRouter{
		DB:     db,
		Engine: r,
	}
}

// Run 启动 Gin 服务器
func (ar *AppRouter) Run(port string) error {
	return ar.Engine.Run(port) // 使用 Gin 路由器的 Run 方法
}

func (ar *AppRouter) Use(handlerFunc gin.HandlerFunc) {
	ar.Engine.Use(handlerFunc)
}
