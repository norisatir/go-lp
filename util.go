package lp

type ConstraintList struct {
	vec []*Constraint
}

func newConstraintList() *ConstraintList {
	cl := &ConstraintList{}
	cl.vec = make([]*Constraint, 0)

	return cl
}

func (self *ConstraintList) AddItem(c *Constraint) {
	self.vec = append(self.vec, c)
}

func (self *ConstraintList) Len() int {
	return len(self.vec)
}

func (self *ConstraintList) AddItemAt(c *Constraint, index int) {
	self.vec = append(self.vec, nil)
	oldConstraint := c
	for i := index; i < self.Len(); i++ {
		oldConstraint, self.vec[i] = self.vec[i], oldConstraint
	}
}

func (self *ConstraintList) RemoveItem(c *Constraint) bool {
	i := self.IndexOf(c)
	if i != -1 {
		self.RemoveItemAt(i)
		return true
	}
	return false
}

func (self *ConstraintList) RemoveItemAt(i int) bool {
	if i >= self.Len() {
		return false
	}
	self.vec = append(self.vec[:i], self.vec[i+1:]...)
	return true
}

func (self *ConstraintList) IndexOf(c *Constraint) int {
	for i, c1 := range self.vec {
		if c1 == c {
			return i
		}
	}
	return -1
}

func (self *ConstraintList) GetAt(index int) *Constraint {
	if index >= self.Len() {
		return nil
	}
	return self.vec[index]
}

func (self *ConstraintList) Clear() {
	self.vec = self.vec[:0]
}

type VariableList struct {
	vec []*Variable
}

func newVariableList() *VariableList {
	vl := &VariableList{}
	vl.vec = make([]*Variable, 0)

	return vl
}

func (self *VariableList) AddItem(v *Variable) {
	self.vec = append(self.vec, v)
}

func (self *VariableList) Len() int {
	return len(self.vec)
}

func (self *VariableList) AddItemAt(v *Variable, index int) {
	self.vec = append(self.vec, nil)
	oldVar := v
	for i := index; i < self.Len(); i++ {
		oldVar, self.vec[i] = self.vec[i], oldVar
	}
}

func (self *VariableList) RemoveItem(v *Variable) bool {
	i := self.IndexOf(v)
	if i != -1 {
		self.RemoveItemAt(i)
		return true
	}
	return false
}

func (self *VariableList) RemoveItemAt(i int) bool {
	if i >= self.Len() {
		return false
	}
	self.vec = append(self.vec[:i], self.vec[i+1:]...)
	return true
}

func (self *VariableList) IndexOf(v *Variable) int {
	for i, v1 := range self.vec {
		if v1 == v {
			return i
		}
	}
	return -1
}

func (self *VariableList) GetAt(index int) *Variable {
	if index >= self.Len() {
		return nil
	}
	return self.vec[index]
}

func (self *VariableList) Clear() {
	self.vec = self.vec[:0]
}

type SummandList struct {
	vec []*Summand
}

func newSummandList() *SummandList {
	sl := &SummandList{}
	sl.vec = make([]*Summand, 0)

	return sl
}

func (self *SummandList) AddItem(s *Summand) {
	self.vec = append(self.vec, s)
}

func (self *SummandList) Len() int {
	return len(self.vec)
}

func (self *SummandList) AddItemAt(s *Summand, index int) {
	self.vec = append(self.vec, nil)
	oldSummand := s
	for i := index; i < self.Len(); i++ {
		oldSummand, self.vec[i] = self.vec[i], oldSummand
	}
}

func (self *SummandList) RemoveItem(s *Summand) bool {
	i := self.IndexOf(s)
	if i != -1 {
		self.RemoveItemAt(i)
		return true
	}
	return false
}

func (self *SummandList) RemoveItemAt(i int) bool {
	if i >= self.Len() {
		return false
	}
	self.vec = append(self.vec[:i], self.vec[i+1:]...)
	return true
}

func (self *SummandList) IndexOf(s *Summand) int {
	for i, s1 := range self.vec {
		if s1 == s {
			return i
		}
	}
	return -1
}

func (self *SummandList) GetAt(index int) *Summand {
	if index >= self.Len() {
		return nil
	}
	return self.vec[index]
}

func (self *SummandList) Clear() {
	self.vec = self.vec[:0]
}

type softInEqList struct {
	vec []*softInEqData
}

func newsoftInEqList() *softInEqList {
	sl := &softInEqList{}
	sl.vec = make([]*softInEqData, 0)

	return sl
}

func (self *softInEqList) AddItem(s *softInEqData) {
	self.vec = append(self.vec, s)
}

func (self *softInEqList) Len() int {
	return len(self.vec)
}

func (self *softInEqList) AddItemAt(s *softInEqData, index int) {
	self.vec = append(self.vec, nil)
	oldsoftInEqData := s
	for i := index; i < self.Len(); i++ {
		oldsoftInEqData, self.vec[i] = self.vec[i], oldsoftInEqData
	}
}

func (self *softInEqList) RemoveItem(s *softInEqData) bool {
	i := self.IndexOf(s)
	if i != -1 {
		self.RemoveItemAt(i)
		return true
	}
	return false
}

func (self *softInEqList) RemoveItemAt(i int) bool {
	if i >= self.Len() {
		return false
	}
	self.vec = append(self.vec[:i], self.vec[i+1:]...)
	return true
}

func (self *softInEqList) IndexOf(s *softInEqData) int {
	for i, s1 := range self.vec {
		if s1 == s {
			return i
		}
	}
	return -1
}

func (self *softInEqList) GetAt(index int) *softInEqData {
	if index >= self.Len() {
		return nil
	}
	return self.vec[index]
}

func (self *softInEqList) Clear() {
	self.vec = self.vec[:0]
}
