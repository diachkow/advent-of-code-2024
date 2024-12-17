package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Registers
var RA, RB, RC int

const (
	OPCODE_ADV int = iota
	OPCODE_BXL
	OPCODE_BST
	OPCODE_JNZ
	OPCODE_BXC
	OPCODE_OUT
	OPCODE_BDV
	OPCODE_CDV
)

func main() {
	program := ReadDebuggerInfo()
	fmt.Printf("A: %v\n", RA)
	fmt.Printf("B: %v\n", RB)
	fmt.Printf("C: %v\n", RC)
	fmt.Printf("Program is %v\n", program)

	fmt.Printf("\n=======================\n-> Executing program <-\n=======================\n\n")

	var output []int
	a := 1
	for {
		RA = a
		output = []int{}

		instructionPtr := 0
		for instructionPtr < len(program) {
			opcode := program[instructionPtr]
			literalOperand := program[instructionPtr+1]

			switch opcode {
			case OPCODE_ADV:
				RA = RA / IntPow(2, toComboOperand(literalOperand))

			case OPCODE_BXL:
				RB = RB ^ literalOperand

			case OPCODE_BST:
				RB = toComboOperand(literalOperand) % 8

			case OPCODE_JNZ:
				if RA != 0 {
					instructionPtr = literalOperand
					continue
				}

			case OPCODE_BXC:
				RB = RB ^ RC

			case OPCODE_OUT:
				res := toComboOperand(literalOperand) % 8
				output = append(output, res)

			case OPCODE_BDV:
				RB = int(float64(RA) / math.Pow(float64(toComboOperand(literalOperand)), 2))
				RB = RA / IntPow(2, toComboOperand(literalOperand))

			case OPCODE_CDV:
				RC = RA / IntPow(2, toComboOperand(literalOperand))
			}

			instructionPtr += 2
		}

		if slices.Equal(output, program) {
			fmt.Printf("Correct A value is %v\n", a)
			break
		}

		if slices.Equal(output, program[len(program)-len(output):]) {
			a *= 8
		} else {
			a++
		}
	}
}

func toComboOperand(literal int) int {
	if literal >= 0 && literal <= 3 {
		return literal
	} else if literal == 4 {
		return RA
	} else if literal == 5 {
		return RB
	} else if literal == 6 {
		return RC
	} else {
		panic("Unexpected value for combo operand!")
	}
}

func ReadDebuggerInfo() []int {
	file, err := os.Open("./debugger.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	program := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Register A:") {
			_, parseErr := fmt.Sscanf(line, "Register A: %d", &RA)
			if parseErr != nil {
				panic(parseErr)
			}
		} else if strings.Contains(line, "Register B:") {
			_, parseErr := fmt.Sscanf(line, "Register B: %d", &RB)
			if parseErr != nil {
				panic(parseErr)
			}
		} else if strings.Contains(line, "Register C:") {
			_, parseErr := fmt.Sscanf(line, "Register C: %d", &RC)
			if parseErr != nil {
				panic(parseErr)
			}
		} else if strings.Contains(line, "Program:") {
			var instructionString string
			_, parseErr := fmt.Sscanf(line, "Program: %s", &instructionString)
			if parseErr != nil {
				panic(parseErr)
			}

			for _, str := range strings.Split(instructionString, ",") {
				num, err := strconv.Atoi(str)
				if err != nil {
					panic(err)
				}
				program = append(program, num)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return program
}

func IntPow(base, exp int) int {
	result := 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}
