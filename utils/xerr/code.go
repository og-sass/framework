package xerr

import (
	"fmt"
)

// 通用错误码
const (
	ErrCodeSuccess             ErrCode = 0   // 成功
	ErrCodeFail                ErrCode = 1   // 失败
	ErrCodeParamError          ErrCode = 400 // 参数错误
	ErrCodeUnauthorized        ErrCode = 401 // 未授权
	ErrCodeForbidden           ErrCode = 403 // 禁止
	ErrCodeNotFound            ErrCode = 404 // 未找到
	ErrCodeServerInternalError ErrCode = 500 // 服务器内部错误
	ErrCodeServiceUnavailable  ErrCode = 503 // 服务不可用
)

type ErrCode int

func (s ErrCode) Int() int {
	return int(s)
}

func (s ErrCode) String() string {
	return fmt.Sprintf("%d", int(s))
}
