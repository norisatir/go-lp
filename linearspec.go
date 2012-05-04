package lp

import "fmt"

type LinearSpec struct {
	variables     *VariableList
	usedVariables *VariableList
	constraints   *ConstraintList
	result        int
	solvingTime   float64
	solver        *ActiveSetSolver
}

func NewLinearSpec() *LinearSpec {
	ls := &LinearSpec{}
	ls.result = ResultError
	ls.solvingTime = 0
	ls.variables = newVariableList()
	ls.usedVariables = newVariableList()
	ls.constraints = newConstraintList()

	ls.solver = NewActiveSetSolver(ls)

	return ls
}

// Adds a new variable to the specification
// if v == 0 then create new default variable in return it.
// Otherwise the returned variable is v.
func (self *LinearSpec) AddVariable(v *Variable) *Variable {
	var var1 *Variable = v
	if var1 == nil {
		var1 = newVariable(self)
	}

	if var1.IsValid() {
		return nil
	}

	self.variables.AddItem(var1)

	if !self.solver.VariableAdded(var1) {
		self.variables.RemoveItem(var1)
		return nil
	}
	var1.isValid = true

	if !self.UpdateRange(var1) {
		self.RemoveVariable(var1)
		return nil
	}

	return var1
}

func (self *LinearSpec) RemoveVariable(v *Variable) bool {

	// must be called first otherwise the index is invalid
	if self.solver.VariableRemoved(v) == false {
		return false
	}

	// do we know the variable?
	if self.variables.RemoveItem(v) == false {
		return false
	}
	self.usedVariables.RemoveItem(v)
	v.isValid = false

	// Invalidate all constraints that use this variable
	markedForInvalidation := newConstraintList()
	constraints := self.Constraints()

	for i := 0; i < constraints.Len(); i++ {
		c := constraints.GetAt(i)
		if !c.IsValid() {
			continue
		}

		summands := c.LeftSide()
		for j := 0; j < summands.Len(); j++ {
			s := summands.GetAt(j)
			if s.Var() == v {
				markedForInvalidation.AddItem(c)
				break
			}
		}
	}

	for i := 0; i < markedForInvalidation.Len(); i++ {
		self.RemoveConstraint(markedForInvalidation.GetAt(i))
	}
	return true
}

func (self *LinearSpec) IndexOf(v *Variable) int {
	return self.usedVariables.IndexOf(v)
}

func (self *LinearSpec) GlobalIndexOf(v *Variable) int {
	return self.variables.IndexOf(v)
}

func (self *LinearSpec) UpdateRange(v *Variable) bool {
	if !self.solver.VariableRangeChanged(v) {
		return false
	}
	return true
}

func (self *LinearSpec) AddConstraint(c *Constraint) bool {
	self.constraints.AddItem(c)

	leftSide := c.LeftSide()

	for i := 0; i < leftSide.Len(); i++ {
		v := leftSide.GetAt(i).Var()
        if v.AddReference() == 1 {
            self.usedVariables.AddItem(v)
        }
	}

	if !self.solver.ConstraintAdded(c) {
		self.RemoveConstraint(c)
		return false
	}
	c.isValid = true
	return true
}

func (self *LinearSpec) RemoveConstraint(c *Constraint) bool {
	self.solver.ConstraintRemoved(c)
    self.constraints.RemoveItem(c)
	c.isValid = false

	leftSide := c.LeftSide()
	for i := 0; i < leftSide.Len(); i++ {
        v := leftSide.GetAt(i).Var()
        if v.RemoveReference() == 0 {
            self.usedVariables.RemoveItem(v)
        }
	}

	return true
}

func (self *LinearSpec) AddConstraint1(summands *SummandList, opType int, rightSide float64) *Constraint {
	return self.AddConstraint3(summands, opType, rightSide, -1, -1)
}

func (self *LinearSpec) AddConstraint2(coeffs []float64, vars []*Variable,
	opType int, rightSide float64) *Constraint {
	return self.AddConstraint4(coeffs, vars, opType, rightSide, -1, -1)
}

func (self *LinearSpec) AddConstraint3(summands *SummandList, opType int, rightSide float64,
	penaltyNeg, penaltyPos float64) *Constraint {
	return self.addConstraint(summands, opType, rightSide, penaltyNeg, penaltyPos)
}

func (self *LinearSpec) AddConstraint4(coeffs []float64, vars []*Variable,
	opType int, rightSide float64,
	penaltyNeg, penaltyPos float64) *Constraint {
	if len(coeffs) != len(vars) {
		return nil
	}
	summands := newSummandList()

	for i, c := range coeffs {
        s := NewSummand(c, vars[i])
		summands.AddItem(s)
	}

	return self.addConstraint(summands, opType, rightSide, penaltyNeg, penaltyPos)
}

func (self *LinearSpec) MinSize(width, height *Variable) Size {
	return self.solver.MinSize(width, height)
}

func (self *LinearSpec) MaxSize(width, height *Variable) Size {
	return self.solver.MaxSize(width, height)
}

func (self *LinearSpec) Solve() int {
	// TODO: Measure solve time
	self.result = self.solver.Solve()
	return self.result
}

// Writes the specification into a text file.
// The file will be overwritten if it exists.
func (self *LinearSpec) Save(filename string) bool {
	return self.solver.SaveModel(filename)
}

//func (self *LinearSpec) CountColumns() int {
//}

// Result get the result type
func (self *LinearSpec) Result() int {
	return self.result
}

// SolvingTime gets the solving time.
func (self *LinearSpec) SolvingTime() float64 {
	return self.solvingTime
}

//func (self *LinearSpec) String() string {
//}

// Constraints get the constraints.
func (self *LinearSpec) Constraints() *ConstraintList {
	return self.constraints
}

func (self *LinearSpec) UsedVariables() *VariableList {
	return self.usedVariables
}

func (self *LinearSpec) AllVariables() *VariableList {
	return self.variables
}

func (self *LinearSpec) updateLeftSide(c *Constraint) bool {
	if !self.solver.LeftSideChanged(c) {
		return false
	}
	return true
}

func (self *LinearSpec) updateRightSide(c *Constraint) bool {
	if !self.solver.RightSideChanged(c) {
		return false
	}
	return true
}

func (self *LinearSpec) updateOperator(c *Constraint) bool {
	if !self.solver.OperatorChanged(c) {
		return false
	}
	return true
}

func (self *LinearSpec) checkSummandList(list *SummandList) bool {
	ok := true
	for i := 0; i < list.Len(); i++ {
		s := list.GetAt(i)
		if s == nil {
			ok = false
			break
		}
	}

	if ok {
		return true
	}

	list.Clear()
	return false
}

func (self *LinearSpec) addConstraint(leftSide *SummandList, opType int,
	rightSide, penaltyNeg, penaltyPos float64) *Constraint {

	c := newConstraint(self, leftSide, opType, rightSide, penaltyNeg, penaltyPos)
	if !self.AddConstraint(c) {
		return nil
	}
	return c
}

func (self *LinearSpec) String() string {
	s := ""
	for i := 0; i < self.variables.Len(); i++ {
		variable := self.variables.GetAt(i)
		s = s + fmt.Sprintf("%v=%v ", variable.String(), variable.Value())
	}
	s = s + "\n"
	for i := 0; i < self.constraints.Len(); i++ {
		c := self.constraints.GetAt(i)
		s = s + fmt.Sprintf("%v: %v\n", i, c.String())
	}

	s = s + "Result="
	switch self.Result() {
	case ResultError:
		s = s + "Error"
	case ResultOptimal:
		s = s + "Optimal"
	case ResultSubOptimal:
		s = s + "SubOptimal"
	case ResultInfeasible:
		s = s + "Infeasible"
	case ResultUnbounded:
		s = s + "Unbounded"
	case ResultDegenerate:
		s = s + "Degenerate"
	case ResultNumFailure:
		s = s + "NumFailure"
	default:
		s = s + string(self.Result())
	}

	return s

}
