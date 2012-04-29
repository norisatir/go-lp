package lp

//import "github.com/norisatir/go-lp/lpsolve"
//import "math"

const (
	OperatorLE = iota
	OperatorGE
	OperatorEQ
)

const (
	OptMinimize = 0
	OptMaximize = 1
)

const (
	ResultNoMemory = -2
	ResultError    = -1
	ResultOptimal  = iota
	ResultSubOptimal
	ResultInfeasible
	ResultUnbounded
	ResultDegenerate
	ResultNumFailure
	ResultUserAbort
	ResultTimeout
	ResultPresolve
	ResultProcFail
	ResultProcBreak
	ResultFeasFound
	ResultNoFeasFound
)

type Size struct {
	W, H float64
}

type SolverLike interface {
	Solve() int
	VariableAdded(*Variable) bool
	VariableRemoved(*Variable) bool
	VariableRangeChanged(*Variable) bool

	ConstraintAdded(*Constraint) bool
	ConstraintRemoved(*Constraint) bool
	LeftSideChanged(*Constraint) bool
	RightSideChanged(*Constraint) bool
	OperatorChanged(*Constraint) bool
	SaveModel(filename string) bool

	MinSize(width, height *Variable) Size
	MaxSize(width, height *Variable) Size
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
