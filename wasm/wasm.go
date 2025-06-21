package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"syscall/js"

	"github.com/dengsgo/math-engine/engine"
	"github.com/shopspring/decimal"
)

var (
	operatorPrecedence = map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
)

// 科学计数法
func ScienceCalc(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return float64(0)
	}
	expr := args[0].String()
	result, err := engine.ParseAndExec(expr)
	if err != nil {
		return float64(0)
	}
	return result
}

// 将中缀表达式转为逆波兰表达式
func infixToRPN(expr string) ([]string, error) {
	// 匹配数字、操作符和括号
	re := regexp.MustCompile(`\d+(\.\d+)?|[+\-*/()]`)
	tokens := re.FindAllString(expr, -1)
	var output []string
	var stack []string

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" || operatorPrecedence[token] > operatorPrecedence[top] {
					break
				}
				output = append(output, top)
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, errors.New("mismatched parentheses")
			}
		default: // number
			output = append(output, token)
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, top)
	}
	return output, nil
}

// 计算逆波兰表达式
func evalRPN(tokens []string) (decimal.Decimal, error) {
	var stack []decimal.Decimal
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return decimal.Zero, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res decimal.Decimal
			switch token {
			case "+":
				res = a.Add(b)
			case "-":
				res = a.Sub(b)
			case "*":
				res = a.Mul(b)
			case "/":
				if b.IsZero() {
					return decimal.Zero, errors.New("division by zero")
				}
				res = a.Div(b)
			}
			stack = append(stack, res)
		default: // number
			d, err := decimal.NewFromString(token)
			if err != nil {
				return decimal.Zero, err
			}
			stack = append(stack, d)
		}
	}
	if len(stack) != 1 {
		return decimal.Zero, errors.New("invalid expression")
	}
	return stack[0], nil
}
func calcExpression(expr string) (decimal.Decimal, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	rpn, err := infixToRPN(expr)
	if err != nil {
		return decimal.Zero, err
	}
	return evalRPN(rpn)
}

func browserJudgeExpr(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return int32(0)
	}
	expr := args[0].String()
	_, err := CalcExpr(expr)
	if err != nil {
		return int32(0)
	}
	return int32(1)
}

// 计算表达式结果：出错返回0
func browserCalcExpr(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return float64(0)
	}
	expr := args[0].String()
	result, err := CalcExpr(expr)
	if err != nil {
		return float64(0)
	}
	return result
}

// 计算表达式的主函数
func CalcExpr(expr string) (float64, error) {
	r, err := calcExpression(expr)
	return r.InexactFloat64(), err
}

func main() {
	expr := "0.1 + 0.2"
	result, err := CalcExpr(expr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(expr, "=", result)
	}
	log.Println("Wasm loaded!")
	js.Global().Set("judgeExpr", js.FuncOf(browserJudgeExpr))
	js.Global().Set("calcExpr", js.FuncOf(browserCalcExpr))
	js.Global().Set("scienceCalc", js.FuncOf(ScienceCalc))
	select {}
}
