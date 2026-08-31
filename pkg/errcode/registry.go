package errcode

import (
	"fmt"
	"net/http"
	"sync"
)

// Definition 描述一个业务错误码的默认文案与 HTTP 状态码。
// HTTPStatus 为 0 时按码段推导。
type Definition struct {
	Code       Code
	Message    string
	HTTPStatus int
}

var (
	registryMu sync.RWMutex
	messages   = make(map[Code]string)
	httpStatus = make(map[Code]int)
)

// Register 注册业务错误码，重复注册同一码会 panic，以便在启动阶段暴露码段冲突。
// 各服务应在包级 init 中调用。
func Register(defs ...Definition) {
	registryMu.Lock()
	defer registryMu.Unlock()
	for _, d := range defs {
		if _, ok := messages[d.Code]; ok {
			panic(fmt.Sprintf("errcode: duplicate registration for code %d", d.Code))
		}
		messages[d.Code] = d.Message
		if d.HTTPStatus != 0 {
			httpStatus[d.Code] = d.HTTPStatus
		}
	}
}

// Message 返回错误码对应的默认文案，未注册的码回落为服务器内部错误文案。
func (c Code) Message() string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	if msg, ok := messages[c]; ok {
		return msg
	}
	return messages[ServerError]
}

// HTTPStatus 返回错误码对应的 HTTP 状态码。
// 业务错误（非通用段）统一返回 400，由 code 字段区分具体原因。
func (c Code) HTTPStatus() int {
	registryMu.RLock()
	s, ok := httpStatus[c]
	registryMu.RUnlock()
	if ok {
		return s
	}
	switch {
	case c >= 30004 && c <= 30006:
		return http.StatusUnauthorized
	case c >= 40001 && c <= 40006:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
