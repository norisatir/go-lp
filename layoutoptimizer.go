package lp

import "errors"

type LayoutOptimizer struct {
    variableCount int
    constraints ConstraintList

    temp1 [][]float64
    temp2 [][]float64
    zTrans [][]float64
    q [][]float64
    activeMatrix [][]float64
    activeMatrixTemp [][]float64
    softConstraints [][]float64
    g [][]float64
    desired []float64
}

func NewLayoutOptimizer(list ConstraintList, variableCount int) *LayoutOptimizer {
    lo := &LayoutOptimizer{}
    lo.SetConstraints(list, variableCount)

    return lo
}

func (self *LayoutOptimizer) SetConstraints(list ConstraintList, variableCount int) bool {
    self.constraints = list
    constraintCount := len(self.constraints)

    if self.variableCount != variableCount {
        self.makeEmpty()
        self.init(variableCount, constraintCount)
    }

    zeroMatrix(self.softConstraints, constraintCount, self.variableCount)
    rightSide := make([]float64, constraintCount)

    // set up soft constraint matrix
    for c := 0; c < len(self.constraints); c++ {
        constraint := self.constraints[c]
        if !constraint.IsSoft() {
            rightSide[c] = 0.0
            continue
        }
        weight := 0.0
        negPenalty := constraint.PenaltyNeg()
        if negPenalty > 0 {
            weight += negPenalty
        }
        posPenalty := constraint.PenaltyPos()
        if posPenalty > 0 {
            weight += posPenalty
        }
        if negPenalty > 0 && posPenalty > 0 {
            weight /= 2
        }

        rightSide[c] = self.rightSide(constraint) * weight
        summands := constraint.LeftSide()
        for s := 0; s < len(summands); s++ {
            summand := summands[s]
            variable := summand.Var().Index()
            if constraint.Op() == OperatorLE {
                self.softConstraints[c][variable] = -summand.Coeff()
            } else {
                self.softConstraints[c][variable] = summand.Coeff()
            }
            self.softConstraints[c][variable] *= weight
        }
    }

    // create g
    transposeMatrix(self.softConstraints, self.temp1, constraintCount, self.variableCount)
    multiplyMatrices(self.temp1, self.softConstraints, self.g, self.variableCount,
        constraintCount, self.variableCount)

    // create d
    multiplyMatrixVector(self.temp1, rightSide, self.variableCount, constraintCount, self.desired)
    negateVector(self.desired, self.variableCount)

    return true
}

func (self *LayoutOptimizer) InitCheck() error {
    if self.temp1 == nil || self.temp2 == nil || self.zTrans == nil || self.q == nil || 
        self.softConstraints == nil || self.g == nil || self.desired == nil {
        
            return errors.New("No memory")
    }
    return nil
}

func (self *LayoutOptimizer) actualValue(constraint *Constraint, values []float64) float64 {
    summands := constraint.LeftSide()
    value := 0.0
    for s := 0; s < len(summands); s++ {
        summand := summands[s]
        variable := summand.Var().Index()
        value += values[variable] * summand.Coeff()
    }

    if constraint.Op() == OperatorLE {
        return -value
    }
    return value
}

func (self *LayoutOptimizer) rightSide(c *Constraint) float64 {
    if c.Op() == OperatorLE {
        return -c.RightSide()
    }
    return c.RightSide()
}

func (self *LayoutOptimizer) makeEmpty() {
    self.temp1 = nil
    self.temp2 = nil
    self.zTrans = nil
    self.softConstraints = nil
    self.q = nil
    self.g = nil
    self.desired = nil
}

func (self *LayoutOptimizer) init(variableCount, nConstraints int) {
    self.variableCount = variableCount
    
    maxExtend := max(variableCount, nConstraints)
    self.temp1 = initMatrixSlice(maxExtend, maxExtend)
    self.temp2 = initMatrixSlice(maxExtend, maxExtend)
    self.zTrans = initMatrixSlice(nConstraints, self.variableCount)
    self.softConstraints = initMatrixSlice(nConstraints, self.variableCount)
    self.q = initMatrixSlice(nConstraints, self.variableCount)
    self.g = initMatrixSlice(nConstraints, nConstraints)

    self.desired = make([]float64, self.variableCount)
}

func (self *LayoutOptimizer) solve(values []float64) bool {
    if values == nil {
        return false
    }

    constraintCount := len(self.constraints)

	// our QP is supposed to be in this form:
	//   min_x 1/2x^TGx + x^Td
	//   s.t. a_i^Tx = b_i,  i \in E
	//        a_i^Tx >= b_i, i \in I

	// init our initial x
    x := make([]float64, self.variableCount)
    for i := 0; i < self.variableCount; i++ {
        x[i] = values[i]
    }

	// init d
	// Note that the values of d and of g result from rewriting the
	// ||x - desired|| we actually want to minimize.
    d := make([]float64, self.variableCount)
    for i := 0; i < self.variableCount; i++ {
        d[i] = self.desired[i]
    }

    // init active set
    activeConstraints := make(ConstraintList, 0)

    for i := 0; i < constraintCount; i++ {
        constraint := self.constraints[i]
        if constraint.IsSoft() { continue }

        actualValue := self.actualValue(constraint, x)

        if fuzzyEquals(actualValue, self.rightSide(constraint)) {
            activeConstraints.AddItem(constraint)
        }
    }

	// The main loop: Each iteration we try to get closer to the optimum
	// solution. We compute a vector p that brings our x closer to the optimum.
	// We do that by computing the QP resulting from our active constraint set,
	// W^k. Afterward each iteration we adjust the active set.

    for true {
		// solve the QP:
		//   min_p 1/2p^TGp + g_k^Tp
		//   s.t. a_i^Tp = 0
		//   with a_i \in activeConstraints
		//        g_k = Gx_k + d
		//        p = x - x_k

        activeCount := len(activeConstraints)
        if activeCount == 0 {
            return false
        }

        // construct a matrix from the active constraints
        am := activeCount
        an := self.variableCount
        independentRows := make([]bool, activeCount)
        zeroMatrix(self.activeMatrix, am, an)

        for i := 0; i < activeCount; i++ {
            constraint := activeConstraints[i]
            if constraint.IsSoft() { continue }

            summands := constraint.LeftSide()
            for s := 0; s < len(summands); s++ {
                summand := summands[i]
                variable := summand.Var().Index()
                if constraint.Op() == OperatorLE {
                    self.activeMatrix[i][variable] = -summand.Coeff()
                } else {
                    self.activeMatrix[i][variable] = summand.Coeff()
                }
            }
        }

        am = removeLinearyDependentRows(self.activeMatrix, self.activeMatrixTemp, 
                independentRows, am, an)

        // gxd = G * x + d
        gxd := make([]float64, self.variableCount)
        multiplyMatrixVector(self.g, x, self.variableCount, self.variableCount, gxd)
        addVectors(gxd, d, self.variableCount)

        p := make([]float64, self.variableCount)

        if !self.solveSubProblem(gxd, am, p) {
            return false
        }

        if isZero(p, self.variableCount) {
			// compute Lagrange multipliers lambda_i
			// if lambda_i >= 0 for all i \in W^k \union inequality constraints,
			// then we're done.
			// Otherwise remove the constraint with the smallest lambda_i
			// from the active set.
			// The derivation of the Lagrangian yields:
			//   \sum_{i \in W^k}(lambda_ia_i) = Gx_k + d
			// Which is an system we can solve:
			//   A^Tlambda = Gx_k + d

			// A^T is over-determined, hence we need to reduce the number of
			// rows before we can solve it.

            independentColumns := make([]bool, an)
            aa := self.temp1
            transposeMatrix(self.activeMatrix, aa, am, an)
            aam := removeLinearyDependentRows(aa, self.temp2, independentColumns, an, am)
            aan := am

            if aam == aan {
                return false
            }

            // also reduce the number of rows on the right hand side
            lambda := make([]float64, aam)
            index := 0
            for i := 0; i < an; i++ {
                if independentColumns[i] {
                    lambda[index] = gxd[i]
                    index++
                }
            }

            success := solve(aa, am, lambda)
            if !success {
                return false
            }

            // find min lambda_i (only, if it's < 0, though)
            minLambda := 0.0
            minIndex := -1
            index = 0
            for i := 0; i < activeCount; i++ {
                if independentRows[i] {
                    constraint := activeConstraints[i]
                    if constraint.Op() != OperatorEQ {
                        if lambda[index] < minLambda {
                            minLambda = lambda[index]
                            minIndex = i
                        }
                    }

                    index++
                }
            }


            // if the min lambda is >= 0, we're done
            if minIndex < 0 || fuzzyEquals(minLambda, 0) {
                self.setResult(x, values)
                return true
            }

            // remove i from the active set
            activeConstraints.RemoveItemAt(minIndex)
        } else {
            // compute alpha_k
            alpha := 0.0
            barrier := -1
            // if alpha_k < 1, add a barrier constraint to W^k
            for i := 0; i < constraintCount; i++ {
                constraint := self.constraints[i]
                if activeConstraints.IndexOf(constraint) != -1 {
                    continue
                }

                divider := self.actualValue(constraint, p)
                if divider > 0 || fuzzyEquals(divider, 0) {
                    continue
                }

                // (b_i - a_i^Tx_k) / a_i^Tp_k
                alphaI := self.rightSide(constraint) - self.actualValue(constraint, x)
                alphaI /= divider
                if alphaI < alpha {
                    alpha = alphaI
                    barrier = i
                }
            }

            if alpha < 1 {
                activeConstraints.AddItem(self.constraints[barrier])
            }

            // x += p * alpha
            addVectorsScaled(x, p, alpha, self.variableCount)
        }
    }
    return false
}

func (self *LayoutOptimizer) solveSubProblem(d []float64, am int, p []float64) bool {
}

/*	Solve solves the quadratic program (QP) given by the constraints added via
	AddConstraint(), the additional constraint \sum_{i=0}^{n-1} x_i = size,
	and the optimization criterion to minimize
	\sum_{i=0}^{n-1} (x_i - desired[i])^2.
	The \a values slice must contain a feasible solution when called and will
	be overwritten with the optimial solution the method computes.
*/
func (self *LayoutOptimizer) Solve(values []float64) (success bool) {
    if values == nil {
        return false
    }

    constraintCount := len(self.constraints)

    // allocate the active constraint matrix and its transposed matrix
    self.activeMatrix = initMatrixSlice(constraintCount, self.variableCount)
    self.activeMatrixTemp = initMatrixSlice(constraintCount, self.variableCount)

    success = self.solve(values)
    return
}

func (self *LayoutOptimizer) setResult(x, values []float64) {
}
