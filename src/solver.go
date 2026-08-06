package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type SolveResult struct {
	Result float64
	Steps  []string
	Valid  bool
	Error  string
}

type solveValue struct {
	value float64
	label string
}

func solveExpression(input string) SolveResult {
	tokens, err := tokenize(input)
	if err != nil {
		return SolveResult{Error: err.Error()}
	}
	rpn, err := toRPN(tokens)
	if err != nil {
		return SolveResult{Error: err.Error()}
	}
	return evaluateRPN(rpn)
}

func evaluateRPN(tokens []Token) SolveResult {
	var stack []solveValue
	var steps []string

	for _, tok := range tokens {
		switch tok.Type {
		case Number:
			val, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				return SolveResult{Error: fmt.Sprintf("invalid number %q", tok.Value)}
			}
			stack = append(stack, solveValue{value: val, label: tok.Value})

		case Operator:
			if len(stack) < 2 {
				return SolveResult{Error: fmt.Sprintf("operator %q is missing an operand", tok.Value)}
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res := evaluateOperator(tok.Value, a.value, b.value)
			step := fmt.Sprintf("%s %s %s = %s", a.label, tok.Value, b.label, formatNumber(res))
			steps = append(steps, step)
			stack = append(stack, solveValue{value: res, label: formatNumber(res)})

		case Function:
			if len(stack) < 1 {
				return SolveResult{Error: fmt.Sprintf("function %q is missing an argument", tok.Value)}
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			res := evaluateFunction(tok.Value, a.value)
			step := fmt.Sprintf("%s(%s) = %s", tok.Value, a.label, formatNumber(res))
			steps = append(steps, step)
			stack = append(stack, solveValue{value: res, label: formatNumber(res)})
		}
	}

	if len(stack) != 1 {
		return SolveResult{Error: "invalid expression"}
	}

	return SolveResult{
		Result: stack[0].value,
		Steps:  steps,
		Valid:  true,
	}
}

func evaluateOperator(operator string, a, b float64) float64 {
	switch operator {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "^":
		return math.Pow(a, b)
	default:
		panic("unknown operator: " + operator)
	}
}

func evaluateFunction(function string, a float64) float64 {
	switch function {
	case "sin":
		return math.Sin(a)
	case "cos":
		return math.Cos(a)
	case "tan":
		return math.Tan(a)
	case "log":
		return math.Log(a)
	case "sqrt":
		return math.Sqrt(a)
	case "abs":
		return math.Abs(a)
	default:
		panic("unknown function: " + function)
	}
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

type LinearEquationResult struct {
	Result float64
	Steps  []string
	Valid  bool
	Error  string
}

type linearValue struct {
	coefficient float64
	constant    float64
}

func solveLinearEquation(input string) (result LinearEquationResult) {
	defer func() {
		if r := recover(); r != nil {
			result = invalidEquation(fmt.Sprintf("invalid equation: %v", r))
		}
	}()

	parts := strings.Split(input, "=")
	if len(parts) != 2 {
		return invalidEquation("equation must contain exactly one equals sign")
	}

	leftText := strings.TrimSpace(parts[0])
	rightText := strings.TrimSpace(parts[1])
	if leftText == "" || rightText == "" {
		return invalidEquation("both sides of the equation must contain an expression")
	}

	left, err := linearExpression(leftText)
	if err != nil {
		return invalidEquation(err.Error())
	}

	right, err := linearExpression(rightText)
	if err != nil {
		return invalidEquation(err.Error())
	}

	coefficient := left.coefficient - right.coefficient
	constant := right.constant - left.constant
	steps := []string{
		fmt.Sprintf("Start with %s = %s", leftText, rightText),
		fmt.Sprintf("Combine like terms: %s = %s", formatLinear(left), formatLinear(right)),
	}

	if right.coefficient != 0 {
		leftWithoutRightX := linearValue{coefficient: coefficient, constant: left.constant}
		steps = append(steps, fmt.Sprintf(
			"Move x terms to the left: %s = %s",
			formatLinear(leftWithoutRightX),
			formatNumber(right.constant),
		))
	}
	if left.constant != 0 {
		steps = append(steps, fmt.Sprintf(
			"Move constants to the right: %s = %s",
			formatLinear(linearValue{coefficient: coefficient}),
			formatNumber(constant),
		))
	}

	if coefficient == 0 {
		if constant == 0 {
			return LinearEquationResult{Steps: steps, Error: "equation has infinitely many solutions"}
		}
		return LinearEquationResult{Steps: steps, Error: "equation has no solution"}
	}

	x := constant / coefficient
	steps = append(steps, fmt.Sprintf(
		"Divide both sides by %s: x = %s",
		formatNumber(coefficient),
		formatNumber(x),
	))

	return LinearEquationResult{
		Result: x,
		Steps:  steps,
		Valid:  true,
	}
}

func invalidEquation(message string) LinearEquationResult {
	return LinearEquationResult{Error: message}
}

func linearExpression(input string) (linearValue, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return linearValue{}, err
	}
	tokens = insertImplicitMultiplication(tokens)
	rpn, err := toRPN(tokens)
	if err != nil {
		return linearValue{}, err
	}
	return evaluateLinearRPN(rpn)
}

func insertImplicitMultiplication(tokens []Token) []Token {
	if len(tokens) < 2 {
		return tokens
	}

	withMultiplication := make([]Token, 0, len(tokens))
	for i, tok := range tokens {
		if i > 0 && needsImplicitMultiplication(tokens[i-1], tok) {
			withMultiplication = append(withMultiplication, Token{Type: Operator, Value: "*"})
		}
		withMultiplication = append(withMultiplication, tok)
	}

	return withMultiplication
}

func needsImplicitMultiplication(left, right Token) bool {
	leftCanMultiply := left.Type == Number || left.Type == Variable || left.Type == RightParen
	rightCanMultiply := right.Type == Number || right.Type == Variable || right.Type == LeftParen
	return leftCanMultiply && rightCanMultiply
}

func evaluateLinearRPN(tokens []Token) (linearValue, error) {
	var stack []linearValue

	for _, tok := range tokens {
		switch tok.Type {
		case Number:
			val, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				return linearValue{}, fmt.Errorf("invalid number %q", tok.Value)
			}
			stack = append(stack, linearValue{constant: val})

		case Variable:
			if tok.Value != "x" {
				return linearValue{}, fmt.Errorf("unsupported variable %q", tok.Value)
			}
			stack = append(stack, linearValue{coefficient: 1})

		case Operator:
			if len(stack) < 2 {
				return linearValue{}, fmt.Errorf("operator %q is missing an operand", tok.Value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := evaluateLinearOperator(tok.Value, a, b)
			if err != nil {
				return linearValue{}, err
			}
			stack = append(stack, res)

		case Function:
			if len(stack) < 1 {
				return linearValue{}, fmt.Errorf("unsupported token %q", tok.Value)
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if hasVariable(a) {
				return linearValue{}, fmt.Errorf("functions with x are not supported in linear equations")
			}
			stack = append(stack, linearValue{constant: evaluateFunction(tok.Value, a.constant)})
		}
	}

	if len(stack) != 1 {
		return linearValue{}, fmt.Errorf("invalid equation expression")
	}
	return stack[0], nil
}

func evaluateLinearOperator(operator string, a, b linearValue) (linearValue, error) {
	switch operator {
	case "+":
		return linearValue{
			coefficient: a.coefficient + b.coefficient,
			constant:    a.constant + b.constant,
		}, nil
	case "-":
		return linearValue{
			coefficient: a.coefficient - b.coefficient,
			constant:    a.constant - b.constant,
		}, nil
	case "*":
		if hasVariable(a) && hasVariable(b) {
			return linearValue{}, fmt.Errorf("nonlinear multiplication is not supported")
		}
		if hasVariable(a) {
			return linearValue{
				coefficient: a.coefficient * b.constant,
				constant:    a.constant * b.constant,
			}, nil
		}
		return linearValue{
			coefficient: b.coefficient * a.constant,
			constant:    b.constant * a.constant,
		}, nil
	case "/":
		if hasVariable(b) {
			return linearValue{}, fmt.Errorf("division by an expression containing x is not supported")
		}
		if b.constant == 0 {
			return linearValue{}, fmt.Errorf("division by zero")
		}
		return linearValue{
			coefficient: a.coefficient / b.constant,
			constant:    a.constant / b.constant,
		}, nil
	case "^":
		if hasVariable(b) {
			return linearValue{}, fmt.Errorf("variable exponents are not supported")
		}
		if !hasVariable(a) {
			return linearValue{constant: math.Pow(a.constant, b.constant)}, nil
		}
		if b.constant == 1 {
			return a, nil
		}
		if b.constant == 0 {
			return linearValue{constant: 1}, nil
		}
		return linearValue{}, fmt.Errorf("powers of x other than 1 are not linear")
	default:
		return linearValue{}, fmt.Errorf("unknown operator %q", operator)
	}
}

func hasVariable(value linearValue) bool {
	return value.coefficient != 0
}

func formatLinear(value linearValue) string {
	coefficient := value.coefficient
	constant := value.constant

	if coefficient == 0 {
		return formatNumber(constant)
	}

	var builder strings.Builder
	switch coefficient {
	case 1:
		builder.WriteString("x")
	case -1:
		builder.WriteString("-x")
	default:
		builder.WriteString(formatNumber(coefficient))
		builder.WriteString("x")
	}

	if constant > 0 {
		builder.WriteString(" + ")
		builder.WriteString(formatNumber(constant))
	} else if constant < 0 {
		builder.WriteString(" - ")
		builder.WriteString(formatNumber(-constant))
	}

	return builder.String()
}
