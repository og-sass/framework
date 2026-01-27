package xerr

type Config struct {
	LanguageDefault string                       `json:"language_default,optional"` // 默认语言
	ErrorMessages   map[string]map[string]string `json:"error_messages,optional"`   // 多语言错误信息
}
