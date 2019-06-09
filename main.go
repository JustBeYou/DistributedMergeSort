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

func main() {
	fmt.Println("Parallel implementation of merge sort algorithm")

	v := randomVector(100, 1000)
	for _, val := range v {
		fmt.Printf("%d ", val)
	}
	fmt.Printf("\n")
}
