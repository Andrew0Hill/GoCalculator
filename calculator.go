package main

import (
		"bufio"
		"fmt"
		"os"
		"stack"

	"math"
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
	right := float64(-1)
	left := float64(-1)
	if !operandStack.IsEmpty() {
		right = operandStack.Pop().(float64)
	}
	if !operandStack.IsEmpty() {
		left = operandStack.Pop().(float64)
	}
	switch op {
		case '+': operandStack.Push(left + right)
		case '-': operandStack.Push(left - right)
		case '*': operandStack.Push(left * right)
		case '/': operandStack.Push(left / right)
		default: panic("illegal operator")
	}
}

func main() {
	// Read a from Stdin.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	openParen := false
	spacePrevious := bool(false)
	//validInputs := [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9','+', '-', '*', '/','(',')','.'}
	//
	for i := 0; i < len(line); {
		switch line[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

				v := float64(0)
				x := int(-1)
				floatPart := bool(false)
				if(spacePrevious){
					panic("Illegal space in number.")
				}
				spacePrevious = false
				for {
					if (line[i] == '.'){
						floatPart = true

					}else if (floatPart){
						v = v + (math.Pow10(x)*float64(line[i]-'0'))
						x--
					}else {
						v = v * 10 + float64(line[i] - '0')
					}
					i++
					if i == len(line) || (!('0' <= line[i] && line[i] <= '9') && (line[i] != '.')) {
						if(i != len(line) && line[i] == ' '){
							spacePrevious = true
						}
						break
					}
				}
				operandStack.Push(v)
			case '+', '-', '*', '/':
				spacePrevious = false
				for !operatorStack.IsEmpty() && precedence(operatorStack.Top().(byte)) >= precedence(line[i]) {
					apply()
				}
				operatorStack.Push(line[i])
				i++
			case '(':
				openParen = true
				spacePrevious = false
				operatorStack.Push(line[i])
				i++
			case ')':
				if(!openParen){
					println("Error: Close parenthesis appears with no Open parenthesis")
					panic("illegal character")
				}
				openParen = false

				for !operatorStack.IsEmpty() && operatorStack.Top().(byte) != '(' {
					apply()
				}
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
	r := operandStack.Pop().(float64)
	fmt.Println(r)
}