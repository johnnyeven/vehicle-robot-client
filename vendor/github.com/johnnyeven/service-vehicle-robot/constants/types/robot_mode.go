package types

//go:generate libtools gen enum RobotMode
// swagger:enum
type RobotMode uint8

// 机器人的AI模式
const (
	ROBOT_MODE_UNKNOWN RobotMode = iota
	ROBOT_MODE__MANUAL           // 人工控制
	ROBOT_MODE__SEARCH           // 搜寻模式
)
