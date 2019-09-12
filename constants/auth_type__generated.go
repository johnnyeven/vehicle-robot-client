package constants

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidAuthType = errors.New("invalid AuthType")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("AuthType", map[string]string{
		"HOST": "主控端认证",
	})
}

func ParseAuthTypeFromString(s string) (AuthType, error) {
	switch s {
	case "":
		return AUTH_TYPE_UNKNOWN, nil
	case "HOST":
		return AUTH_TYPE__HOST, nil
	}
	return AUTH_TYPE_UNKNOWN, InvalidAuthType
}

func ParseAuthTypeFromLabelString(s string) (AuthType, error) {
	switch s {
	case "":
		return AUTH_TYPE_UNKNOWN, nil
	case "主控端认证":
		return AUTH_TYPE__HOST, nil
	}
	return AUTH_TYPE_UNKNOWN, InvalidAuthType
}

func (AuthType) EnumType() string {
	return "AuthType"
}

func (AuthType) Enums() map[int][]string {
	return map[int][]string{
		int(AUTH_TYPE__HOST): {"HOST", "主控端认证"},
	}
}
func (v AuthType) String() string {
	switch v {
	case AUTH_TYPE_UNKNOWN:
		return ""
	case AUTH_TYPE__HOST:
		return "HOST"
	}
	return "UNKNOWN"
}

func (v AuthType) Label() string {
	switch v {
	case AUTH_TYPE_UNKNOWN:
		return ""
	case AUTH_TYPE__HOST:
		return "主控端认证"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*AuthType)(nil)

func (v AuthType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidAuthType
	}
	return []byte(str), nil
}

func (v *AuthType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseAuthTypeFromString(string(bytes.ToUpper(data)))
	return
}
