package types

//go:generate libtools gen enum HolderDirection
// swagger:enum
type HolderDirection uint8

// 云台方向
const (
	HOLDER_DIRECTION_UNKNOWN   HolderDirection = iota
	HOLDER_DIRECTION__HORIZEN                  // 水平
	HOLDER_DIRECTION__VERTICAL                 // 垂直
)
