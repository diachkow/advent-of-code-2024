package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

	output := []int{}

	instructionPtr := 0
	for instructionPtr < len(program) {
		opcode := program[instructionPtr]
		literalOperand := program[instructionPtr+1]

		switch opcode {
		case OPCODE_ADV:
			fmt.Printf("Executing `adv` operation with operand %v\n", literalOperand)
			RA = RA / IntPow(2, toComboOperand(literalOperand))

		case OPCODE_BXL:
			fmt.Printf("Executing `bxl` operation with operand %v\n", literalOperand)
			RB = RB ^ literalOperand

		case OPCODE_BST:
			fmt.Printf("Executing `bst` operation with operand %v\n", literalOperand)
			RB = toComboOperand(literalOperand) % 8

		case OPCODE_JNZ:
			fmt.Printf("Executing opcode `jnz` with operand %v (RA=%v)\n", literalOperand, RA)
			if RA != 0 {
				instructionPtr = literalOperand
				continue
			}

		case OPCODE_BXC:
			fmt.Printf("Executing opcode `bxc` with operand %v\n", literalOperand)
			RB = RB ^ RC

		case OPCODE_OUT:
			fmt.Printf("Executing opcode `out` with literalOperand %v\n", literalOperand)
			res := toComboOperand(literalOperand) % 8
			output = append(output, res)

		case OPCODE_BDV:
			fmt.Printf("Executing opcode `bdv` with operand %v\n", literalOperand)
			RB = int(float64(RA) / math.Pow(float64(toComboOperand(literalOperand)), 2))
			RB = RA / IntPow(2, toComboOperand(literalOperand))

		case OPCODE_CDV:
			fmt.Printf("Executing opcode `cdv` with literalOperand %v\n", literalOperand)
			RC = RA / IntPow(2, toComboOperand(literalOperand))
		}

		instructionPtr += 2
	}

	fmt.Printf("\n======================\n-> End of execution <-\n======================\n\n")

	fmt.Printf("A: %v\n", RA)
	fmt.Printf("B: %v\n", RB)
	fmt.Printf("C: %v\n", RC)

	outputStrings := make([]string, len(output))
	for i, num := range output {
		outputStrings[i] = strconv.Itoa(num)
	}

	fmt.Printf("\nOutput: %v\n", strings.Join(outputStrings, ","))
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
