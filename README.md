# seq-common-service

序列（seq）系列后端服务的公共基础库，被 `seq-user-service`、`seq-admin-service`、`seq-universe-service` 共同引用。

只放**与具体业务、具体基础设施无关**的通用能力；依赖 `internal/infra`（配置、数据库、缓存、日志）的代码（如 HTTP 中间件）仍留在各服务内。

## 包一览

| 包               | 说明                                                             |
| ---------------- | ---------------------------------------------------------------- |
| `pkg/errcode`    | 业务错误码类型 `Code`、业务错误 `Error`，以及错误码注册表        |
| `pkg/response`   | 统一 HTTP 响应体与输出方法                                       |
| `pkg/binder`     | gin 参数绑定封装，绑定失败直接输出中文参数错误                   |
| `pkg/validator`  | 自定义校验规则（phone/password/username/alphanumdash）与中文翻译 |
| `pkg/ctxkeys`    | gin Context / HTTP Header 键名常量                               |
| `pkg/jwt`        | Access / Refresh 双令牌签发与解析                                |
| `pkg/hash`       | bcrypt 密码哈希与强度校验                                        |
| `pkg/idgen`      | Snowflake ID 生成                                                |
| `pkg/pagination` | 分页参数归一化与 SQL Offset/Limit                                |

## 错误码约定

本库提供三段通用码，各服务在 `internal/errcodes` 中定义自己的业务码并在 `init` 注册：

| 码段           | 归属                           |
| -------------- | ------------------------------ |
| `1xxxx`        | 通用错误（本库）               |
| `3xxxx`        | 认证域（本库）                 |
| `4xxxx`        | 权限域（本库）                 |
| `2xxxx` 及其他 | 各服务自有领域码（服务内定义） |

服务侧写法：

```go
package errcodes

import "github.com/seq-tech/seq-common-service/pkg/errcode"

const (
    UserNotFound   errcode.Code = 20001
    UsernameExists errcode.Code = 20002
)

func init() {
    errcode.Register(
        errcode.Definition{Code: UserNotFound, Message: "用户不存在"},
        errcode.Definition{Code: UsernameExists, Message: "用户名已被占用"},
    )
}
```

`Register` 对重复码会 panic，以便在启动阶段暴露码段冲突。`Definition.HTTPStatus` 留空时按码段推导（`30004~30006` 为 401，其余业务码为 400）。

## 版本策略

语义化版本 tag，服务侧通过 `go get github.com/seq-tech/seq-common-service@vX.Y.Z` 升级。破坏性调整（错误码归属、公开 API 签名）需要升次版本号并同步三个服务。

## 本地开发

需要在服务里联调本库未发布的改动时，临时在服务的 `go.mod` 中加：

```
replace github.com/seq-tech/seq-common-service => ../seq-common-service
```

联调结束、tag 发布后必须移除该 replace。
