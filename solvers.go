package lp

import "math"
import "fmt"

type softInEqData struct {
	slack              *Summand
	constraint         *Constraint
	minSlackConstraint *Constraint
}

type QPSolver struct {
	inEqSlackConstraints *softInEqList
	ls                   *LinearSpec
}

func newQPSolver(ls *LinearSpec) *QPSolver {
	qs := &QPSolver{}
	qs.inEqSlackConstraints = newsoftInEqList()
	qs.ls = ls

	return qs
}

func (self *QPSolver) ConstraintAdded(constraint *Constraint) bool {
	if self.isSoftInequality(constraint) == false {
		return true
	}

	data := new(softInEqData)
	coeff := -1.0
	if constraint.Op() == OperatorGE {
		coeff = 1
	}

	data.slack = &Summand{coeff, self.ls.AddVariable(nil)}
	data.slack.Var().SetRange(0, 20000)
	data.minSlackConstraint = self.ls.AddConstraint4([]float64{1.0}, []*Variable{data.slack.Var()},
		OperatorEQ, 0.0, constraint.PenaltyNeg(), constraint.PenaltyPos())

	data.constraint = constraint
	self.inEqSlackConstraints.AddItem(data)

	leftSide := constraint.LeftSide()
	leftSide.AddItem(data.slack)
	constraint.SetLeftSide(leftSide)

	return true
}

func (self *QPSolver) ConstraintRemoved(constraint *Constraint) bool {
	if self.isSoftInequality(constraint) == false {
		return true
	}

	for i := 0; i < self.inEqSlackConstraints.Len(); i++ {
		data := self.inEqSlackConstraints.GetAt(i)
		if data.constraint != constraint {
			continue
		}

		leftSide := constraint.LeftSide()
		leftSide.RemoveItem(data.slack)
		constraint.SetLeftSide(leftSide)
		self.destroyData(data)

		self.inEqSlackConstraints.RemoveItemAt(i)
		break
	}
	return true
}

func (self *QPSolver) destroyData(data *softInEqData) {
	self.ls.RemoveConstraint(data.minSlackConstraint)
	self.ls.RemoveVariable(data.slack.Var())

	data.slack = nil
}

func (self *QPSolver) isSoftInequality(constraint *Constraint) bool {
	if constraint.PenaltyNeg() <= 0 && constraint.PenaltyPos() <= 0 {
		return false
	}
	if constraint.Op() != OperatorEQ {
		return true
	}
	return false
}

type ActiveSetSolver struct {
	*QPSolver
	variables             *VariableList
	constraints           *ConstraintList
	variableGEConstraints *ConstraintList
	variableLEConstraints *ConstraintList
}

func NewActiveSetSolver(ls *LinearSpec) *ActiveSetSolver {
	ass := &ActiveSetSolver{}
	ass.QPSolver = newQPSolver(ls)

	ass.variables = ls.UsedVariables()
	ass.constraints = ls.Constraints()

	ass.variableGEConstraints = newConstraintList()
	ass.variableLEConstraints = newConstraintList()

	return ass
}

func (self *ActiveSetSolver) Solve() int {
	nConstraints := self.constraints.Len()
	nVariables := self.variables.Len()

	if nVariables > nConstraints {
		return ResultInfeasible
	}

	// First find an initial solution and the optimize it using the
	// active set method
	system := NewEquationSystem(nConstraints, nVariables+nConstraints)

	slackIndex := nVariables
	// set constraint matrix and add slack variables if necessary
	rowIndex := 0
	for c := 0; c < nConstraints; c++ {
		constraint := self.constraints.GetAt(c)
		if constraint.IsSoft() {
			continue
		}
		leftSide := constraint.LeftSide()
		*(system.B(rowIndex)) = constraint.RightSide()
		for sIndex := 0; sIndex < leftSide.Len(); sIndex++ {
			summand := leftSide.GetAt(sIndex)
			coefficient := summand.Coeff()
			*(system.A(rowIndex, summand.VariableIndex())) = coefficient
		}
		if constraint.Op() == OperatorLE {
			*(system.A(rowIndex, slackIndex)) = 1.0
			slackIndex++
		} else if constraint.Op() == OperatorGE {
			*(system.A(rowIndex, slackIndex)) = -1.0
			slackIndex++
		}
		rowIndex++
	}

	system.SetRows(rowIndex)
	system.RemoveLinearlyDependentRows()
	system.RemoveUnusedVariables()

	if !solveEq(system) {
		return ResultInfeasible
	}

	results := make([]float64, nVariables+nConstraints)
	system.Results(results, nVariables+nConstraints)
	optimizer := NewLayoutOptimizer(self.constraints, nVariables)
	optimizer.Solve(results)
    fmt.Println(results)

	// back to the variables
	for i := 0; i < nVariables; i++ {
		self.variables.GetAt(i).SetValue(results[i])
	}

	return ResultOptimal
}

func (self *ActiveSetSolver) VariableAdded(variable *Variable) bool {
	return true
}

func (self *ActiveSetSolver) VariableRemoved(variable *Variable) bool {
	self.variableGEConstraints.RemoveItemAt(variable.GlobalIndex())
	self.variableLEConstraints.RemoveItemAt(variable.GlobalIndex())
	return true
}

func (self *ActiveSetSolver) VariableRangeChanged(variable *Variable) bool {
	min := variable.Min()
	max := variable.Max()
	variableIndex := variable.GlobalIndex()

	constraintGE := self.variableGEConstraints.GetAt(variableIndex)
	constraintLE := self.variableLEConstraints.GetAt(variableIndex)

	if constraintGE == nil && min > -20000 {
		constraintGE = self.ls.AddConstraint2([]float64{1.0}, []*Variable{variable},
			OperatorGE, 0.0)
		if constraintGE == nil {
			return false
		}
		self.variableGEConstraints.RemoveItemAt(variableIndex)
		self.variableGEConstraints.AddItemAt(constraintGE, variableIndex)
	}

	if constraintLE == nil && max < 20000 {
		constraintLE = self.ls.AddConstraint2([]float64{1.0}, []*Variable{variable},
			OperatorLE, 20000)
		if constraintLE == nil {
			return false
		}
		self.variableLEConstraints.RemoveItemAt(variableIndex)
		self.variableLEConstraints.AddItemAt(constraintLE, variableIndex)
	}

	if constraintGE != nil {
		constraintGE.SetRightSide(min)
	}
	if constraintLE != nil {
		constraintLE.SetRightSide(max)
	}
	return true
}

func (self *ActiveSetSolver) ConstraintAdded(constraint *Constraint) bool {
	return self.QPSolver.ConstraintAdded(constraint)
}

func (self *ActiveSetSolver) ConstraintRemoved(constraint *Constraint) bool {
	return self.QPSolver.ConstraintRemoved(constraint)
}

func (self *ActiveSetSolver) LeftSideChanged(constraint *Constraint) bool {
	return true
}

func (self *ActiveSetSolver) RightSideChanged(constraint *Constraint) bool {
	return true
}

func (self *ActiveSetSolver) OperatorChanged(constraint *Constraint) bool {
	return true
}

func (self *ActiveSetSolver) SaveModel(fileName string) bool {
	return false
}

func (self *ActiveSetSolver) removeSoftConstraint(list *ConstraintList) {
	allConstraints := self.ls.Constraints()
	for i := 0; i < allConstraints.Len(); i++ {
		constraint := allConstraints.GetAt(i)
		if !constraint.IsSoft() {
			continue
		}
		if self.ls.RemoveConstraint(constraint) == true {
			list.AddItem(constraint)
		}
	}
}

func (self *ActiveSetSolver) addSoftConstraint(list *ConstraintList) {
	for i := 0; i < list.Len(); i++ {
		constraint := list.GetAt(i)
		if self.ls.AddConstraint(constraint) == false {
			constraint = nil
		}
	}
}

func (self *ActiveSetSolver) MinSize(width, height *Variable) Size {
	softConstraints := newConstraintList()
	self.removeSoftConstraint(softConstraints)

	heightConstraint := self.ls.AddConstraint4([]float64{1.0}, []*Variable{height},
		OperatorEQ, 0, 5, 5)
	widthConstraint := self.ls.AddConstraint4([]float64{1.0}, []*Variable{width},
		OperatorEQ, 0, 5, 5)
	result := self.Solve()
	self.ls.RemoveConstraint(heightConstraint)
	self.ls.RemoveConstraint(widthConstraint)

	self.addSoftConstraint(softConstraints)

	if result == ResultUnbounded {
		return Size{0, 0}
	}
	if result != ResultOptimal {
		fmt.Println("Could not solve the layout specification (%d).", result)
	}

	return Size{width.Value(), height.Value()}
}

func (self *ActiveSetSolver) MaxSize(width, height *Variable) Size {
	softConstraints := newConstraintList()
	self.removeSoftConstraint(softConstraints)

	hugeValue := 32000.00
	heightConstraint := self.ls.AddConstraint4([]float64{1.0}, []*Variable{height},
		OperatorEQ, hugeValue, 5, 5)
	widthConstraint := self.ls.AddConstraint4([]float64{1.0}, []*Variable{width},
		OperatorEQ, hugeValue, 5, 5)
	result := self.Solve()
	self.ls.RemoveConstraint(heightConstraint)
	self.ls.RemoveConstraint(widthConstraint)

	self.addSoftConstraint(softConstraints)

	if result == ResultUnbounded {
		return Size{math.MaxFloat64, math.MaxFloat64}
	}
	if result != ResultOptimal {
		fmt.Println("Could not solve the layout specification (%d).", result)
	}

	return Size{width.Value(), height.Value()}

}
