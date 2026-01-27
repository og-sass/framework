package metadata

import (
	"context"
	"github.com/spf13/cast"
)

type Metadata string

const (
	Authorization      Metadata = "Authorization"       // 授权
	UserId             Metadata = "user_id"             // 用户id
	Username           Metadata = "username"            // 用户名
	ChannelId          Metadata = "channel_id"          // 渠道id
	SiteId             Metadata = "site_id"             // 站点id
	Language           Metadata = "language"            // 语言
	IP                 Metadata = "ip"                  // ip
	Currency           Metadata = "currency"            // 币种
	Domain             Metadata = "domain"              // 域名
	Region             Metadata = "region"              // 区域
	DeviceId           Metadata = "Device-Id"           // 设备id
	DeviceType         Metadata = "Device-Type"         // 设备类型
	DeviceOS           Metadata = "Device-OS"           // 设备操作系统
	BrowserFingerprint Metadata = "browser_fingerprint" // 浏览器指纹
	AppVersion         Metadata = "App-Version"         // app版本
	UserAgent          Metadata = "User-Agent"          // 浏览器用户代理
)

// GetKey 获取元数据key
func (s Metadata) GetKey() string {
	return string(s)
}

// GetValue 获取元数据
func (s Metadata) GetValue(ctx context.Context) any {
	return ctx.Value(s)
}

// GetString 获取元数据字符串
func (s Metadata) GetString(ctx context.Context) string {
	return cast.ToString(ctx.Value(s))
}

// GetInt64 获取元数据int64
func (s Metadata) GetInt64(ctx context.Context) int64 {
	return cast.ToInt64(ctx.Value(s))
}

// SetValue 设置元数据
func (s Metadata) SetValue(ctx context.Context, val any) context.Context {
	return context.WithValue(ctx, s, val)
}

// SetValues 批量设置元数据
func SetValues(ctx context.Context, vals map[Metadata]any) context.Context {
	for k, v := range vals {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}
