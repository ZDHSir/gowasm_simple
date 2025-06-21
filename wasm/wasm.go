package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/dengsgo/math-engine/engine"
)

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
	r, err := engine.ParseAndExec(expr)
	return r, err
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
	select {}
}
