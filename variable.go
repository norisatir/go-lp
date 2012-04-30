package lp

import "math"
import "strconv"

// Variable struct contains minimum and maximum values
type Variable struct {
	ls              *LinearSpec
	min, max, value float64
	label           string
	isValid         bool
}

func newVariable(ls *LinearSpec) *Variable {
	v := &Variable{}
	v.ls = ls
	v.label = ""
	v.value = math.NaN()
	v.min = -20000
	v.max = 20000
	v.isValid = false

	return v
}

// Index returns the index of the variable
func (self *Variable) Index() int {
	return self.ls.IndexOf(self)
}

func (self *Variable) GlobalIndex() int {
	return self.ls.GlobalIndexOf(self)
}

// LS gets the current linear specification
func (self *Variable) LS() *LinearSpec {
	return self.ls
}

// Value gets the value
func (self *Variable) Value() float64 {
	return self.value
}

// SetValue sets the value
func (self *Variable) SetValue(val float64) {
	self.value = val
}

// Min gets the minimum value of the variable
func (self *Variable) Min() float64 {
	return self.min
}

// SetMin sets the minimum value of the variable
func (self *Variable) SetMin(min float64) {
	self.SetRange(min, self.max)
}

// Max gets the maximum value of the variable
func (self *Variable) Max() float64 {
	return self.max
}

// SetMax sets the maximum value of the variable
func (self *Variable) SetMax(max float64) {
	self.SetRange(self.min, max)
}

// SetRange sets the minimum and maximum values of the variable
func (self *Variable) SetRange(min, max float64) {
	if !self.isValid {
		return
	}

	self.min = min
	self.max = max
	self.ls.UpdateRange(self)
}

// Label returns Variable label
func (self *Variable) Label() string {
	return self.label
}

// SetLabel sets variable label
func (self *Variable) SetLabel(label string) {
	self.label = label
}

// Returns index of variable as String
func (self *Variable) String() string {
	resStr := ""

	if self.label != "" {
		resStr = resStr + self.label
		if !self.isValid {
			resStr = resStr + " (invalid)"
		}
	} else {
		resStr = resStr + " Variable"
		if !self.isValid {
			resStr = resStr + "(invalid)"
		} else {
			resStr = resStr + strconv.Itoa(self.Index())
		}
	}

	return resStr
}

func (self *Variable) IsEqual(v *Variable) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint2([]float64{1.0, -1.0}, []*Variable{self, v}, OperatorEQ, 0.0)
}

func (self *Variable) IsSmallerOrEqual(v *Variable) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint2([]float64{1.0, -1.0}, []*Variable{self, v}, OperatorLE, 0.0)
}

func (self *Variable) IsGreaterOrEqual(v *Variable) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint2([]float64{-1.0, 1.0}, []*Variable{v, self}, OperatorGE, 0.0)
}

func (self *Variable) IsEqual1(v *Variable, penaltyNeg, penaltyPos float64) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint4([]float64{1.0, -1.0}, []*Variable{self, v}, OperatorEQ, 0.0,
		penaltyNeg, penaltyPos)
}

func (self *Variable) IsSmallerOrEqual1(v *Variable, penaltyNeg, penaltyPos float64) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint4([]float64{1.0, -1.0}, []*Variable{self, v}, OperatorLE, 0.0,
		penaltyNeg, penaltyPos)
}

func (self *Variable) IsGreaterOrEqual1(v *Variable, penaltyNeg, penaltyPos float64) *Constraint {
	if !self.isValid {
		return nil
	}
	return self.ls.AddConstraint4([]float64{-1.0, 1.0}, []*Variable{v, self}, OperatorGE, 0.0,
		penaltyNeg, penaltyPos)
}

func (self *Variable) IsValid() bool {
	return self.isValid
}

func (self *Variable) Invalidate() {
	if !self.isValid {
		return
	}
	self.isValid = false
	self.ls.RemoveVariable(self)
}

///////////////////////////////////////////////////////////
