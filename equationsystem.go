package lp

import "fmt"

type EquationSystem struct {
    rowIndices []int
    columnIndices []int
    matrix [][]float64
    b []float64
    rows int
    columns int
}

func NewEquationSystem(rows, columns int) *EquationSystem {
    es := &EquationSystem{}
    es.rows = rows
    es.columns = columns

    es.matrix = initMatrixSlice(es.rows, es.columns)
    es.b = make([]float64, es.columns)

    for i := 0; i < es.columns; i++ {
        es.b[i] = 0
    }
    zeroMatrix(es.matrix, es.rows, es.columns)

    es.rowIndices = make([]int, es.rows)
    es.columnIndices = make([]int, es.columns)

    for i := 0; i < es.rows; i++ {
        es.rowIndices[i] = i
    }
    for i := 0; i < es.columns; i++ {
        es.columnIndices[i] = i
    }

    return es
}

func (self *EquationSystem) SetRows(rows int) {
    self.rows = rows
}

func (self *EquationSystem) Rows() int {
    return self.rows
}

func (self *EquationSystem) Columns() int {
    return self.columns
}

func (self *EquationSystem) A(row, column int) *float64 {
    return &self.matrix[self.rowIndices[row]][self.columnIndices[column]]
}

func (self *EquationSystem) B(row int) *float64 {
    return &self.b[row]
}

func (self *EquationSystem) Results(results []float64, size int) {
    for i := 0; i < size; i++ {
        results[i] = 0
    }
    for i := 0; i < self.columns; i++ {
        index := self.columnIndices[i]
        if index < self.rows {
            results[index] = self.b[i]
        }
    }
}

func (self *EquationSystem) SwapColumn(i, j int) {
    self.columnIndices[i], self.columnIndices[j] = 
    self.columnIndices[j], self.columnIndices[i]
}

func (self *EquationSystem) SwapRow(i, j int) {
    self.rowIndices[i], self.rowIndices[j] =
    self.rowIndices[j], self.rowIndices[i]

    self.b[i], self.b[j] = 
    self.b[j], self.b[i]
}

func (self *EquationSystem) GaussianElimination() bool {
    // basic solve
    for i := 0; i < self.rows; i++ {
        // find none zero element
        swapRow := -1
        for r := i; r < self.rows; r++ {
            value := self.matrix[self.rowIndices[r]][self.columnIndices[i]]
            if fuzzyEquals(value, 0) { continue }
            swapRow = r
            break
        }

        if swapRow == -1 {
            swapColumn := -1
            for c := i+1; c < self.columns; c++ {
                value := self.matrix[self.rowIndices[i]][self.columnIndices[c]]
                if fuzzyEquals(value, 0) { continue }
                swapRow = i
                swapColumn = c
                break
            }

            if swapColumn == -1 {
                return false
            }
            self.SwapColumn(i, swapColumn)
        }
        if i != swapRow {
            self.SwapRow(i, swapRow)
        }

        // normalize
        self.eliminateColumn(i, i+1, self.rows - 1)
    }

    return true
}

func (self *EquationSystem) GaussJordan() bool {
    if !self.GaussianElimination() {
        return false
    }

    for i := self.rows-1; i >= 0; i-- {
        self.eliminateColumn(i, 0, i-1)
    }
    return true
}

func (self *EquationSystem) GaussJordan1(i int) {
    self.eliminateColumn(i, 0, self.rows-1)
}

func (self *EquationSystem) RemoveLinearlyDependentRows() {
    oldB := make([]float64, self.rows)
    for r := 0; r < self.rows; r++ {
        oldB[r] = self.b[r]
    }

    temp := initMatrixSlice(self.rows, self.columns)
    independentRows := make([]bool, self.rows)

    // copy to temp
    copyMatrix(self.matrix, temp, self.rows, self.columns)
    nIndependent := computeDependencies(temp, self.rows, self.columns, independentRows)
    if nIndependent == self.rows { return }

    // remove the rows
    for i := 0; i < self.rows; i++ {
        if !independentRows[i] {
            lastDepRow := -1
            for d := self.rows-1; d < i; d-- {
                if independentRows[d] {
                    lastDepRow = d
                    break
                }
            }
            if lastDepRow < 0 { break }
            self.SwapRow(i, lastDepRow)
            self.rows--
        }
    }
    self.rows = nIndependent
}

func (self *EquationSystem) RemoveUnusedVariables() {
    for c := 0; c < self.columns; c++ {
        used := false
        for r := 0; r < self.rows; r++ {
            if !fuzzyEquals(self.matrix[r][self.columnIndices[c]], 0) {
                used = true
                break
            }
        }
        if used { continue }

        self.SwapColumn(c, self.columns-1)
        self.columns--
        c--
    }
}

func (self *EquationSystem) MoveColumnRight(i, target int) {
    index := self.columnIndices[i]
    for c := i; c < target; c++ {
        self.columnIndices[c] = self.columnIndices[c+1]
    }
    self.columnIndices[target] = index
}

func (self *EquationSystem) Print() {
    for m := 0; m < self.rows; m++ {
        for n := 0; n < self.columns; n++ {
            fmt.Printf("%.1f ", self.matrix[self.rowIndices[m]][self.columnIndices[n]])
        }
        fmt.Printf("= %.1f\n", self.b[m])
    }
}

func (self *EquationSystem) eliminateColumn(column, startRow, endRow int) {
    value := self.matrix[self.rowIndices[column]][self.columnIndices[column]]
    if value != 1.0 {
        for j := column; j < self.columns; j++ {
            self.matrix[self.rowIndices[column]][self.columnIndices[j]] /= value
        }
        self.b[column] /= value
    }

    for r := startRow; r < endRow+1; r++ {
        if r == column { continue }
        q := -self.matrix[self.rowIndices[r]][self.columnIndices[column]]
        // don't need to to anything, since matrix is typically sparse
        // this should save some work
        if fuzzyEquals(q, 0) { continue }
        for c := column; c < self.columns; c++ {
            self.matrix[self.rowIndices[r]][self.columnIndices[c]] += 
            self.matrix[self.rowIndices[column]][self.columnIndices[c]] * q
        }
        self.b[r] += self.b[column] * q
    }
}
