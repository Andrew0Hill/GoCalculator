package main

import (
	"bufio"
	"fmt"
	"os"
	"stack"
	"math"
	"reflect"
)

var operatorStack = stack.NewStack()
var operandStack = stack.NewStack()

func precedence(op byte) uint8 {
	switch op {
		case '(', ')': return 0
		case '+', '-': return 1
		case '*', '/': return 2
		default: panic("illegal operator")
	}
}

func apply() {
	op := operatorStack.Pop().(byte)
	var right interface{}
	var left interface{}
	floatT := reflect.TypeOf(0.0)
	intT := reflect.TypeOf(0)
	if !operandStack.IsEmpty() {
		right = operandStack.Pop()
	}
	if !operandStack.IsEmpty() {
		left = operandStack.Pop()
	}
	if(reflect.TypeOf(right) == reflect.TypeOf(left)){
		fmt.Println("Both types are of type " + reflect.TypeOf(right).String())
		if(reflect.TypeOf(right) == floatT){
			switch op {
			case '+': operandStack.Push(left.(float64) + right.(float64))
			case '-': operandStack.Push(left.(float64) - right.(float64))
			case '*': operandStack.Push(left.(float64) * right.(float64))
			case '/': operandStack.Push(left.(float64) / right.(float64))
			default: panic("illegal operator")
			}
		}else if (reflect.TypeOf(right) == intT) {
			switch op {
			case '+': operandStack.Push(left.(int) + right.(int))
			case '-': operandStack.Push(left.(int) - right.(int))
			case '*': operandStack.Push(left.(int) * right.(int))
			case '/': operandStack.Push(left.(int) / right.(int))
			default: panic("illegal operator")
			}
		}else{
			panic("Error: Invalid operand type.")
		}
	}else{
		fmt.Println("Left operand is of type: " + reflect.TypeOf(left).String() + ".")
		fmt.Println("Right operand is of type: " + reflect.TypeOf(right).String() + ".")
		rt := reflect.ValueOf(right)
		rt = reflect.Indirect(rt)
		rt = rt.Convert(floatT)
		lt := reflect.ValueOf(left)
		lt = reflect.Indirect(lt)
		lt = lt.Convert(floatT)
		switch op {
		case '+': operandStack.Push(lt.Float() + rt.Float())
		case '-': operandStack.Push(lt.Float() - rt.Float())
		case '*': operandStack.Push(lt.Float() * rt.Float())
		case '/': operandStack.Push(lt.Float() / rt.Float())
		default: panic("illegal operator")
		}
	}


}

func main() {
	// Read a from Stdin.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	openCloseParens := 0
	spacePrevious := bool(false)
	numberPrevious := bool(false)
	for i := 0; i < len(line); {
		switch line[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':

				floatValue := float64(0)
				intValue := int(0)
				x := int(-1)
				numberPrevious = true
				floatPart := bool(false)
				if(spacePrevious){
					panic("Illegal space in number.")
				}
				spacePrevious = false
				for {
					if (line[i] == '.'){
						floatPart = true
						floatValue = float64(intValue)

					}else if (floatPart){
						floatValue = floatValue + (math.Pow10(x)*float64(line[i]-'0'))
						x--
					}else {
						intValue = intValue * 10 + int(line[i] - '0')
					}
					i++
					if i == len(line) || (!('0' <= line[i] && line[i] <= '9') && (line[i] != '.')) {
						if(i != len(line) && line[i] == ' '){
							spacePrevious = true
						}
						break
					}
				}
				if(floatPart){
					pushVal := floatValue
					fmt.Println("Type of " + reflect.TypeOf(pushVal).String() + " is being pushed to the stack.")
					operandStack.Push(pushVal)
				}else {
					pushVal := intValue
					fmt.Println("Type of " + reflect.TypeOf(pushVal).String() + " is being pushed to the stack.")
					operandStack.Push(pushVal)
				}


			case '+', '-', '*', '/':
				if(!numberPrevious){
					panic("Error: invalid operator.")
				}
				if(operandStack.IsEmpty()){
					panic("Error: Operator must be applied to two operands.")
				}
				numberPrevious = false
				spacePrevious = false
				for !operatorStack.IsEmpty() && precedence(operatorStack.Top().(byte)) >= precedence(line[i]) {
					apply()
				}
				operatorStack.Push(line[i])
				i++
			case '(':
				openCloseParens++
				spacePrevious = false
				operatorStack.Push(line[i])
				i++
			case ')':
				openCloseParens--
				for !operatorStack.IsEmpty() && operatorStack.Top().(byte) != '(' {
					apply()
				}
				operatorStack.Pop()
				i++
			case ' ':
				i++
			default:
				println("Invalid Character '" + string(line[i]) + "'. Please try again.")
				panic("illegal character")
		}
	}
	for !operatorStack.IsEmpty() {
		if operatorStack.Top().(byte) == '(' {
			operatorStack.Pop()
		}else {
			apply()
		}
	}
	r := operandStack.Pop()
	if(!operandStack.IsEmpty()){
		panic("Error: Invalid expression")
	}
	if(openCloseParens != 0){
		panic("Error: Mismatched Parens")
	}
	fmt.Println("Final result value is of type: " + reflect.TypeOf(r).String() + ".")
	fmt.Println(r)
}