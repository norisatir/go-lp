package lp

import "math"

const (
    EqualsEpsilon float64 = 0.000001
)

func fuzzyEquals(a, b float64) bool {
    return math.Abs(a-b) < EqualsEpsilon
}

type Matrix struct {
    a [][]float64
    Rows int
    Cols int
}

func NewMatrix(m, n int) *Matrix {
    mat := &Matrix{}
    mat.Rows = m
    mat.Cols = n
    mat.a = make([][]float64, m)

    values := make([]float64, m*n)

    for i,j := 0,0; i < m; i,j = i+1,j+n {
        mat.a[i] = values[j:j+n]
    }

    return mat
}
