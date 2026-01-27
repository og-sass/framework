package xerr

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stringx"
	"sync"
	"sync/atomic"
)

type Error struct {
	Code ErrCode `json:"code"`
	Data any     `json:"data"`
	Msg  string  `json:"message"`
}

var (
	conf atomic.Value
	once sync.Once
)

func Must(c *Config) {
	once.Do(func() {
		conf.Store(c)
	})
}

// NewError 业务错误
func NewError(code ErrCode, data any) Error {
	return Error{
		Code: code,
		Data: data,
		Msg:  "",
	}
}

// NewParamError 参数错误
func NewParamError(formatMsg string, data any, formatMsgArgs ...any) Error {
	return Error{
		Code: ErrCodeParamError,
		Data: data,
		Msg:  fmt.Sprintf(formatMsg, formatMsgArgs...),
	}
}

// NewUnauthorizedError 未授权错误
func NewUnauthorizedError() Error {
	return Error{
		Code: ErrCodeUnauthorized,
		Data: nil,
		Msg:  "",
	}
}

// NewForbiddenError 禁止访问错误
func NewForbiddenError(formatMsg string, data any, formatMsgArgs ...any) Error {
	return Error{
		Code: ErrCodeForbidden,
		Data: data,
		Msg:  fmt.Sprintf(formatMsg, formatMsgArgs...),
	}
}

// NewServerInternalError 服务内部错误
func NewServerInternalError() Error {
	return Error{
		Code: ErrCodeServerInternalError,
		Data: nil,
		Msg:  "",
	}
}

// NewServiceUnreachableError 服务不可用错误
func NewServiceUnreachableError(formatMsg string, data any, formatMsgArgs ...any) Error {
	return Error{
		Code: ErrCodeServiceUnavailable,
		Data: data,
		Msg:  fmt.Sprintf(formatMsg, formatMsgArgs...),
	}
}

func (err Error) Error() string {
	return err.Msg
}

// GetMessage 获取多语言错误信息
func (err Error) GetMessage(language string) string {
	// 配置为空，返回默认错误信息
	confVal, ok := conf.Load().(*Config) // 原子读取
	if !ok || confVal == nil {
		return err.Msg
	}

	// 未找到code，返回默认错误信息
	languageMsgMap, ok := confVal.ErrorMessages[err.Code.String()]
	if !ok {
		return err.Msg
	}

	// 语言参数为空，则使用默认语言
	if stringx.HasEmpty(language) {
		language = confVal.LanguageDefault
	}
	// 语言参数还是为空，则使用默认语言
	if stringx.HasEmpty(language) {
		return err.Msg
	}

	// 返回多语言错误信息
	msg, ok := languageMsgMap[language]
	if ok {
		return msg
	}

	// 该语言未配置错误信息，尝试获取默认语言的错误信息
	if stringx.NotEmpty(confVal.LanguageDefault) && confVal.LanguageDefault != language {
		if msg, ok = languageMsgMap[confVal.LanguageDefault]; ok {
			return msg
		}
	}

	return err.Msg
}
