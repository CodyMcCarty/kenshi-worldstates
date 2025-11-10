package internal

const (
	CheckMark = "\u2713"
	CrossMark = "\u2717"
)

type BoolExpr func() bool

func allTrue(rs []Cond) bool {
	for _, r := range rs {
		eval := r.Eval()
		if !eval == r.Want {
			return false
		}
	}
	return true
}

func AppendUnique[T comparable](slice []T, val T) []T {
	for _, v := range slice {
		if v == val {
			return slice // already there, do nothing
		}
	}
	return append(slice, val)
}
