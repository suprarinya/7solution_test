package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func findMaxPath(pyramid [][]int) (int, []int) {
	n := len(pyramid)
	if n == 0 {
		return 0, nil
	}

	maxSum := make([][]int, n)
	for i := range maxSum {
		maxSum[i] = make([]int, len(pyramid[i]))
		copy(maxSum[i], pyramid[i])
	}

	path := make([][]int, n-1) 
	for i := range path {
		path[i] = make([]int, len(pyramid[i]))
	}

	for i := n - 2; i >= 0; i-- {
		for j := 0; j < len(pyramid[i]); j++ {
			if maxSum[i+1][j] > maxSum[i+1][j+1] {
				maxSum[i][j] += maxSum[i+1][j]
				path[i][j] = j 
			} else {
				maxSum[i][j] += maxSum[i+1][j+1]
				path[i][j] = j + 1 
			}
		}
	}

	resultPath := make([]int, n)
	j := 0
	for i := 0; i < n-1; i++ {
		resultPath[i] = pyramid[i][j]
		j = path[i][j] 
	}
	resultPath[n-1] = pyramid[n-1][j] 

	return maxSum[0][0], resultPath
}

func main() {	
	file, err := ioutil.ReadFile("hard.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var pyramid [][]int
	err = json.Unmarshal(file, &pyramid)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	maxTotal, path := findMaxPath(pyramid)
	fmt.Println("Maximum Total Sum:", maxTotal)
	fmt.Println("Path:", path)
}

