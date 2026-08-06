package main

import (
	"fmt"
	"math"
)

func SquareRoot(a float64) float64 {
	if a < 0 {
		fmt.Println("No negatives SQRT")
		return 0
	}
	return math.Sqrt(a)
}
