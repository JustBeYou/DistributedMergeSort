package main

import (
	"fmt"
	"crypto/rand"
	"math/big"
)

func randomInteger(max int64) int64 {
	ret, _ := rand.Int(rand.Reader, big.NewInt(max))
	return ret.Int64() 
}

func randomVector(size, max int64) []int64 {
	ret := make([]int64, size)
	for i := int64(0); i < size; i++ {
		ret[i] = randomInteger(max)
	}
	return ret
}

// ascending order
func clasicMergeSort(vec []int64) []int64 {
	sorted := make([]int64, len(vec))
	
	if len(vec) == 1 {
		sorted[0] = vec[0]
	} else if len(vec) == 2 {
		if vec[0] < vec[1] {
			sorted[0], sorted[1] = vec[0], vec[1]
		} else {
			sorted[0], sorted[1] = vec[1], vec[0]
		}
	} else {
		A := clasicMergeSort(vec[:len(vec)/2])
		B := clasicMergeSort(vec[len(vec)/2:])
		i, j, idx := 0, 0, 0

		for i < len(A) && j < len(B) {
			if A[i] < B[j] {
				sorted[idx] = A[i]
				idx++
				i++
			} else if A[i] > B[j]{
				sorted[idx] = B[j]
				idx++
				j++
			} else {
				sorted[idx], sorted[idx + 1] = A[i], B[j]
				i++
				j++
			}
		}

		for i < len(A) {
			sorted[idx] = A[i]
			i++
			idx++
		}
		
		for j < len(B) {
			sorted[idx] = B[j]
			j++
			idx++
		}
	}

	return sorted
}

func main() {
	fmt.Println("Parallel implementation of merge sort algorithm")

	v := randomVector(10, 1000)
	fmt.Println(v)
	sorted_v := clasicMergeSort(v)
	fmt.Println(sorted_v)

}
