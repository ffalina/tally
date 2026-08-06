package main

func Factorial(a int) int {
	if a < 0 {
		return 0
	}

	b := 1
	for i := 2; i <= a; i++ {
		b *= i
	}

	return b
}

func Combination(a int, b int) int {
	if b > a || b < 0 {
		return 0
	}

	if b > a-b {
		b = a - b
	}

	c := 1
	for i := 0; i < b; i++ {
		c *= (a - i)
		c /= (i + 1)
	}

	return c
}

func Permuration(a int, b int) int {
	if b > a || b < 0 {
		return 0
	}

	c := 1
	for i := 0; i < b; i++ {
		c *= (a - i)
	}

	return c
}
