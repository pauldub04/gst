package client

import (
	"sync"

	"golang.org/x/exp/rand"
)

func GenerateMatrix(seed uint64, rows int32, cols int32, min int32, max int32) [][]int32 {
	matrix := make([][]int32, rows)

	var wg sync.WaitGroup
	for i := int32(0); i < rows; i++ {
		wg.Add(1)
		go func(rowIndex int32) {
			defer wg.Done()
			matrix[rowIndex] = GenerateVector(seed+uint64(rowIndex), cols, min, max)
		}(i)
	}

	wg.Wait()
	return matrix
}

func GenerateVector(seed uint64, len int32, min int32, max int32) []int32 {
	randGen := rand.New(rand.NewSource(seed))
	vector := make([]int32, len)
	for i := range vector {
		vector[i] = randGen.Int31n(max-min+1) + min
	}
	return vector
}
