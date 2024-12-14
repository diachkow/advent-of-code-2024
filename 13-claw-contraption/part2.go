package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ClawMachine struct {
	ButtonA, ButtonB, Prize struct{ X, Y int }
}

const MAX_BUTTON_PRESSES = 100

func main() {
	machines := ReadMachineConfigurations()

	totalTokens := 0

	for _, machine := range machines {
		Ax, Ay := machine.ButtonA.X, machine.ButtonA.Y
		Bx, By := machine.ButtonB.X, machine.ButtonB.Y
		Px, Py := machine.Prize.X, machine.Prize.Y

		// Using substitution method:
		//
		// { Ax*a + Bx*b = Px
		// { Ay*a + By*b = Py
		//
		// a = (Px - Bx*b) / Ax
		// b = (Py - Ay*a) / By
		//
		// b = (Py - Ay * ((Px - Bx*b) / Ax)) / By
		// b*By = Py - Ay * ((Px - Bx*b) / Ax)
		// b*By = Py - (Ay*Px - Ay*Bx*b) / Ax
		// b*By*Ax = Py*Ax - Ay*Px - Ay*Bx*b
		// b*By*Ax + Ay*Bx*b = Py*Ax - Ay*Px
		// b * (By*Ax + Ay*Bx) = Py*Ax - Ay*Px
		// b = (Py*Ax - Ay*Px) / (By*Ax + Ay*Bx)
		b := (Ax*Py - Ay*Px) / (Ax*By - Ay*Bx)
		a := (Px - b*Bx) / Ax

		if a*Ax+b*Bx == Px && a*Ay+b*By == Py {
			totalTokens += 3*a + b
		}
	}

	fmt.Printf("Total tokens = %v\n", totalTokens)
}

func ReadMachineConfigurations() []ClawMachine {
	file, err := os.Open("./machines.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	machines := make([]ClawMachine, 0)
	buffer := strings.Builder{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			} else {
				break
			}
		}

		if line == "\n" {
			continue
		}

		buffer.WriteString(line)

		if bufStr := buffer.String(); strings.Count(bufStr, "\n") == 3 {
			var machine ClawMachine
			_, err := fmt.Sscanf(
				bufStr,
				"Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d",
				&machine.ButtonA.X, &machine.ButtonA.Y,
				&machine.ButtonB.X, &machine.ButtonB.Y,
				&machine.Prize.X, &machine.Prize.Y,
			)
			if err != nil {
				panic(err)
			} else {
				machine.Prize.X += 10000000000000
				machine.Prize.Y += 10000000000000
				machines = append(machines, machine)
			}

			buffer.Reset()
		}
	}

	return machines
}
