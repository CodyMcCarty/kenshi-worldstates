package internal

type BoolExpr func() bool // lazy boolean (re-evaluated each time)
type Cond func() bool     // condition (also lazy)

func Equal(b BoolExpr, v bool) Cond { return func() bool { return b() == v } }

//	func And(cs ...Cond) Cond {
//		return func() bool {
//			for _, c := range cs {
//				if !c() {
//					return false
//				}
//			}
//			return true
//		}
//	}

func allTrue(cs []Cond) bool {
	for _, c := range cs {
		if !c() {
			return false
		}
	}
	return true
}
