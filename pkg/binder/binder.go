// Package binder 封装 gin 参数绑定，绑定失败时直接输出统一的中文参数错误。
package binder

import (
	"github.com/gin-gonic/gin"

	"github.com/seq-tech/seq-common-service/pkg/errcode"
	"github.com/seq-tech/seq-common-service/pkg/response"
	"github.com/seq-tech/seq-common-service/pkg/validator"
)

// JSON 绑定 JSON 请求体，返回 false 表示已输出错误响应，handler 应直接 return。
func JSON(c *gin.Context, req any) bool {
	return handle(c, c.ShouldBindJSON(req))
}

// Query 绑定 URL query 参数。
func Query(c *gin.Context, req any) bool {
	return handle(c, c.ShouldBindQuery(req))
}

// URI 绑定路径参数。
func URI(c *gin.Context, req any) bool {
	return handle(c, c.ShouldBindUri(req))
}

func handle(c *gin.Context, err error) bool {
	if err == nil {
		return true
	}
	response.Fail(c, errcode.New(errcode.InvalidParam).
		WithMessage("%s", validator.Translate(err)).
		WithCause(err))
	return false
}
