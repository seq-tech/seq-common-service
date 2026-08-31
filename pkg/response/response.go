// Package response 定义统一 HTTP 响应结构与输出方法。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seq-tech/seq-common-service/pkg/ctxkeys"
	"github.com/seq-tech/seq-common-service/pkg/errcode"
)

// Body 统一响应体。
type Body struct {
	Code      errcode.Code `json:"code"`
	Message   string       `json:"message"`
	Data      any          `json:"data,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

// PageData 分页响应数据。
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// Success 输出成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:      errcode.Success,
		Message:   errcode.Success.Message(),
		Data:      data,
		RequestID: requestID(c),
	})
}

// OK 输出无数据的成功响应。
func OK(c *gin.Context) { Success(c, nil) }

// Page 输出分页成功响应。
func Page(c *gin.Context, list any, total int64, page, pageSize int) {
	Success(c, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}

// Fail 按业务错误输出失败响应。非业务错误会被归一化为服务器内部错误，
// 底层错误详情只写日志，不返回给客户端。
func Fail(c *gin.Context, err error) {
	be := errcode.From(err)
	if be == nil {
		OK(c)
		return
	}
	// 交给 accesslog / recovery 中间件记录原始错误。
	_ = c.Error(err) //nolint:errcheck // gin 收集错误链，无需处理返回值

	message := be.Message
	if be.Code == errcode.ServerError {
		message = errcode.ServerError.Message()
	}
	c.AbortWithStatusJSON(be.HTTPStatus(), Body{
		Code:      be.Code,
		Message:   message,
		RequestID: requestID(c),
	})
}

// FailCode 按错误码直接输出失败响应。
func FailCode(c *gin.Context, code errcode.Code) { Fail(c, errcode.New(code)) }

// FailCodeMsg 按错误码 + 自定义文案输出失败响应。
func FailCodeMsg(c *gin.Context, code errcode.Code, msg string) {
	Fail(c, errcode.New(code).WithMessage("%s", msg))
}

func requestID(c *gin.Context) string {
	if v, ok := c.Get(ctxkeys.RequestID); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
