package main

import (
	"fmt"
	"math"
)

func Log(a float64) float64 {
	if a <= 0 {
		fmt.Println("Log only works for positive numbers")
		return 0
	}
	return math.Log(a)
}

func Log10(a float64) float64 {
	if a <= 0 {
		fmt.Println("Log only works for positive numbers")
		return 0
	}
	return math.Log10(a)
}

func LogBase(a float64, base float64) float64 {
	return math.Log(a) / math.Log(base)
}
