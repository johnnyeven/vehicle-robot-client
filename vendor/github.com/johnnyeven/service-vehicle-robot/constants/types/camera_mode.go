package types

//go:generate libtools gen enum CameraMode
// swagger:enum
type CameraMode uint8

// 摄像头模式
const (
	CAMERA_MODE_UNKNOWN           CameraMode = iota
	CAMERA_MODE__NORMAL                      // 普通模式（仅视频传输）
	CAMERA_MODE__OBJECT_DETECTIVE            // 物体识别模式
)
