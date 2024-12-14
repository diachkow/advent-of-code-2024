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

	solvableMachines := make(map[ClawMachine]bool)
	totalTokens := 0

	for _, machine := range machines {
		type solution struct{ ap, bp, tokens int }
		const MAX_TOKENS = 99999
		s := solution{MAX_BUTTON_PRESSES, MAX_BUTTON_PRESSES, MAX_TOKENS}

		for ap := 1; ap <= MAX_BUTTON_PRESSES; ap++ {
			for bp := 1; bp <= MAX_BUTTON_PRESSES; bp++ {
				xPoints := machine.ButtonA.X*ap + machine.ButtonB.X*bp
				yPoints := machine.ButtonA.Y*ap + machine.ButtonB.Y*bp
				if xPoints == machine.Prize.X && yPoints == machine.Prize.Y {
					tokens := 3*ap + bp

					if tokens <= s.tokens {
						s = solution{ap, bp, tokens}
					}
				}
			}
		}

		if s.tokens != MAX_TOKENS {
			solvableMachines[machine] = true
			totalTokens += s.tokens
		}
	}

	fmt.Printf("Number of solvable machines: %v\n", len(solvableMachines))
	fmt.Printf("Total tokens is %v\n", totalTokens)
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
				fmt.Printf("Read machine: %v\n", machine)
				machines = append(machines, machine)
			}

			buffer.Reset()
		}
	}

	return machines
}
