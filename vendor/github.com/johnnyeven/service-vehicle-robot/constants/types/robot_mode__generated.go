package types

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidRobotMode = errors.New("invalid RobotMode")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("RobotMode", map[string]string{
		"MANUAL": "人工控制",
		"SEARCH": "搜寻模式",
	})
}

func ParseRobotModeFromString(s string) (RobotMode, error) {
	switch s {
	case "":
		return ROBOT_MODE_UNKNOWN, nil
	case "MANUAL":
		return ROBOT_MODE__MANUAL, nil
	case "SEARCH":
		return ROBOT_MODE__SEARCH, nil
	}
	return ROBOT_MODE_UNKNOWN, InvalidRobotMode
}

func ParseRobotModeFromLabelString(s string) (RobotMode, error) {
	switch s {
	case "":
		return ROBOT_MODE_UNKNOWN, nil
	case "人工控制":
		return ROBOT_MODE__MANUAL, nil
	case "搜寻模式":
		return ROBOT_MODE__SEARCH, nil
	}
	return ROBOT_MODE_UNKNOWN, InvalidRobotMode
}

func (RobotMode) EnumType() string {
	return "RobotMode"
}

func (RobotMode) Enums() map[int][]string {
	return map[int][]string{
		int(ROBOT_MODE__MANUAL): {"MANUAL", "人工控制"},
		int(ROBOT_MODE__SEARCH): {"SEARCH", "搜寻模式"},
	}
}
func (v RobotMode) String() string {
	switch v {
	case ROBOT_MODE_UNKNOWN:
		return ""
	case ROBOT_MODE__MANUAL:
		return "MANUAL"
	case ROBOT_MODE__SEARCH:
		return "SEARCH"
	}
	return "UNKNOWN"
}

func (v RobotMode) Label() string {
	switch v {
	case ROBOT_MODE_UNKNOWN:
		return ""
	case ROBOT_MODE__MANUAL:
		return "人工控制"
	case ROBOT_MODE__SEARCH:
		return "搜寻模式"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*RobotMode)(nil)

func (v RobotMode) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidRobotMode
	}
	return []byte(str), nil
}

func (v *RobotMode) UnmarshalText(data []byte) (err error) {
	*v, err = ParseRobotModeFromString(string(bytes.ToUpper(data)))
	return
}
