package request

// DiaryDto 代表日记的请求模型
type DiaryDto struct {
	Content string `json:"content"` // 日记内容
	TagIds  []uint `json:"tagIds"`
}
