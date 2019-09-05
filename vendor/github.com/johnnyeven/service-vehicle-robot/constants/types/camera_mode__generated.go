package types

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidCameraMode = errors.New("invalid CameraMode")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("CameraMode", map[string]string{
		"NORMAL":           "普通模式（仅视频传输）",
		"OBJECT_DETECTIVE": "物体识别模式",
	})
}

func ParseCameraModeFromString(s string) (CameraMode, error) {
	switch s {
	case "":
		return CAMERA_MODE_UNKNOWN, nil
	case "NORMAL":
		return CAMERA_MODE__NORMAL, nil
	case "OBJECT_DETECTIVE":
		return CAMERA_MODE__OBJECT_DETECTIVE, nil
	}
	return CAMERA_MODE_UNKNOWN, InvalidCameraMode
}

func ParseCameraModeFromLabelString(s string) (CameraMode, error) {
	switch s {
	case "":
		return CAMERA_MODE_UNKNOWN, nil
	case "普通模式（仅视频传输）":
		return CAMERA_MODE__NORMAL, nil
	case "物体识别模式":
		return CAMERA_MODE__OBJECT_DETECTIVE, nil
	}
	return CAMERA_MODE_UNKNOWN, InvalidCameraMode
}

func (CameraMode) EnumType() string {
	return "CameraMode"
}

func (CameraMode) Enums() map[int][]string {
	return map[int][]string{
		int(CAMERA_MODE__NORMAL):           {"NORMAL", "普通模式（仅视频传输）"},
		int(CAMERA_MODE__OBJECT_DETECTIVE): {"OBJECT_DETECTIVE", "物体识别模式"},
	}
}
func (v CameraMode) String() string {
	switch v {
	case CAMERA_MODE_UNKNOWN:
		return ""
	case CAMERA_MODE__NORMAL:
		return "NORMAL"
	case CAMERA_MODE__OBJECT_DETECTIVE:
		return "OBJECT_DETECTIVE"
	}
	return "UNKNOWN"
}

func (v CameraMode) Label() string {
	switch v {
	case CAMERA_MODE_UNKNOWN:
		return ""
	case CAMERA_MODE__NORMAL:
		return "普通模式（仅视频传输）"
	case CAMERA_MODE__OBJECT_DETECTIVE:
		return "物体识别模式"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*CameraMode)(nil)

func (v CameraMode) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidCameraMode
	}
	return []byte(str), nil
}

func (v *CameraMode) UnmarshalText(data []byte) (err error) {
	*v, err = ParseCameraModeFromString(string(bytes.ToUpper(data)))
	return
}
