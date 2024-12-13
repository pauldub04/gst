package client

import "sync"

func MultiplyMatrixVector(matrix [][]int32, vector []int32) []int32 {
	rows := len(matrix)
	result := make([]int32, rows)

	var wg sync.WaitGroup

	for i := 0; i < rows; i++ {
		wg.Add(1)
		go func(rowIndex int) {
			defer wg.Done()
			for j := range vector {
				result[rowIndex] += matrix[rowIndex][j] * vector[j]
			}
		}(i)
	}

	wg.Wait()
	return result
}
