package types

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidNodeType = errors.New("invalid NodeType")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("NodeType", map[string]string{
		"HOST":  "主控端",
		"ROBOT": "机器人端",
	})
}

func ParseNodeTypeFromString(s string) (NodeType, error) {
	switch s {
	case "":
		return NODE_TYPE_UNKNOWN, nil
	case "HOST":
		return NODE_TYPE__HOST, nil
	case "ROBOT":
		return NODE_TYPE__ROBOT, nil
	}
	return NODE_TYPE_UNKNOWN, InvalidNodeType
}

func ParseNodeTypeFromLabelString(s string) (NodeType, error) {
	switch s {
	case "":
		return NODE_TYPE_UNKNOWN, nil
	case "主控端":
		return NODE_TYPE__HOST, nil
	case "机器人端":
		return NODE_TYPE__ROBOT, nil
	}
	return NODE_TYPE_UNKNOWN, InvalidNodeType
}

func (NodeType) EnumType() string {
	return "NodeType"
}

func (NodeType) Enums() map[int][]string {
	return map[int][]string{
		int(NODE_TYPE__HOST):  {"HOST", "主控端"},
		int(NODE_TYPE__ROBOT): {"ROBOT", "机器人端"},
	}
}
func (v NodeType) String() string {
	switch v {
	case NODE_TYPE_UNKNOWN:
		return ""
	case NODE_TYPE__HOST:
		return "HOST"
	case NODE_TYPE__ROBOT:
		return "ROBOT"
	}
	return "UNKNOWN"
}

func (v NodeType) Label() string {
	switch v {
	case NODE_TYPE_UNKNOWN:
		return ""
	case NODE_TYPE__HOST:
		return "主控端"
	case NODE_TYPE__ROBOT:
		return "机器人端"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*NodeType)(nil)

func (v NodeType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidNodeType
	}
	return []byte(str), nil
}

func (v *NodeType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseNodeTypeFromString(string(bytes.ToUpper(data)))
	return
}
