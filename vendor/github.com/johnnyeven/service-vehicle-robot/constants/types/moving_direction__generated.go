package types

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidMovingDirection = errors.New("invalid MovingDirection")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("MovingDirection", map[string]string{
		"BACKWARD":   "后退",
		"FORWARD":    "前进",
		"STOP":       "停止",
		"TURN_LEFT":  "左转",
		"TURN_RIGHT": "右转",
	})
}

func ParseMovingDirectionFromString(s string) (MovingDirection, error) {
	switch s {
	case "":
		return MOVING_DIRECTION_UNKNOWN, nil
	case "BACKWARD":
		return MOVING_DIRECTION__BACKWARD, nil
	case "FORWARD":
		return MOVING_DIRECTION__FORWARD, nil
	case "STOP":
		return MOVING_DIRECTION__STOP, nil
	case "TURN_LEFT":
		return MOVING_DIRECTION__TURN_LEFT, nil
	case "TURN_RIGHT":
		return MOVING_DIRECTION__TURN_RIGHT, nil
	}
	return MOVING_DIRECTION_UNKNOWN, InvalidMovingDirection
}

func ParseMovingDirectionFromLabelString(s string) (MovingDirection, error) {
	switch s {
	case "":
		return MOVING_DIRECTION_UNKNOWN, nil
	case "后退":
		return MOVING_DIRECTION__BACKWARD, nil
	case "前进":
		return MOVING_DIRECTION__FORWARD, nil
	case "停止":
		return MOVING_DIRECTION__STOP, nil
	case "左转":
		return MOVING_DIRECTION__TURN_LEFT, nil
	case "右转":
		return MOVING_DIRECTION__TURN_RIGHT, nil
	}
	return MOVING_DIRECTION_UNKNOWN, InvalidMovingDirection
}

func (MovingDirection) EnumType() string {
	return "MovingDirection"
}

func (MovingDirection) Enums() map[int][]string {
	return map[int][]string{
		int(MOVING_DIRECTION__BACKWARD):   {"BACKWARD", "后退"},
		int(MOVING_DIRECTION__FORWARD):    {"FORWARD", "前进"},
		int(MOVING_DIRECTION__STOP):       {"STOP", "停止"},
		int(MOVING_DIRECTION__TURN_LEFT):  {"TURN_LEFT", "左转"},
		int(MOVING_DIRECTION__TURN_RIGHT): {"TURN_RIGHT", "右转"},
	}
}
func (v MovingDirection) String() string {
	switch v {
	case MOVING_DIRECTION_UNKNOWN:
		return ""
	case MOVING_DIRECTION__BACKWARD:
		return "BACKWARD"
	case MOVING_DIRECTION__FORWARD:
		return "FORWARD"
	case MOVING_DIRECTION__STOP:
		return "STOP"
	case MOVING_DIRECTION__TURN_LEFT:
		return "TURN_LEFT"
	case MOVING_DIRECTION__TURN_RIGHT:
		return "TURN_RIGHT"
	}
	return "UNKNOWN"
}

func (v MovingDirection) Label() string {
	switch v {
	case MOVING_DIRECTION_UNKNOWN:
		return ""
	case MOVING_DIRECTION__BACKWARD:
		return "后退"
	case MOVING_DIRECTION__FORWARD:
		return "前进"
	case MOVING_DIRECTION__STOP:
		return "停止"
	case MOVING_DIRECTION__TURN_LEFT:
		return "左转"
	case MOVING_DIRECTION__TURN_RIGHT:
		return "右转"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*MovingDirection)(nil)

func (v MovingDirection) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidMovingDirection
	}
	return []byte(str), nil
}

func (v *MovingDirection) UnmarshalText(data []byte) (err error) {
	*v, err = ParseMovingDirectionFromString(string(bytes.ToUpper(data)))
	return
}
