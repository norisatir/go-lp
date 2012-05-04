package lp

import "testing"
import "fmt"

func tstLinearSpec(t *testing.T) {
	fmt.Println("Test linear spec")
	ls := NewLinearSpec()
	x1 := ls.AddVariable(nil)
	x2 := ls.AddVariable(nil)

	ls.AddConstraint2([]float64{1.0}, []*Variable{x1}, OperatorLE, 108)
	c2 := ls.AddConstraint2([]float64{1.0}, []*Variable{x2}, OperatorGE, 113)

	fmt.Println("Num of Variables: ", ls.solver.variables.Len())
	fmt.Println("Num of Constraints: ", ls.solver.constraints.Len())

	ls.Solve()
	fmt.Println(ls.String())

	ls.RemoveConstraint(c2)
	ls.Solve()
	fmt.Println(ls.String())

	c2 = ls.AddConstraint2([]float64{1.0}, []*Variable{x2}, OperatorGE, 113)
	ls.Solve()
	fmt.Println(ls.String())
}

func printResults(valist *VariableList) {
	for i := 0; i < valist.Len(); i++ {
		fmt.Printf("Variable %v = %f\n", i, valist.GetAt(i).Value())
	}
}

func TestSoftConstraints(t *testing.T) {
	fmt.Println("Test SoftConstraints")

	ls := NewLinearSpec()

	x1 := ls.AddVariable(nil)
	x1.SetLabel("label_x1")

	x2 := ls.AddVariable(nil)
	x2.SetLabel("label_x2")

	x3 := ls.AddVariable(nil)
	x3.SetLabel("label_x3")

	ls.AddConstraint2([]float64{1.0}, []*Variable{x1}, OperatorEQ, 0)
	ls.AddConstraint2([]float64{1.0, -1.0}, []*Variable{x1, x2}, OperatorLE, 0)
	ls.AddConstraint2([]float64{1.0, -1.0}, []*Variable{x2, x3}, OperatorLE, 0)
	ls.AddConstraint2([]float64{1.0, -1.0}, []*Variable{x3, x1}, OperatorEQ, 20)

	ls.AddConstraint4([]float64{1.0, -1.0}, []*Variable{x2, x1}, OperatorEQ, 10, 5, 5)
	c6 := ls.AddConstraint4([]float64{1.0, -1.0}, []*Variable{x3, x2}, OperatorEQ, 5, 5, 5)

    printResults(ls.UsedVariables())

	ls.Solve()
	fmt.Println("ls: ", ls.String())
	printResults(ls.AllVariables())

	ls.RemoveConstraint(c6)
	ls.Solve()
	fmt.Println("ls: ", ls.String())
	printResults(ls.UsedVariables())
}
