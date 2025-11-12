package types

import "slices"

type ValueType uint8
const (
	ValueTypeNone ValueType = iota
	ValueTypeInt
	ValueTypeString
	ValueTypePath
)

func (vt ValueType) Format() string {
	return ValueTypeFormat(vt)
}

func (vt ValueType) Int() bool {
	return vt == ValueTypeInt
}

func (vt ValueType) String() bool {
	return vt == ValueTypeString
}

func (vt ValueType) Path() bool {
	return vt == ValueTypePath
}

func (vt ValueType) Primary() bool {
	return vt == ValueTypeInt || vt == ValueTypeString
}

func (vt ValueType) Any() bool {
	return vt == ValueTypeInt || vt == ValueTypeString || vt == ValueTypePath
}

func (vt ValueType) None() bool {
	return vt == ValueTypeNone
}

func (vt ValueType) SafeNone() bool {
	return !vt.Any()
}

func (vt ValueType) Is(vts...ValueType) bool {
	return slices.Contains(vts, vt)
}

func (vt *ValueType) Validate() bool {
	switch *vt {
	case ValueTypeInt:		return true
	case ValueTypeString:	return true
	case ValueTypePath:		return true
	case ValueTypeNone:		return true
	default:
		*vt = ValueTypeNone
		return false
	}
}

func ParseValueType(word *Word) ValueType {
	switch word.UnwrapName() {
	case "int":				return ValueTypeInt
	case "string":			return ValueTypeString
	case "path":			return ValueTypePath
	default:				return ValueTypeNone
	}
}

func ValueTypeFormat(vt ValueType) string {
	switch vt {
	case ValueTypeNone:		return "None"
	case ValueTypeInt:		return "Int"
	case ValueTypeString:	return "String"
	case ValueTypePath:		return "Path"
	default:				return "Unknown"
	}
}
