package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// setup scanner to read input from keyboard
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter encoded text (e.g., LLRR=):")
	scanner.Scan()
	encoded := strings.TrimSpace(scanner.Text())

	// create the number sequence array with enough space
	numbers := make([]int, len(encoded)+1)

	// start from the end and work backwards to determine minimal numbers
	for i := len(encoded) - 1; i >= 0; i-- {
		switch encoded[i] {
			case 'L':
				// left number is greater than the right
				numbers[i] = numbers[i+1] + 1
			case 'R':
				// right number is greater than the left
				// If right = 0 and increment all numbers from right
				if numbers[i+1] == 0 {
					numbers[i] = 0
					for j := i + 1; j <= len(encoded); j++ {
						numbers[j]++
					}
				} else {
					numbers[i] = 0
				}
			case '=':
				// number is equal
				numbers[i] = numbers[i+1]
		}
	}

	// calc sum
	sum := 0
	for _, num := range numbers {
		sum += num
	}

	fmt.Printf("The lowest number sequence: ")
	for _, num := range numbers {
		fmt.Print(num)
	}
	fmt.Println()
	fmt.Printf("Sum of all numbers: %d\n", sum)
}
