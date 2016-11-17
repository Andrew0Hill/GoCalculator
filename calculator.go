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
	// Create two interfaces
	var right interface{}
	var left interface{}
	// Create two values to hold the Type of a Float64 and Int
	floatT := reflect.TypeOf(0.0)
	intT := reflect.TypeOf(0)
	if !operandStack.IsEmpty() {
		right = operandStack.Pop()
	}else {
		panic("Error: Operand stack is empty")
	}
	if !operandStack.IsEmpty() {
		left = operandStack.Pop()
	}else {
		panic("Error: Operand stack is empty")
	}
	// If types are the same
	if(reflect.TypeOf(right) == reflect.TypeOf(left)){
		// Print which type they are
		fmt.Println("Both types are of type " + reflect.TypeOf(right).String())
		// Compare to Float64 type
		if(reflect.TypeOf(right) == floatT){
			// Push result onto stack as floating point number
			switch op {
			case '+': operandStack.Push(left.(float64) + right.(float64))
			case '-': operandStack.Push(left.(float64) - right.(float64))
			case '*': operandStack.Push(left.(float64) * right.(float64))
			case '/': operandStack.Push(left.(float64) / right.(float64))
			default: panic("illegal operator")
			}
		}else if (reflect.TypeOf(right) == intT) {
			// Push result onto stack as integer number
			switch op {
			case '+': operandStack.Push(left.(int) + right.(int))
			case '-': operandStack.Push(left.(int) - right.(int))
			case '*':operandStack.Push(left.(int) * right.(int))
			case '/':
				// Check if division will result in a remainder. If so,
				// We need to use floating point.
				if((left.(int) % right.(int)) == 0){
					operandStack.Push(left.(int) / right.(int))
				}else {
					fmt.Println("Integer division will result in non-zero remainder,\nSwitching to floating point.")
					// Get value of operands
					rt := reflect.ValueOf(right)
					lt := reflect.ValueOf(left)
					// Get data from Value
					rt = reflect.Indirect(rt)
					lt = reflect.Indirect(lt)
					// Convert data to Float64 type
					rt = rt.Convert(floatT)
					lt = lt.Convert(floatT)
					// Push result onto stack as Float64
					operandStack.Push(lt.Float() / rt.Float())
				}

			default: panic("illegal operator")
			}
		}else{
			panic("Error: Invalid operand type.")
		}
	}else{
		// Get value of each interface
		rt := reflect.ValueOf(right)
		lt := reflect.ValueOf(left)
		// Get the data held by the Value
		rt = reflect.Indirect(rt)
		lt = reflect.Indirect(lt)
		// Convert the actual data to floating point type
		rt = rt.Convert(floatT)
		lt = lt.Convert(floatT)

		// Print out the types of each operand
		fmt.Println("Left operand is of type: " + reflect.TypeOf(left).String() + ".")
		fmt.Println("Right operand is of type: " + reflect.TypeOf(right).String() + ".")
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
				// Create a floating point variable and an integer variable.
				floatValue := float64(0)
				intValue := int(0)
				// Holds the power of 10 for floating point
				x := int(-1)
				// Set error checking flag.
				numberPrevious = true
				// Flag tells us if the number has a floating point part or not.
				floatPart := bool(false)
				if(spacePrevious){
					// If space occurs before operator, panic.
					panic("Illegal space in number.")
				}
				// Set spacePrevious to be false.
				spacePrevious = false
				for {
					// Check if multiple decimal points appear in the expression
					if (line[i] == '.'){
						if(floatPart){
							fmt.Println("Error: multiple decimal points")
							panic("Why did you do this to me")
						}
						floatPart = true
						floatValue = float64(intValue)

					}else if (floatPart){
						// Calculate the floating port part of the number
						floatValue = floatValue + (math.Pow10(x)*float64(line[i]-'0'))
						x--
					}else {
						// Otherwise, treat number as integer
						intValue = intValue * 10 + int(line[i] - '0')
					}
					i++
					// Check for end of the number
					if i == len(line) || (!('0' <= line[i] && line[i] <= '9') && (line[i] != '.')) {
						if(i != len(line) && line[i] == ' '){
							spacePrevious = true
						}
						break
					}
				}
				if(floatPart){
					// If number if floating point, push it to the stack as interface value
					pushVal := floatValue
					fmt.Println("Type of " + reflect.TypeOf(pushVal).String() + " is being pushed to the stack.")
					operandStack.Push(pushVal)
				}else {
					// If number is int, push it to the stack as interface value
					pushVal := intValue
					fmt.Println("Type of " + reflect.TypeOf(pushVal).String() + " is being pushed to the stack.")
					operandStack.Push(pushVal)
				}


			case '+', '-', '*', '/':
				// If operator occurs without a number before it, panic
				if(!numberPrevious){
					panic("Error: invalid operator.")
				}
				// If there is an insufficient number of operands, panic
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
				// Increment openCloseParens to keep track of parenthesis, and push onto the stack.
				openCloseParens++
				spacePrevious = false
				operatorStack.Push(line[i])
				i++
			case ')':
				// Apply the previous expression until we hit an open parenthesis
				openCloseParens--
				for !operatorStack.IsEmpty() && operatorStack.Top().(byte) != '(' {
					apply()
				}
				// Pop the open parenthesis from the stack.
				operatorStack.Pop()
				i++
			case ' ':
				// Error checking for spaces is handled inside the case statement for number input
				i++
			default:
				// Print error if character is illegal
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
	// If the operand stack isn't empty after we pop the result value, something went wrong
	if(!operandStack.IsEmpty()){
		panic("Error: Invalid expression")
	}
	// If the openCloseParens variable isn't set to 0 after we finish, The entered an expression
	// with mismatches parentheses.
	if(openCloseParens != 0){
		panic("Error: Mismatched Parens")
	}
	// Print out the type of the final value, and its result.
	fmt.Println("Final result value is of type: " + reflect.TypeOf(r).String() + ".")
	fmt.Println(r)
}