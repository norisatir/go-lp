package lp

type ConstraintList []*Constraint

func (self *ConstraintList) AddItem(c *Constraint) {
	*self = append(*self, c)
}

func (self *ConstraintList) RemoveItem(c *Constraint) bool {
	for i, c1 := range *self {
		if c == c1 {
			*self = append((*self)[:i], (*self)[i+1:]...)
			return true
		}
	}
	return false
}

func (self *ConstraintList) RemoveItemAt(i int) bool {
    if i >= len(*self) {
        return false
    }
    *self = append((*self)[:i], (*self)[i+1:]...)
    return true
}

func (self ConstraintList) IndexOf(c *Constraint) int {
	for i, c1 := range self {
		if c1 == c {
			return i
		}
	}
	return -1
}

type VariableList []*Variable

func (self *VariableList) AddItem(v *Variable) {
	*self = append(*self, v)
}

func (self *VariableList) RemoveItem(v *Variable) bool {
	for i, v1 := range *self {
		if v == v1 {
			*self = append((*self)[:i], (*self)[i+1:]...)
			return true
		}
	}
	return false
}

func (self VariableList) IndexOf(v *Variable) int {
	for i, v1 := range self {
		if v1 == v {
			return i
		}
	}
	return -1
}

type LinearSpec struct {
	variables     VariableList
	usedVariables VariableList
	constraints   ConstraintList
	result        int
	solvingTime   float64
	solver        SolverLike
}

func NewLinearSpec() *LinearSpec {
	ls := &LinearSpec{}
	ls.result = ResultError
	ls.solvingTime = 0
	// TODO: new activeSolver
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
	var markedForInvalidation ConstraintList = make(ConstraintList, 0)
	constraints := self.Constraints()

	for _, c := range constraints {
		if !c.IsValid() {
			continue
		}

		summands := c.LeftSide()
		for _, s := range summands {
			if s.Var() == v {
				markedForInvalidation.AddItem(c)
				break
			}
		}
	}

	for _, m := range markedForInvalidation {
		self.RemoveConstraint(m)
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

	for _, s := range c.LeftSide() {
		self.usedVariables.AddItem(s.Var())
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

	if !self.constraints.RemoveItem(c) {
		return false
	}
	c.isValid = false

	for _, s := range c.LeftSide() {
		self.usedVariables.RemoveItem(s.Var())
	}

	return true
}

func (self *LinearSpec) AddConstraint1(summands SummandList, opType int, rightSide float64) *Constraint {
	return self.AddConstraint3(summands, opType, rightSide, -1, -1)
}

func (self *LinearSpec) AddConstraint2(coeffs []float64, vars []*Variable,
	opType int, rightSide float64) *Constraint {
	return self.AddConstraint4(coeffs, vars, opType, rightSide, -1, -1)
}

func (self *LinearSpec) AddConstraint3(summands SummandList, opType int, rightSide float64,
	penaltyNeg, penaltyPos float64) *Constraint {
	return self.addConstraint(summands, opType, rightSide, penaltyNeg, penaltyPos)
}

func (self *LinearSpec) AddConstraint4(coeffs []float64, vars []*Variable,
	opType int, rightSide float64,
	penaltyNeg, penaltyPos float64) *Constraint {
	if len(coeffs) != len(vars) {
		return nil
	}
	summands := make(SummandList, 0)

	for i, c := range coeffs {
		summands.AddItem(NewSummand(c, vars[i]))
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
func (self *LinearSpec) Constraints() ConstraintList {
	return self.constraints
}

func (self *LinearSpec) UsedVariables() VariableList {
	return self.usedVariables
}

func (self *LinearSpec) AllVariables() VariableList {
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

func (self *LinearSpec) checkSummandList(list SummandList) bool {
	ok := true
	for _, s := range list {
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

func (self *LinearSpec) addConstraint(leftSide SummandList, opType int,
	rightSide, penaltyNeg, penaltyPos float64) *Constraint {

	c := newConstraint(self, leftSide, opType, rightSide, penaltyNeg, penaltyPos)
	if !self.AddConstraint(c) {
		return nil
	}
	return c
}
