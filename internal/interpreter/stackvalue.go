package interpreter

type StackValue struct {
	tstring		*string
	tint		*int
	tpath		*string
}

func StackString(val string) StackValue {
	return StackValue{
		tstring: &val,
	}
}

func StackInt(val int) StackValue {
	return StackValue{
		tint: &val,
	}
}

func StackPath(val string) StackValue {
	return StackValue{
		tpath: &val,
	}
}

func (sv StackValue) IsString() bool {
	return sv.tstring != nil
}

func (sv StackValue) IsInt() bool {
	return sv.tint != nil
}

func (sv StackValue) IsPath() bool {
	return sv.tpath != nil
}

func (sv StackValue) IsPrimary() bool {
	return sv.tstring != nil || sv.tint != nil
}

func (sv StackValue) IsAny() bool {
	return sv.tstring != nil || sv.tint != nil || sv.tpath != nil
}

func (sv StackValue) String() (string, bool) {
	if sv.tstring != nil {
		return *sv.tstring, true
	}

	return "", false
}

func (sv StackValue) Int() (int, bool) {
	if sv.tint != nil {
		return *sv.tint, true
	}

	return 0, false
}

func (sv StackValue) Path() (string, bool) {
	if sv.tpath != nil {
		return *sv.tpath, true
	}

	return "", false
}
