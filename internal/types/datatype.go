package types

type DataType struct {
	tstring		*string
	tint		*int
	tpath		*string
}

func TypeString(val string) DataType {
	return DataType{
		tstring: &val,
	}
}

func TypeInt(val int) DataType {
	return DataType{
		tint: &val,
	}
}

func TypePath(val string) DataType {
	return DataType{
		tpath: &val,
	}
}

func (sv DataType) IsString() bool {
	return sv.tstring != nil
}

func (sv DataType) IsInt() bool {
	return sv.tint != nil
}

func (sv DataType) IsPath() bool {
	return sv.tpath != nil
}

func (sv DataType) IsPrimary() bool {
	return sv.tstring != nil || sv.tint != nil
}

func (sv DataType) IsAny() bool {
	return sv.tstring != nil || sv.tint != nil || sv.tpath != nil
}

func (sv DataType) String() (string, bool) {
	if sv.tstring != nil {
		return *sv.tstring, true
	}

	return "", false
}

func (sv DataType) Int() (int, bool) {
	if sv.tint != nil {
		return *sv.tint, true
	}

	return 0, false
}

func (sv DataType) Path() (string, bool) {
	if sv.tpath != nil {
		return *sv.tpath, true
	}

	return "", false
}
