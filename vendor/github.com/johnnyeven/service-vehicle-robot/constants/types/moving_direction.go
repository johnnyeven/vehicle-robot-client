package types

//go:generate libtools gen enum MovingDirection
// swagger:enum
type MovingDirection uint8

// 移动方向
const (
	MOVING_DIRECTION_UNKNOWN     MovingDirection = iota
	MOVING_DIRECTION__FORWARD                    // 前进
	MOVING_DIRECTION__BACKWARD                   // 后退
	MOVING_DIRECTION__TURN_LEFT                  // 左转
	MOVING_DIRECTION__TURN_RIGHT                 // 右转
	MOVING_DIRECTION__STOP                       // 停止
)
