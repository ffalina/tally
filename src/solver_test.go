package main

import "testing"

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
