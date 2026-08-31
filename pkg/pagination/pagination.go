// Package pagination 提供统一的分页参数处理。
package pagination

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// Query 分页请求参数，可内嵌到各模块的列表 DTO 中。
type Query struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// Normalize 修正非法分页参数，防止越界查询。
func (q *Query) Normalize() {
	if q.Page < defaultPage {
		q.Page = defaultPage
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
}

// Offset 返回 SQL OFFSET。
func (q Query) Offset() int {
	q.Normalize()
	return (q.Page - 1) * q.PageSize
}

// Limit 返回 SQL LIMIT。
func (q Query) Limit() int {
	q.Normalize()
	return q.PageSize
}
