package types

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidHolderDirection = errors.New("invalid HolderDirection")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("HolderDirection", map[string]string{
		"HORIZEN":  "水平",
		"VERTICAL": "垂直",
	})
}

func ParseHolderDirectionFromString(s string) (HolderDirection, error) {
	switch s {
	case "":
		return HOLDER_DIRECTION_UNKNOWN, nil
	case "HORIZEN":
		return HOLDER_DIRECTION__HORIZEN, nil
	case "VERTICAL":
		return HOLDER_DIRECTION__VERTICAL, nil
	}
	return HOLDER_DIRECTION_UNKNOWN, InvalidHolderDirection
}

func ParseHolderDirectionFromLabelString(s string) (HolderDirection, error) {
	switch s {
	case "":
		return HOLDER_DIRECTION_UNKNOWN, nil
	case "水平":
		return HOLDER_DIRECTION__HORIZEN, nil
	case "垂直":
		return HOLDER_DIRECTION__VERTICAL, nil
	}
	return HOLDER_DIRECTION_UNKNOWN, InvalidHolderDirection
}

func (HolderDirection) EnumType() string {
	return "HolderDirection"
}

func (HolderDirection) Enums() map[int][]string {
	return map[int][]string{
		int(HOLDER_DIRECTION__HORIZEN):  {"HORIZEN", "水平"},
		int(HOLDER_DIRECTION__VERTICAL): {"VERTICAL", "垂直"},
	}
}
func (v HolderDirection) String() string {
	switch v {
	case HOLDER_DIRECTION_UNKNOWN:
		return ""
	case HOLDER_DIRECTION__HORIZEN:
		return "HORIZEN"
	case HOLDER_DIRECTION__VERTICAL:
		return "VERTICAL"
	}
	return "UNKNOWN"
}

func (v HolderDirection) Label() string {
	switch v {
	case HOLDER_DIRECTION_UNKNOWN:
		return ""
	case HOLDER_DIRECTION__HORIZEN:
		return "水平"
	case HOLDER_DIRECTION__VERTICAL:
		return "垂直"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*HolderDirection)(nil)

func (v HolderDirection) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidHolderDirection
	}
	return []byte(str), nil
}

func (v *HolderDirection) UnmarshalText(data []byte) (err error) {
	*v, err = ParseHolderDirectionFromString(string(bytes.ToUpper(data)))
	return
}
