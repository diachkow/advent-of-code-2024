package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var OPERATORS = []string{"+", "*", "||"}

type Equation struct {
	TestValue int
	Operands  []int
}

func (e Equation) String() string {
	return fmt.Sprintf("%v: %v", e.TestValue, e.Operands)
}

func TestEquation(eq Equation) bool {
	if len(eq.Operands) == 0 {
		return false
	}
	if len(eq.Operands) == 1 {
		return eq.Operands[0] == eq.TestValue
	}

	for _, operators := range Product(Repeat(OPERATORS, len(eq.Operands)-1)) {
		if eq.TestValue == EvaluateEquation(eq.Operands, operators) {
			return true
		}
	}
	return false
}

func EvaluateEquation(operands []int, operators []string) int {
	if len(operands) != len(operators)+1 {
		panic(fmt.Errorf("Invalid operands/operators ration: %v/%v", len(operands), len(operators)))
	}

	result := operands[0]
	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case "+":
			result += operands[i+1]
		case "*":
			result *= operands[i+1]
		case "||":
			result = MergeNumbers(result, operands[i+1])
		}
	}

	return result
}

func MergeNumbers(x, y int) int {
	return StringToInteger(strconv.Itoa(x) + strconv.Itoa(y))
}

func RunEquationsTest(eqs []Equation) int {
	var total int
	for _, eq := range eqs {
		if TestEquation(eq) {
			total += eq.TestValue
		}
	}
	return total
}

func main() {
	equations := ReadEquations()
	// equations := []Equation{Equation{TestValue: 7290, Operands: []int{6, 8, 6, 15}}}
	total := RunEquationsTest(equations)
	fmt.Printf("Total value of all calibrated equations is: %v\n", total)
}

func ReadEquations() []Equation {
	file, err := os.Open("./equations.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var res []Equation
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ": ")
		if len(splitted) != 2 {
			panic("Invalid equation string")
		}
		testValString, operandsString := splitted[0], splitted[1]

		testVal := StringToInteger(testValString)
		operandsAsStrings := strings.Split(operandsString, " ")

		operands := make([]int, len(operandsAsStrings))
		for i, ops := range operandsAsStrings {
			operands[i] = StringToInteger(ops)
		}

		res = append(res, Equation{testVal, operands})
	}

	return res
}

func StringToInteger(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func Repeat[T comparable](sets []T, times int) [][]T {
	result := make([][]T, times)
	for i := 0; i < times; i++ {
		result[i] = make([]T, len(sets))
		copy(result[i], sets)
	}
	return result
}

func Product[T comparable](sets [][]T) [][]T {
	if len(sets) == 0 {
		return [][]T{{}}
	}

	firstSet, restSets := sets[0], sets[1:]
	restProduct := Product(restSets)

	results := make([][]T, 0)
	for _, el := range firstSet {
		for _, comb := range restProduct {
			newComb := append([]T{el}, comb...)
			results = append(results, newComb)
		}
	}

	return results
}
