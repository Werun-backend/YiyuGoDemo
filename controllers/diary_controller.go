package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"awesomeProject1/request"
	"awesomeProject1/response"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math"
	"strconv"
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
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	var diary models.Diary
	var diaryDto request.DiaryDto

	// 绑定请求中的 JSON 数据到 diary 变量中
	if err := c.ShouldBindJSON(&diaryDto); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的输入"))
		return
	}

	// 检查日记内容是否为空
	if diaryDto.Content == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "日记内容不能为空"))
		return
	}
	diary.Content = diaryDto.Content

	// 检查 user_id 是否在 users 表中存在
	var user models.User
	if err := ctrl.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "用户未找到"))
		return
	}
	diary.UserID = userID

	// 查找标签,未创建的接口返回错误
	var tags []models.Tag
	for _, tagId := range diaryDto.TagIds {
		var tag models.Tag
		if err := ctrl.DB.First(&tag, "id = ?", tagId).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		tags = append(tags, tag)
	}

	diary.Tags = tags

	// 使用 GORM 创建日记记录
	result := ctrl.DB.Create(&diary)

	// 检查是否有错误发生
	if result.Error != nil {
		// 如果有错误，返回错误信息和 HTTP 500 状态码

		response.WriteJSON(c, response.NewResponse(2, nil, "创建日记失败"))
		return
	}

	// 如果创建成功，返回创建的日记和 HTTP 201 状态码
	response.WriteJSON(c, response.NewResponse(0, diary, "日记创建成功"))
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	// QueryParams 用于绑定查询参数的结构体
	type QueryParams struct {
		Page      int    `form:"page"`
		PageSize  int    `form:"pageSize"`
		SortField string `form:"sortField"`
		SortBy    string `form:"sortBy"`
		QueryType string `form:"queryType"`
		TagID     string `form:"tagId"`
		Content   string `form:"content"`
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
	}

	// 绑定查询参数
	var queryParams QueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的查询参数"))
		return
	}

	// 临时调试输出
	log.Printf("QueryParams: %+v", queryParams)

	//验证分页参数
	if queryParams.Page < 1 || queryParams.PageSize < 1 {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的页码或页面大小"))
		return
	}

	// 构建基础查询
	var result *gorm.DB
	result = ctrl.DB.Model(&models.Diary{}).Where("user_id = ?", userID)

	// 根据标签 ID 查询
	if queryParams.TagID != "" {
		tagID, err := strconv.Atoi(queryParams.TagID)
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "无效的标签 ID"))
			return
		}
		result = result.Joins("JOIN diary_tags ON diaries.id = diary_tags.diary_id").
			Where("diary_tags.tag_id = ?", tagID)
	}

	// 根据内容查询
	if queryParams.Content != "" {
		result = result.Where("content LIKE ?", "%"+queryParams.Content+"%")
	}

	// 根据起始时间查询
	if queryParams.StartTime != "" {
		startTime, err := time.Parse("2006-01-02", queryParams.StartTime)
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "无效的开始时间"))
			return
		}
		result = result.Where("created_at >= ?", startTime)
	}

	// 根据结束时间查询
	if queryParams.EndTime != "" {
		endTime, err := time.Parse("2006-01-02", queryParams.EndTime)
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "无效的结束时间"))
			return
		}
		result = result.Where("created_at <= ?", endTime)
	}

	// 应用分页和排序
	offset := (queryParams.Page - 1) * queryParams.PageSize
	var diaries []models.Diary
	var count int64
	var err error

	// 先获取总数
	err = result.Preload("Tags").Count(&count).Error
	if err != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "获取日记列表失败"))
	} else {
		// 计算总页数
		totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))

		// 应用分页
		result = result.Preload("Tags").Offset(offset).Limit(queryParams.PageSize)

		// 应用排序
		if queryParams.SortField != "" && queryParams.SortBy != "" {
			sortStr := queryParams.SortField + " " + queryParams.SortBy
			result = result.Order(sortStr)
		}

		// 获取日记列表
		err = result.Find(&diaries).Error
		if err != nil {
			response.WriteJSON(c, response.NewResponse(2, nil, "获取日记列表失败"))
		} else {

			diaryResponses := make([]response.DiaryResponse, len(diaries))
			for i, diary := range diaries {
				diaryResponses[i] = response.DiaryResponse{
					ID:        diary.ID,
					Content:   diary.Content,
					CreatedAt: diary.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: diary.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				}
			}

			response.WriteJSON(c, response.NewResponse(0, gin.H{"diaries": diaryResponses, "currentPage": queryParams.Page, "totalPages": totalPages}, "获取日记列表成功"))
		}
	}
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
