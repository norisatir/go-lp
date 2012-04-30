package lp

// Hard linear constraint, i.e. one that must be satisfied.
// May render a specification infeasible.
type Constraint struct {
	ls                     *LinearSpec
	leftSide               *SummandList
	opType                 int
	rightSide              float64
	penaltyNeg, penaltyPos float64
	dNegObjSummand         *Summand
	dPosObjSummand         *Summand
	label                  string
	isValid                bool
}

func newConstraint(ls *LinearSpec, summands *SummandList, opType int,
	rightSide float64, penaltyNeg, penaltyPos float64) *Constraint {

	c := &Constraint{}
	c.ls = ls
	c.opType = opType
	c.rightSide = rightSide
	c.penaltyNeg = penaltyNeg
	c.penaltyPos = penaltyPos
	c.dNegObjSummand = nil
	c.dPosObjSummand = nil
	c.isValid = true
	c.SetLeftSide(summands)

	return c
}

// Index gets the index of the constraint
func (self *Constraint) Index() int {
	return self.ls.Constraints().IndexOf(self)
}

// LeftSide gets the left side of the constraint.
func (self *Constraint) LeftSide() *SummandList {
	return self.leftSide
}

// SetLeftSide sets the summands on the left side of the constraint
func (self *Constraint) SetLeftSide(summands *SummandList) {
	if !self.isValid {
		return
	}

	// check left side
	for i := 0; i < summands.Len(); i++ {
		s := summands.GetAt(i)
		for a := i + 1; a < summands.Len(); a++ {
			nextSummand := summands.GetAt(a)
			if s.Var() == nextSummand.Var() {
				s.SetCoeff(s.Coeff() + nextSummand.Coeff())
				summands.RemoveItem(nextSummand)
				a--
			}
		}
	}
	self.leftSide = summands
	self.ls.updateLeftSide(self)
}

func (self *Constraint) SetLeftSide1(coeffs []float64, vars []*Variable) {
	if !self.isValid {
		return
	}
	if len(coeffs) != len(vars) {
		return
	}

	self.LeftSide().Clear()
	for i, c := range coeffs {
		self.leftSide.AddItem(NewSummand(c, vars[i]))
	}
	self.SetLeftSide(self.leftSide)
}

// Op gets the operator used for this constraint.
func (self *Constraint) Op() int {
	return self.opType
}

// SetOp sets the operator used for this constraint.
func (self *Constraint) SetOp(opType int) {
	if !self.isValid {
		return
	}
	self.opType = opType
	self.ls.updateOperator(self)
}

// RightSide gets the constant value that is on the right side of the operator.
func (self *Constraint) RightSide() float64 {
	return self.rightSide
}

// SetRightSide sets the constant value that is on the right side of the operator.
func (self *Constraint) SetRightSide(value float64) {
	if !self.isValid {
		return
	}
	if self.rightSide == value {
		return
	}
	self.ls.updateRightSide(self)
}

// PenaltyNeg gets the penalty coefficient for negative deviations.
func (self *Constraint) PenaltyNeg() float64 {
	return self.penaltyNeg
}

// SetPenaltyNeg sets the penalty coefficient for negative deviations from the soft
// constraint's exact solution, i.e. if the left side is too large.
func (self *Constraint) SetPenaltyNeg(value float64) {
	self.penaltyNeg = value
	self.ls.updateLeftSide(self)
}

// PenaltyPos gets the penalty coefficient for positive deviations.
func (self *Constraint) PenaltyPos() float64 {
	return self.penaltyPos
}

// SetPenaltyPos sets the penalty coefficient for negative deviations from the soft
// constraint's exact solution, i.e. if the left side is too small.
func (self *Constraint) SetPenaltyPos(value float64) {
	self.penaltyPos = value
	self.ls.updateLeftSide(self)
}

func (self *Constraint) Label() string {
	return self.label
}

func (self *Constraint) SetLabel(label string) {
	self.label = label
}

// DNeg gets the slack variable for the negative variations.
func (self *Constraint) DNeg() *Variable {
	if self.dNegObjSummand == nil {
		return nil
	}
	return self.dNegObjSummand.Var()
}

// DPos gets the slack variable for the positive variations.
func (self *Constraint) DPos() *Variable {
	if self.dPosObjSummand == nil {
		return nil
	}
	return self.dPosObjSummand.Var()
}

func (self *Constraint) IsSoft() bool {
	if self.opType == OperatorEQ {
		return false
	}
	if self.penaltyNeg > 0.0 || self.penaltyPos > 0.0 {
		return true
	}
	return false
}

func (self *Constraint) IsValid() bool {
	return self.isValid
}

func (self *Constraint) Invalidate() {
	if !self.isValid {
		return
	}
	self.isValid = false
	self.ls.RemoveConstraint(self)
}

func (self *Constraint) String() string {
	s := "Constraint "
	s = s + self.label
	// TODO: more info

	return s
}
