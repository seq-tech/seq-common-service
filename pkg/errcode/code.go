// Package errcode 定义跨服务通用的业务错误码与业务错误类型。
//
// 约定：
//   - 0 表示成功；
//   - 1xxxx 通用错误、3xxxx 认证域、4xxxx 权限域由本包提供；
//   - 各服务自有的业务码在服务内定义，并在 init 阶段调用 Register 注册；
//   - handler 层统一通过 response.Fail 输出，不直接写 HTTP 状态码。
package errcode

import "net/http"

// Code 业务错误码。
type Code int

// 通用错误码。
const (
	Success        Code = 0
	ServerError    Code = 10000
	InvalidParam   Code = 10001
	Unauthorized   Code = 10002
	Forbidden      Code = 10003
	NotFound       Code = 10004
	TooManyRequest Code = 10005
	Timeout        Code = 10006
	Conflict       Code = 10007
	NotImplemented Code = 10008
)

// 认证域错误码。
const (
	InvalidCredential  Code = 30001
	InvalidToken       Code = 30004
	TokenExpired       Code = 30005
	RefreshTokenReused Code = 30006
	AccountLocked      Code = 30008
	OldPasswordWrong   Code = 30009
	WeakPassword       Code = 30010
	SameAsOldPassword  Code = 30011
)

// 权限域错误码。
const (
	RoleNotFound       Code = 40001
	RoleCodeExists     Code = 40002
	PermissionNotFound Code = 40003
	PermissionExists   Code = 40004
	BuiltinRoleProtect Code = 40005
	RoleInUse          Code = 40006
)

func init() {
	Register(
		Definition{Code: Success, Message: "成功", HTTPStatus: http.StatusOK},
		Definition{Code: ServerError, Message: "服务器内部错误", HTTPStatus: http.StatusInternalServerError},
		Definition{Code: InvalidParam, Message: "请求参数错误", HTTPStatus: http.StatusBadRequest},
		Definition{Code: Unauthorized, Message: "未认证或登录已失效", HTTPStatus: http.StatusUnauthorized},
		Definition{Code: Forbidden, Message: "没有访问权限", HTTPStatus: http.StatusForbidden},
		Definition{Code: NotFound, Message: "资源不存在", HTTPStatus: http.StatusNotFound},
		Definition{Code: TooManyRequest, Message: "请求过于频繁，请稍后再试", HTTPStatus: http.StatusTooManyRequests},
		Definition{Code: Timeout, Message: "请求超时", HTTPStatus: http.StatusGatewayTimeout},
		Definition{Code: Conflict, Message: "资源状态冲突", HTTPStatus: http.StatusConflict},
		Definition{Code: NotImplemented, Message: "该功能尚未实现", HTTPStatus: http.StatusNotImplemented},

		Definition{Code: InvalidCredential, Message: "账号或密码错误"},
		Definition{Code: InvalidToken, Message: "令牌无效"},
		Definition{Code: TokenExpired, Message: "令牌已过期"},
		Definition{Code: RefreshTokenReused, Message: "刷新令牌已失效，请重新登录"},
		Definition{Code: AccountLocked, Message: "登录失败次数过多，账号已临时锁定"},
		Definition{Code: OldPasswordWrong, Message: "原密码错误"},
		Definition{Code: WeakPassword, Message: "密码强度不足"},
		Definition{Code: SameAsOldPassword, Message: "新密码不能与原密码相同"},

		Definition{Code: RoleNotFound, Message: "角色不存在"},
		Definition{Code: RoleCodeExists, Message: "角色标识已存在"},
		Definition{Code: PermissionNotFound, Message: "权限不存在"},
		Definition{Code: PermissionExists, Message: "权限标识已存在"},
		Definition{Code: BuiltinRoleProtect, Message: "内置角色不允许修改或删除"},
		Definition{Code: RoleInUse, Message: "角色已被占用，无法删除"},
	)
}
