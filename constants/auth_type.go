package constants

//go:generate libtools gen enum AuthType
// swagger:enum
type AuthType uint8

// 认证类型
const (
	AUTH_TYPE_UNKNOWN AuthType = iota
	AUTH_TYPE__HOST            // 主控端认证
)
