package lp

import "math"

const (
    EqualsEpsilon float64 = 0.000001
)

func fuzzyEquals(a, b float64) bool {
    return math.Abs(a-b) < EqualsEpsilon
}

func initMatrixSlice(m, n int) (A [][]float64) {
    A = make([][]float64, m)
    values := make([]float64, m*n)

    for i,j := 0,0; i < m; i,j = i+1,j+n {
        A[i] = values[j:j+n]
    }
    return
}


// multiplyMatrixVector multiplies matrix with vector
//         y = Ax
func multiplyMatrixVector(A [][]float64, vec []float64, m,n int) (y []float64) {
    sum := 0.0
    y = make([]float64, m)
    for i := 0; i < m; i++ {
        for k := 0; k < n; k++ {
            sum = sum + A[i][k] * vec[k]
        }
        y[i] = sum
        sum = 0.0
    }

    return
}

func multiplyMatrices(A, B [][]float64, m, n, l int) (C [][]float64) {
    sum := 0.0

    C = initMatrixSlice(m, l)
    for i := 0; i < m; i++ {
        for j := 0; j < l; j++ {
            for k := 0; k < n; k++ {
                sum = sum + A[i][k] * B[k][j]
            }
            C[i][j] = sum
            sum = 0
        }
    }
    return
}

func transposeMatrix(A [][]float64, m, n int) (B [][]float64) {
    B = initMatrixSlice(n, m)
    
    for i := 0; i < m; i++ {
        for k := 0; k < n; k++ {
            B[k][i] = A[i][k]
        }
    }
    return
}

func zeroMatrix(A [][]float64, m, n int) {
    for i := 0; i < m; i++ {
        for k := 0; k < n; k++ {
            A[i][k] = 0
        }
    }
}

func copyMatrix(A [][]float64, m, n int) (B [][]float64) {
    B = initMatrixSlice(m, n)

    for i := 0; i < m; i++ {
        for k := 0; k < n; k++ {
            B[i][k] = A[i][k]
        }
    }
    return
}

func multiplyOptimizationMatrixVector(x []float64, n int) (y []float64) {
    y = make([]float64, n)
    if n == 1 {
        y[0] = x[0]
        return
    }

    y[0] = 2 * x[0] - x[1]
    for i := 1; i < n-1; i++ {
        y[i] = 2 * x[i] - x[i - 1] - x[i + 1]
    }

    y[n-1] = x[n - 1] - x[n - 2]
    return
}

func multiplyOptimizationMatrixMatrix(A [][]float64, m, n int) (B [][]float64) {
    B = initMatrixSlice(m, n)

    if m == 1 {
        copy(B[0], A[0])
    }

    for k := 0; k < n; k++ {
        B[0][k] = 2 * A[0][k] - A[1][k]
        for i := 1; i < m-1; i++ {
            B[i][k] = 2 * A[i][k] - A[i-1][k] - A[i+1][k]
        }
        B[m-1][k] = A[m-1][k] - A[m-2][k]
    }
    return
}

func solve(a [][]float64, n int) (bool, []float64) {
    b := make([]float64, n)

    // index slice for row permutation
    indices := make([]int, n)
    for i := 0; i < n; i++ {
        indices[i] = i
    }

    // forward elimination
    for i := 0; i < n-1; i++ {
        // find pivot
        pivot := i
        pivotValue := math.Abs(a[indices[i]][i])
        for j := i+1; j < n; j++ {
            index := indices[j]
            value := math.Abs(a[index][i])
            if value > pivotValue {
                pivot = j
                pivotValue = value
            }
        }

        if fuzzyEquals(pivotValue, 0) {
            return false, nil
        }

        if pivot != i {
            indices[i], indices[pivot] = indices[pivot], indices[i]
            b[i], b[pivot] = b[pivot], b[i]
        }
        pivot = indices[i]

        // eliminate
        for j := i+1; j < n; j++ {
            index := indices[j]
            q := -a[index][i] / a[pivot][i]
            for k := i+1; k < n; k++ {
                a[index][k] += a[pivot][k] * q
            }
            b[j] += b[i] * q
        }
    }

    // backwards substitution
    for i := n-1; i >= 0; i-- {
        index := indices[i]
        sum := b[i]
        for j := i+1; j < n; j++ {
            sum -= a[index][j] * b[j]
        }
        
        if fuzzyEquals(a[index][i], 0) {
            return false, nil
        }
        b[i] = sum / a[index][i]
    }
    return true, b
}

func computeDependencies(a [][]float64, m, n int) (int, []bool) {
    // index array for row permutation
    indices := make([]int, m)
    for i := 0; i < m; i++ {
        indices[i] = i
    }

    independent := make([]bool, m)

    // forward elimination
    iterations := m
    if m > n {
        iterations = n
    }
    i := 0
    column := 0
    for ; i < iterations && column < n; i++ {
        // find next pivot
        pivot := i
        for true {
            pivotValue := math.Abs(a[indices[i]][column])
            for j := i+1; j < m; j++ {
                index := indices[j]
                value := math.Abs(a[index][column])
                if value > pivotValue {
                    pivot = j
                    pivotValue = value
                }
            }

            if fuzzyEquals(pivotValue, 0) { break }
            column++

            if column >= n {
                break
            }
        }

        if column == n {
            break
        }

        if pivot != i {
            indices[i], indices[pivot] = indices[pivot], indices[i]
        }
        pivot = indices[i]

        independent[pivot] = true

        // eliminate
        for j := i+1; j < m; j++ {
            index := indices[j]
            q := -a[index][column] / a[pivot][column]
            a[index][column] = 0
            for k := column+1; k < n; k++ {
                a[index][k] += a[pivot][k] * q
            }
        }
        column++
    }

    for j := i; j < m; j++ {
        independent[indices[j]] = false
    }
    return i, independent
}

func removeLinearyDependentRows(A [][]float64, m, n int) int {
    // copy to temp
    temp := copyMatrix(A, m, n)

    count, independentRows := computeDependencies(temp, m, n)
    if count == m {
        return count
    }

    // remove the rows
    index := 0
    for i := 0; i < m; i++ {
        if independentRows[i] {
            if index < i {
                for k := 0; k < n; k++ {
                    A[index][k] = A[i][k]
                }
            }
            index++
        }
    }

    return count
}

// QR decomposition using Householder transformations.
func qrDecomposition(a [][]float64, m, n int, d []float64, q [][]float64) bool {
    if m < n {
        return false
    }

    for j := 0; j < n; j++ {
        // inner product of the first vector x of the (j,j) minor
        innerProductU := 0.0
        for i := j+1; i < m; i++ {
            innerProductU = innerProductU + a[i][j] * a[i][j]
        }
        innerProduct := innerProductU + a[j][j] * a[j][j]
        if fuzzyEquals(innerProduct, 0) {
            return false
        }

        // alpha (norm of x with opposite signedness of x_1) and thus r_{j,j}
        alpha := -math.Sqrt(innerProduct)
        if a[j][j] < 0 {
            alpha = math.Sqrt(innerProduct)
        }
        d[j] = alpha

        beta := 1 / (alpha * a[j][j] - innerProduct)

        // u = x - alpha * e_1
        // (u is a[j..n][j]
        a[j][j] -= alpha

        // left-multiply A_k with Q_k, thus obtaining a row of R and the A_{k+1}
        // for the next iteration
        for k := j+1; k < n; k++ {
            sum := 0.0
            for i := j; i < m; i++ {
                sum += a[i][j] * a[i][k]
            }
            sum *= beta

            for i := j; i < m; i++ {
                a[i][k] += a[i][j] * sum
            }
        }

        // v = u/|u|
        innerProductU += a[j][j] * a[j][j]
        beta2 := -2 / innerProductU

        // right-multiply Q with Q_k
        // Q_k = I - 2vv^T
        // Q * Q_k = Q - 2 * Q * vv^T
        if j == 0 {
            for k := 0; k < m; k++ {
                for i := 0; i < m; i++ {
                    q[k][i] = beta2 * a[k][0] * a[i][0]
                }
                q[k][k] += 1
            }
        } else {
            for k := 0; k < m; k++ {
                sum := 0.0
                for i := j; i < m; i++ {
                    sum += q[k][i] * a[i][j]
                }
                sum *= beta2

                for i := j; i < m; i++ {
                    q[k][i] += sum * a[i][j]
                }
            }
        }

    }
    return true
}
