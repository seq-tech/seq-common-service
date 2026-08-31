// Package ctxkeys 集中定义 gin.Context / context.Context 中的共享键名，
// 避免各层用魔法字符串取值。
package ctxkeys

// gin.Context 键名。
const (
	RequestID      = "request_id"
	UserID         = "auth_user_id"
	Username       = "auth_username"
	TokenID        = "auth_token_id"
	TokenExpiresAt = "auth_token_exp"
	ClientIP       = "client_ip"
)

// HeaderRequestID 请求链路 ID 的 HTTP 头。
const HeaderRequestID = "X-Request-Id"
