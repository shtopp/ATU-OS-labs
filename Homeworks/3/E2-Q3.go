package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Merge slices
func merge(left, right []int) []int {
	size, i, j := len(left)+len(right), 0, 0
	merged := make([]int, size)

	for k := 0; k < size; k++ {
		if i >= len(left) {
			merged[k] = right[j]
			j++
		} else if j >= len(right) {
			merged[k] = left[i]
			i++
		} else if left[i] < right[j] {
			merged[k] = left[i]
			i++
		} else {
			merged[k] = right[j]
			j++
		}
	}
	return merged
}

// parallel merge sort
func ParallelMergeSort(arr []int) []int {

	len := len(arr)

	if len <= 1 {
		return arr
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	var wg sync.WaitGroup
	wg.Add(2)

	mid := len / 2
	var left, right []int

	//split
	go func() {
		defer wg.Done()
		left = ParallelMergeSort(arr[:mid])
	}()

	go func() {
		defer wg.Done()
		right = ParallelMergeSort(arr[mid:])
	}()

	wg.Wait()

	//merge
	return merge(left, right)
}

func main() {
	// Generate a random unsorted slice
	arr := rand.Perm(1000)

	//perform the simple merge sort
	startTime := time.Now()
	sorted := ParallelMergeSort(arr)
	elapsedTime := time.Since(startTime)

	fmt.Println(sorted[:10])
	fmt.Println("Elapsed Time:", elapsedTime)
}
