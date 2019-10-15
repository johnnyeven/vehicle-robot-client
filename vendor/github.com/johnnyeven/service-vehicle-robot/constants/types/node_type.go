package types

//go:generate libtools gen enum NodeType
// swagger:enum
type NodeType uint8

// 端类型
const (
	NODE_TYPE_UNKNOWN NodeType = iota
	NODE_TYPE__HOST            // 主控端
	NODE_TYPE__ROBOT           // 机器人端
)
