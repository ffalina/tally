package main

import (
	"fmt"
	"math"
)

func Sin(a float64) float64 {
	b := a * math.Pi / 180
	return math.Sin(b)
}

func Cos(a float64) float64 {
	b := a * math.Pi / 180
	return math.Cos(b)
}

func Tan(a float64) float64 {
	b := a * math.Pi / 180
	return math.Tan(b)
}

func Asin(a float64) float64 {
	if a < -1 || a > 1 {
		fmt.Println("asin input must be between -1 and 1")
		return math.NaN()
	}
	b := math.Asin(a)
	return b * 180 / math.Pi
}

func Acos(a float64) float64 {
	if a < -1 || a > 1 {
		fmt.Println("acos input must be between -1 and 1")
		return math.NaN()
	}
	b := math.Acos(a)
	return b * 180 / math.Pi
}

func Atan(a float64) float64 {
	b := math.Atan(a)
	return b * 180 / math.Pi
}
