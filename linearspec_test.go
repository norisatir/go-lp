package lp

import "testing"
import "fmt"

func TestLinearSpec(t *testing.T) {
	ls := NewLinearSpec()
	x1 := ls.AddVariable(nil)
	x2 := ls.AddVariable(nil)

	//c1 := ls.AddConstraint2([]float64{1.0}, []*Variable{x1}, OperatorLE, 108)
	ls.AddConstraint2([]float64{1.0}, []*Variable{x1}, OperatorLE, 108)
	c2 := ls.AddConstraint2([]float64{1.0}, []*Variable{x2}, OperatorGE, 113)

	ls.Solve()
	t.Log(ls.String())
	fmt.Println(ls.String())

	ls.RemoveConstraint(c2)
	ls.Solve()
	fmt.Println(ls.String())

	c2 = ls.AddConstraint2([]float64{1.0}, []*Variable{x2}, OperatorGE, 113)
	ls.Solve()
	fmt.Println(ls.String())
}
