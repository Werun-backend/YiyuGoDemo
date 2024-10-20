package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"awesomeProject1/request"
	"awesomeProject1/response"
	"awesomeProject1/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// DiaryController 处理日记相关的请求
type DiaryController struct {
	DB *gorm.DB
}

// NewDiaryController 创建日记控制器实例
func NewDiaryController(db *gorm.DB) *DiaryController {
	return &DiaryController{DB: db}
}

// Routes 注册日记相关的路由
func (ctrl *DiaryController) Routes(r *gin.Engine) {
	//添加jwt校验
	diaries := r.Group("/diaries", middleware.JwtMiddleware())
	{
		diaries.POST("", ctrl.CreateDiary)       // 创建日记
		diaries.GET("", ctrl.GetDiaries)         // 获取日记列表
		diaries.GET("/:id", ctrl.GetDiary)       // 根据 ID 获取日记
		diaries.PUT("/:id", ctrl.UpdateDiary)    // 更新日记
		diaries.DELETE("/:id", ctrl.DeleteDiary) // 删除日记
		//diaries.POST("/login", ctrl.Login)       // 用户登录
	}
}

// CreateDiary 创建新的日记
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
	var diaryDto request.DiaryDto
	// 绑定请求中的 JSON 数据到 diaryDto 变量中
	if err := c.ShouldBindJSON(&diaryDto); err != nil {
		response.WriteJSON(c, response.UserErrorResponse(nil, "无效的输入"))
		return
	}
	// 检查日记内容是否为空
	if diaryDto.Content == "" {
		response.WriteJSON(c, response.UserErrorResponse(nil, "日记内容不能为空"))
		return
	}
	// 查找标签,未创建的接口返回错误
	var tags []models.Tag
	for _, tagId := range diaryDto.TagIds {
		var tag models.Tag
		if err := ctrl.DB.First(&tag, "id = ?", tagId).Error; err != nil {
			response.WriteJSON(c, response.UserErrorResponse(nil, "标签未找到"))
			return
		}
		tags = append(tags, tag)
	}
	// 获取当前用户ID
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	// 创建日记记录
	diary := models.Diary{
		Content: diaryDto.Content,
		Tags:    tags,
		UserID:  userID, // 设置外键
	}
	// 使用 GORM 创建日记记录
	result := ctrl.DB.Create(&diary)
	// 检查是否有错误发生
	if result.Error != nil {
		response.WriteJSON(c, response.InternalErrorResponse(nil, "创建日记失败"))
		return
	}
	// 转换为 DTO
	diaryDto = request.DiaryToDto(diary)
	//创建成功返回响应
	response.WriteJSON(c, response.SuccessResponse(diaryDto))
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	// 绑定查询参数
	var queryParams struct {
		Page     int `form:"page" binding:"required,gte=1"`
		PageSize int `form:"pageSize" binding:"required,gte=1"`
		//UserID    string `form:"userId" binding:"required"`
		TagID     string `form:"tagId"`
		Content   string `form:"content"`
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
		SortField string `form:"sortField"`
		SortBy    string `form:"sortBy"`
	}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.WriteJSON(c, response.UserErrorNoMsgResponse("无效的查询参数"))
		fmt.Printf(err.Error())
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.WriteJSON(c, response.UserErrorNoMsgResponse("用户未登录"))
		c.Abort()
		return
	}
	// 定义 modifier 函数
	modifier := func(db *gorm.DB) *gorm.DB {
		query := db.Where("user_id = ?", userID)
		if queryParams.TagID != "" {
			query = query.Where("tag_id = ?", queryParams.TagID)
		}
		if queryParams.Content != "" {
			query = query.Where("content LIKE ?", "%"+queryParams.Content+"%")
		}
		if queryParams.StartTime != "" {
			query = query.Where("created_at >= ?", queryParams.StartTime)
		}
		if queryParams.EndTime != "" {
			query = query.Where("created_at <= ?", queryParams.EndTime)
		}
		// 设置默认排序字段为创建时间，排序方式为降序
		if queryParams.SortField == "" || queryParams.SortBy == "" {
			query = query.Order("created_at desc")
		} else {
			query = query.Order(queryParams.SortField + " " + queryParams.SortBy)
		}
		return query
	}

	// 调用分页服务
	paginationResult, err := service.Paginate(ctrl.DB, queryParams.Page, queryParams.PageSize, &models.Diary{}, modifier)
	if err != nil {
		response.WriteJSON(c, response.InternalErrorResponse(nil, "获取日记列表失败"))
		return
	}

	// 转换日记数据为DiaryVo
	diaryVos := make([]response.DiaryVo, len(paginationResult.Records))
	for i, diary := range paginationResult.Records {
		diaryVos[i].Copy(diary)
	}
	// 创建包含分页信息的响应VO
	paginatedDiaryVo := response.PaginatedDiaryVo{}
	paginatedDiaryVo.Copy(diaryVos, response.PaginationData{
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		TotalCount: paginationResult.Total,
		TotalPages: paginationResult.TotalPages,
	})

	// 返回分页响应
	response.WriteJSON(c, response.NewResponse(200, paginatedDiaryVo, "success"))
}

// UpdateDiary 更新日记内容
func (ctrl *DiaryController) UpdateDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var update models.Diary
	if err := c.BindJSON(&update); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的更新数据"))
		return
	}

	var diary models.Diary
	result := ctrl.DB.First(&diary, "id = ? AND user_id = ?", diaryID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "更新日记时数据库错误"))
		return
	}

	//先判断标签是否存在
	var tags []models.Tag
	for _, tag := range update.Tags {
		var tag1 models.Tag
		if err := ctrl.DB.First(&tag1, "name = ?", tag.Name).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		tags = append(tags, tag1)
	}
	diary.Tags = tags

	// 更新字段...
	diary.Content = update.Content // 假设只更新内容字段
	diary.UpdatedAt = time.Now()   // 更新更新时间

	result = ctrl.DB.Save(&diary)
	if result.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "更新日记失败"))
		return
	}
	//更新标签
	err := ctrl.DB.Model(&diary).Association("Tags").Replace(&tags)
	if err != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "更新标签失败"))
		return
	}
	response.WriteJSON(c, response.NewResponse(0, diary, "日记更新成功"))
}

// GetDiary 根据ID获取单个日记
func (ctrl *DiaryController) GetDiary(c *gin.Context) {
	diaryID := c.Param("id") // 从URL参数中获取日记ID
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var diary models.Diary
	result := ctrl.DB.Where("user_id = ?", userID).Preload("Tags").Take(&diary, diaryID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到"))
		} else {
			response.WriteJSON(c, response.NewResponse(2, nil, "获取日记时数据库错误"))
		}
		return
	}
	response.WriteJSON(c, response.NewResponse(0, diary, "获取日记成功"))
}

// DeleteDiary 软删除日记（标记为已删除）
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	// 尝试获取日记记录
	var diary models.Diary
	result := ctrl.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", diaryID, userID).First(&diary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到或已被删除"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "删除日记时数据库错误"))
		return
	}

	// 执行软删除
	updateResult := ctrl.DB.Model(&diary).Update("deleted_at", time.Now())
	if updateResult.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "删除日记失败"))
	} else if updateResult.RowsAffected == 0 {
		response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到或已被删除"))
	} else {
		response.WriteJSON(c, response.NewResponse(0, nil, "日记删除成功"))
	}
}
