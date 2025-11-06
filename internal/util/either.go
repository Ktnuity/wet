package util

type Either[L, R any] struct {
	left *L
	right *R
}

func NewLeft[L, R any](val L) Either[L, R] {
	return Either[L, R]{left: &val}
}

func NewRight[L, R any](val R) Either[L, R] {
	return Either[L, R]{right: &val}
}

func (e Either[L, R]) IsLeft() bool { return e.left != nil }
func (e Either[L, R]) IsRight() bool { return e.right != nil }
func (e Either[L, R]) IsEither() bool { return e.left != nil || e.right != nil }

func (e Either[L, R]) UnpackLeft() L { return *e.left }
func (e Either[L, R]) UnpackRight() R { return *e.right }

func (e Either[L, R]) Left() (L, bool) {
	if e.left != nil {
		return *e.left, true
	}

	var left L
	return left, false
}

func (e Either[L, R]) Right() (R, bool) {
	if e.right != nil {
		return *e.right, true
	}

	var right R
	return right, false
}
