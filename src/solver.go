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
			argc := tok.Argc
			if argc < 1 {
				argc = 1
			}
			if len(stack) < argc {
				return SolveResult{Error: fmt.Sprintf("function %q requires %d argument(s)", tok.Value, argc)}
			}
			args := make([]float64, argc)
			labels := make([]string, argc)
			for i := argc - 1; i >= 0; i-- {
				args[i] = stack[len(stack)-1].value
				labels[i] = stack[len(stack)-1].label
				stack = stack[:len(stack)-1]
			}

			res, err := evaluateFunction(tok.Value, args)
			if err != nil {
				return SolveResult{Error: err.Error()}
			}
			step := fmt.Sprintf("%s(%s) = %s", tok.Value, strings.Join(labels, ", "), formatNumber(res))
			steps = append(steps, step)
			stack = append(stack, solveValue{value: res, label: formatNumber(res)})

		case Postfix:
			if len(stack) < 1 {
				return SolveResult{Error: fmt.Sprintf("%q is missing an operand", tok.Value)}
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			res, err := evaluateFactorial(a.value)
			if err != nil {
				return SolveResult{Error: err.Error()}
			}
			step := fmt.Sprintf("%s%s = %s", a.label, tok.Value, formatNumber(res))
			steps = append(steps, step)
			stack = append(stack, solveValue{value: res, label: formatNumber(res)})

		case Constant:
			val, err := evaluateConstant(tok.Value)
			if err != nil {
				return SolveResult{Error: err.Error()}
			}
			stack = append(stack, solveValue{value: val, label: tok.Value})
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
		return Add(a, b)
	case "-":
		return Sub(a, b)
	case "*":
		return Multiply(a, b)
	case "/":
		return Divide(a, b)
	case "^":
		return Power(a, b)
	default:
		panic("unknown operator: " + operator)
	}
}

var unaryFunctions = map[string]bool{
	"sin": true, "cos": true, "tan": true,
	"asin": true, "acos": true, "atan": true,
	"log10": true, "sqrt": true, "abs": true,
	"floor": true, "ceil": true, "round": true,
}

// evaluateFunction dispatches named functions to the corresponding helpers.
// Most take exactly one argument; "log" also accepts a second argument (the
// base), and "C"/"P" (combination/permutation) always take two.
func evaluateFunction(function string, args []float64) (float64, error) {
	if unaryFunctions[function] && len(args) != 1 {
		return 0, fmt.Errorf("%s requires exactly 1 argument", function)
	}

	switch function {
	case "sin":
		return Sin(args[0]), nil
	case "cos":
		return Cos(args[0]), nil
	case "tan":
		return Tan(args[0]), nil
	case "asin":
		if args[0] < -1 || args[0] > 1 {
			return 0, fmt.Errorf("asin input must be between -1 and 1")
		}
		return Asin(args[0]), nil
	case "acos":
		if args[0] < -1 || args[0] > 1 {
			return 0, fmt.Errorf("acos input must be between -1 and 1")
		}
		return Acos(args[0]), nil
	case "atan":
		return Atan(args[0]), nil
	case "log":
		if len(args) == 2 {
			if args[1] <= 0 || args[1] == 1 {
				return 0, fmt.Errorf("log base must be positive and not equal to 1")
			}
			if args[0] <= 0 {
				return 0, fmt.Errorf("log input must be positive")
			}
			return LogBase(args[0], args[1]), nil
		}
		if len(args) != 1 {
			return 0, fmt.Errorf("log requires 1 or 2 arguments")
		}
		if args[0] <= 0 {
			return 0, fmt.Errorf("log input must be positive")
		}
		return Log(args[0]), nil
	case "log10":
		if args[0] <= 0 {
			return 0, fmt.Errorf("log10 input must be positive")
		}
		return Log10(args[0]), nil
	case "sqrt":
		if args[0] < 0 {
			return 0, fmt.Errorf("sqrt input must be non-negative")
		}
		return SquareRoot(args[0]), nil
	case "abs":
		return Abs(args[0]), nil
	case "floor":
		return Floor(args[0]), nil
	case "ceil":
		return Ceil(args[0]), nil
	case "round":
		return Round(args[0]), nil
	case "C":
		if len(args) != 2 {
			return 0, fmt.Errorf("C requires 2 arguments: C(n, r)")
		}
		return float64(Combination(int(args[0]), int(args[1]))), nil
	case "P":
		if len(args) != 2 {
			return 0, fmt.Errorf("P requires 2 arguments: P(n, r)")
		}
		return float64(Permutation(int(args[0]), int(args[1]))), nil
	default:
		return 0, fmt.Errorf("unknown function: %s", function)
	}
}

func evaluateFactorial(a float64) (float64, error) {
	if a < 0 || a != math.Trunc(a) {
		return 0, fmt.Errorf("factorial requires a non-negative integer")
	}
	return float64(Factorial(int(a))), nil
}

func evaluateConstant(name string) (float64, error) {
	switch name {
	case "pi":
		return Pi(), nil
	case "e":
		return E(), nil
	default:
		return 0, fmt.Errorf("unknown constant: %s", name)
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
	leftCanMultiply := left.Type == Number || left.Type == Variable || left.Type == Constant || left.Type == RightParen
	rightCanMultiply := right.Type == Number || right.Type == Variable || right.Type == Constant || right.Type == LeftParen
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
			argc := tok.Argc
			if argc < 1 {
				argc = 1
			}
			if len(stack) < argc {
				return linearValue{}, fmt.Errorf("function %q requires %d argument(s)", tok.Value, argc)
			}
			args := make([]float64, argc)
			hasVar := false
			for i := argc - 1; i >= 0; i-- {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if hasVariable(top) {
					hasVar = true
				}
				args[i] = top.constant
			}
			if hasVar {
				return linearValue{}, fmt.Errorf("functions with x are not supported in linear equations")
			}
			res, err := evaluateFunction(tok.Value, args)
			if err != nil {
				return linearValue{}, err
			}
			stack = append(stack, linearValue{constant: res})

		case Postfix:
			if len(stack) < 1 {
				return linearValue{}, fmt.Errorf("%q is missing an operand", tok.Value)
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if hasVariable(top) {
				return linearValue{}, fmt.Errorf("factorial of x is not supported in linear equations")
			}
			res, err := evaluateFactorial(top.constant)
			if err != nil {
				return linearValue{}, err
			}
			stack = append(stack, linearValue{constant: res})

		case Constant:
			val, err := evaluateConstant(tok.Value)
			if err != nil {
				return linearValue{}, err
			}
			stack = append(stack, linearValue{constant: val})
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
