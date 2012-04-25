package lp

type Summand struct {
	coeff float64
	v     *Variable
}

func NewSummand(coeff float64, v *Variable) *Summand {
	s := &Summand{}
	s.coeff = coeff
	s.v = v

	return s
}

func (self *Summand) Coeff() float64 {
	return self.coeff
}

func (self *Summand) SetCoeff(coeff float64) {
	self.coeff = coeff
}

func (self *Summand) Var() *Variable {
	return self.v
}

func (self *Summand) SetVar(v *Variable) {
	self.v = v
}

func (self *Summand) VariableIndex() int {
	return self.v.Index()
}
