package interpreter

import "github.com/ktnuity/wet/internal/types"

type Proc struct {
	Name		string
	Start		int64
	End			int64
	Token		*types.Token
}
