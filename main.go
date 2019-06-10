package main

import (
	"fmt"
	"crypto/rand"
	"math/big"
	"time"
	"sync"
	"os"
	"strconv"
)

func randomInteger(max int64) int64 {
	ret, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic("Can't get random data")
	}
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
				idx += 2
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

func parallelMergeSort_A(vec []int64, ch chan int64) {
	if len(vec) == 1 {
		ch <- vec[0]
	} else if len(vec) == 2 {
		if vec[0] < vec[1] {
			ch <- vec[0]
			ch <- vec[1]
		} else {
			ch <- vec[1]
			ch <- vec[0]
		}
	} else {
		A := make(chan int64)
		B := make(chan int64)
		
		go parallelMergeSort_A(vec[:len(vec)/2], A)
		go parallelMergeSort_A(vec[len(vec)/2:], B)
		
		get_new_A, get_new_B := true, true
		var val_A, val_B int64
		var ok_A, ok_B bool
		for {
			if get_new_A {
				val_A, ok_A = <-A
			}
			if get_new_B {
				val_B, ok_B = <-B
			}

			if ok_A == false && ok_B == true {
				A = nil
				ch <- val_B
				break
			} else if ok_A == true && ok_B == false {
				B = nil
				ch <- val_A
				break
			} else if ok_A == false && ok_B == false {
				A, B = nil, nil
				break
			}

			if val_A < val_B {
				ch <- val_A
				get_new_A = true
				get_new_B = false
			} else if val_A > val_B {
				ch <- val_B
				get_new_A = false
				get_new_B = true
			} else {
				ch <- val_A
				ch <- val_B
				get_new_A = true
				get_new_B = true
			}
		}

		if A != nil {
			for i := range A {
				ch <- i
			}
		}

		if B != nil {
			for i:= range B {
				ch <- i
			}
		}
	}

	close(ch)
}

func parallelMergeSort_B(vec, sorted []int64, wg *sync.WaitGroup, level int64) {
	if len(vec) == 1 {
		sorted[0] = vec[0]
	} else if len(vec) == 2 {
		if vec[0] < vec[1] {
			sorted[0], sorted[1] = vec[0], vec[1]
		} else {
			sorted[0], sorted[1] = vec[1], vec[0]
		}
	} else {
		A, B := make([]int64, len(sorted)/2), make([]int64, len(sorted) - len(sorted)/2)
		parallelMergeSort_B(vec[:len(vec)/2], A, wg, level + 1)
		parallelMergeSort_B(vec[len(vec)/2:], B, wg, level + 1)
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
				idx += 2
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

	if level == 0 {
		wg.Done()
	}
}


func mergeArrays(arrays [][]int64, sorted []int64) {
	indexes := make([]int, len(arrays))
	idx := 0
	good := true

	for good {
		good = false

		min_val, min_pos := int64(1e15), -1
		for i := 0; i < len(arrays); i++ {
			if indexes[i] >= len(arrays[i]) {
				continue
			}

			if arrays[i][indexes[i]] < min_val {
				min_val = arrays[i][indexes[i]]
				min_pos = i
			}
		}

		if min_pos != -1 {
			good = true
			sorted[idx] = arrays[min_pos][indexes[min_pos]]
			idx++
			indexes[min_pos]++
		}
	}
}

func main() {
	fmt.Println("Parallel implementation of merge sort algorithm")
	if len(os.Args) < 2 {
		panic("Too few arguments. Usage: ./exe CPUs")
	}

	v := randomVector(3 * 10e6, 10e8)
	fmt.Println("Generated the vector")
	//fmt.Printf("Initial vector: ")
	//fmt.Println(v)

	start := time.Now()
	sorted_v := clasicMergeSort(v)
	t := time.Now()
	elapsed := t.Sub(start)

	
	//fmt.Printf("Sorted vector: ")
	//fmt.Println(sorted_v)
	fmt.Printf("Classic merge sort: ")
	fmt.Println(elapsed)

	/*ch := make(chan int64) VERY VERY SLOOOOOW

	start = time.Now()
	go parallelMergeSort_A(v, ch)
	
	i := 0
	for val := range ch {
		if val != sorted_v[i] {
			panic("Parallel merge sort A failed")
		}
		i++
	}
	t = time.Now()
	elapsed = t.Sub(start)

	fmt.Printf("Parallel (A) merge sort: ")
	fmt.Println(elapsed)*/

	CPUs, _ := strconv.Atoi(os.Args[1])
	sorted_K := make([][]int64, CPUs)
	sorted_B := make([]int64, len(v))
	for i := 0; i < CPUs - 1; i++ {
		sorted_K[i] = make([]int64, len(v) / CPUs)
	}
	sorted_K[CPUs - 1] = make([]int64, len(v) - (CPUs - 1) * (len(v) / CPUs))

	start = time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(CPUs)
	for i := 0; i < CPUs - 1; i++ {
		go parallelMergeSort_B(v[i * (len(v) / CPUs):(i + 1) * (len(v) / CPUs)], sorted_K[i], wg, 0)
	}
	parallelMergeSort_B(v[(CPUs - 1) * (len(v) / CPUs):], sorted_K[CPUs - 1], wg, 0)
	wg.Wait()
	mergeArrays(sorted_K, sorted_B)

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Parallel (B) merge sort: ")
	fmt.Println(elapsed)

	for i := 0; i < len(v); i++ {
		if sorted_B[i] != sorted_v[i] {
			panic("Parallel merge sort B failed")
		}
	}
}
