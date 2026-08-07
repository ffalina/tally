package main

import (
	"math"
	"testing"
)

func TestSolveExpressionReportsMalformedExpressions(t *testing.T) {
	tests := []string{
		"1 + 2)",
		"(1 + 2",
		"1 +",
		"()",
		"sqrt()",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if result.Valid {
				t.Fatalf("solveExpression(%q) returned valid result %v", input, result.Result)
			}
			if result.Error == "" {
				t.Fatalf("solveExpression(%q) returned no error", input)
			}
		})
	}
}

func TestSolveExpressionReportsInvalidNumbers(t *testing.T) {
	tests := []string{
		"1.2.3",
		".",
		"-.",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if result.Valid {
				t.Fatalf("solveExpression(%q) returned valid result %v", input, result.Result)
			}
			if result.Error == "" {
				t.Fatalf("solveExpression(%q) returned no error", input)
			}
		})
	}
}

func TestSolveExpressionStillSolvesValidExpression(t *testing.T) {
	result := solveExpression("2 + 3 * 4")
	if !result.Valid {
		t.Fatalf("solveExpression returned error: %s", result.Error)
	}
	if result.Result != 14 {
		t.Fatalf("solveExpression result = %v, want 14", result.Result)
	}
}

func TestSolveExpressionUsesRightAssociativePowers(t *testing.T) {
	result := solveExpression("2^3^2")
	if !result.Valid {
		t.Fatalf("solveExpression returned error: %s", result.Error)
	}
	if result.Result != 512 {
		t.Fatalf("solveExpression result = %v, want 512", result.Result)
	}
}

func TestSolveExpressionKeepsOtherOperatorsLeftAssociative(t *testing.T) {
	result := solveExpression("10 - 3 - 2")
	if !result.Valid {
		t.Fatalf("solveExpression returned error: %s", result.Error)
	}
	if result.Result != 5 {
		t.Fatalf("solveExpression result = %v, want 5", result.Result)
	}
}

func TestSolveExpressionFactorial(t *testing.T) {
	tests := map[string]float64{
		"5!":    120,
		"0!":    1,
		"2^3!":  64, // factorial binds tighter than ^, so this is 2^(3!)
		"3! + 1": 7,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if !result.Valid {
				t.Fatalf("solveExpression(%q) returned error: %s", input, result.Error)
			}
			if result.Result != want {
				t.Fatalf("solveExpression(%q) = %v, want %v", input, result.Result, want)
			}
		})
	}
}

func TestSolveExpressionFactorialRejectsNonIntegers(t *testing.T) {
	tests := []string{"2.5!", "(-3)!"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if result.Valid {
				t.Fatalf("solveExpression(%q) returned valid result %v", input, result.Result)
			}
		})
	}
}

func TestSolveExpressionDegreeTrig(t *testing.T) {
	tests := map[string]float64{
		"sin(90)":  1,
		"cos(0)":   1,
		"asin(1)":  90,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if !result.Valid {
				t.Fatalf("solveExpression(%q) returned error: %s", input, result.Error)
			}
			if math.Abs(result.Result-want) > 1e-9 {
				t.Fatalf("solveExpression(%q) = %v, want %v", input, result.Result, want)
			}
		})
	}
}

func TestSolveExpressionLogWithBase(t *testing.T) {
	result := solveExpression("log(8, 2)")
	if !result.Valid {
		t.Fatalf("solveExpression returned error: %s", result.Error)
	}
	if result.Result != 3 {
		t.Fatalf("solveExpression result = %v, want 3", result.Result)
	}
}

func TestSolveExpressionReportsDomainErrors(t *testing.T) {
	tests := []string{"sqrt(-4)", "log(-1)", "log(8, 1)", "asin(2)"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if result.Valid {
				t.Fatalf("solveExpression(%q) returned valid result %v", input, result.Result)
			}
			if result.Error == "" {
				t.Fatalf("solveExpression(%q) returned no error", input)
			}
		})
	}
}

func TestSolveExpressionCombinationAndPermutation(t *testing.T) {
	tests := map[string]float64{
		"C(5, 2)": 10,
		"P(5, 2)": 20,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			result := solveExpression(input)
			if !result.Valid {
				t.Fatalf("solveExpression(%q) returned error: %s", input, result.Error)
			}
			if result.Result != want {
				t.Fatalf("solveExpression(%q) = %v, want %v", input, result.Result, want)
			}
		})
	}
}

func TestSolveExpressionConstants(t *testing.T) {
	result := solveExpression("2 * pi")
	if !result.Valid {
		t.Fatalf("solveExpression returned error: %s", result.Error)
	}
	want := 2 * Pi()
	if result.Result != want {
		t.Fatalf("solveExpression result = %v, want %v", result.Result, want)
	}
}

func TestSolveLinearEquationSupportsConstants(t *testing.T) {
	result := solveLinearEquation("2x + pi = 10")
	if !result.Valid {
		t.Fatalf("solveLinearEquation returned error: %s", result.Error)
	}
	want := (10 - Pi()) / 2
	if math.Abs(result.Result-want) > 1e-9 {
		t.Fatalf("solveLinearEquation result = %v, want %v", result.Result, want)
	}
}

func TestSolveLinearEquationReportsNoSolution(t *testing.T) {
	result := solveLinearEquation("x + 1 = x + 2")
	if result.Valid {
		t.Fatalf("solveLinearEquation(%q) returned valid result %v", "x + 1 = x + 2", result.Result)
	}
	if result.Error == "" {
		t.Fatalf("solveLinearEquation returned no error")
	}
}
