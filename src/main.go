package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	Number TokenType = iota
	Operator
	Function
	Variable
	LeftParen
	RightParen
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	i := 0
	expectNumber := true

	for i < len(input) {
		ch := input[i]

		if ch == ' ' {
			i++
			continue
		}

		if unicode.IsDigit(rune(ch)) || ch == '.' || (ch == '-' && expectNumber && i+1 < len(input) && (unicode.IsDigit(rune(input[i+1])) || input[i+1] == '.')) {
			start := i
			hasDigit := false
			if ch == '-' {
				i++
			}
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				if unicode.IsDigit(rune(input[i])) {
					hasDigit = true
				}
				i++
			}
			value := input[start:i]
			if !hasDigit {
				return nil, fmt.Errorf("invalid number %q", value)
			}
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				return nil, fmt.Errorf("invalid number %q", value)
			}
			tokens = append(tokens, Token{Number, value})
			expectNumber = false
			continue
		}

		if unicode.IsLetter(rune(ch)) {
			start := i
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			value := input[start:i]
			tokenType := Function
			if value == "x" {
				tokenType = Variable
			}
			tokens = append(tokens, Token{tokenType, value})
			expectNumber = false
			continue
		}

		if strings.ContainsRune("+-*/^", rune(ch)) {
			tokens = append(tokens, Token{Operator, string(ch)})
			i++
			expectNumber = true
			continue
		}

		if ch == '(' {
			tokens = append(tokens, Token{LeftParen, "("})
			i++
			expectNumber = true
			continue
		}
		if ch == ')' {
			tokens = append(tokens, Token{RightParen, ")"})
			i++
			expectNumber = false
			continue
		}

		return nil, fmt.Errorf("unknown character: %s", string(ch))
	}

	return tokens, nil
}

var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

func isRightAssociative(operator string) bool {
	return operator == "^"
}

func toRPN(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	for _, tok := range tokens {
		switch tok.Type {

		case Number, Variable:
			output = append(output, tok)

		case Function:
			stack = append(stack, tok)

		case Operator:
			if _, ok := precedence[tok.Value]; !ok {
				return nil, fmt.Errorf("unknown operator: %s", tok.Value)
			}
			for len(stack) > 0 {
				top := stack[len(stack)-1]

				shouldPopOperator := top.Type == Operator &&
					(precedence[top.Value] > precedence[tok.Value] ||
						(precedence[top.Value] == precedence[tok.Value] && !isRightAssociative(tok.Value)))
				if top.Type == Function || shouldPopOperator {

					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, tok)

		case LeftParen:
			stack = append(stack, tok)

		case RightParen:
			for len(stack) > 0 && stack[len(stack)-1].Type != LeftParen {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]

			if len(stack) > 0 && stack[len(stack)-1].Type == Function {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].Type == LeftParen || stack[len(stack)-1].Type == RightParen {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evalExpression(input string) float64 {
	result := solveExpression(input)
	if !result.Valid {
		panic(result.Error)
	}
	return result.Result
}

func main() {
	tests := []string{
		"2 + 3 * 4",
		"sqrt(16)",
		"sin(0)",
		"2^3 + 1",
		"abs(-5) + 3",
		"3 + 4 * 2 / (1 - 5)^2",
	}

	for _, t := range tests {
		solution := solveExpression(t)
		fmt.Println(t)
		for i, step := range solution.Steps {
			fmt.Printf("Step %d: %s\n", i+1, step)
		}
		fmt.Println("Answer:", formatNumber(solution.Result))
		fmt.Println()
	}
}
