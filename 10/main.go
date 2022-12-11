package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const ADD_CYCLES = 2
const NOOP_CYCLES = 1

type Instruction struct {
	value  int
	cycles int
}

func buildAdd(value int) Instruction {
	return Instruction{value: value, cycles: ADD_CYCLES}
}

func buildNoop() Instruction {
	return Instruction{value: 0, cycles: NOOP_CYCLES}
}

func parseInstructions(scanner *bufio.Scanner) []Instruction {
	instructions := make([]Instruction, 0)
	for scanner.Scan() {
		line := scanner.Text()
		instructionSplit := strings.Split(line, " ")
		if len(instructionSplit) == 1 {
			instructions = append(instructions, buildNoop())
		} else {
			value, _ := strconv.Atoi(instructionSplit[1])
			instructions = append(instructions, buildAdd(value))
		}
	}
	return instructions
}

func popQueue(instructions []Instruction) (Instruction, []Instruction) {
	var currentInstruction Instruction
	if len(instructions) > 0 {
		currentInstruction = instructions[0]
		instructions = instructions[1:]
	} else {
		currentInstruction = buildNoop()
	}
	return currentInstruction, instructions
}

func runCycles(instructions []Instruction, cycles int) []int {
	x := 1
	cycleXs := make([]int, 0)
	currentInstruction := instructions[0]
	instructions = instructions[1:]
	currentCycle := 1
	for currentCycle < cycles+1 {
		cycleXs = append(cycleXs, x)
		currentInstruction.cycles--
		if currentInstruction.cycles == 0 {
			x += currentInstruction.value
			currentInstruction, instructions = popQueue(instructions)
		}
		cycles--
	}
	return cycleXs
}

func sumInterestingSignalStrengths(instructions []Instruction) int {
	interestingCycles := []int{20, 60, 100, 140, 180, 220}
	total := 0
	cycleXs := runCycles(instructions, 220)
	for _, cycle := range interestingCycles {
		total += cycle * cycleXs[cycle-1]
	}
	return total
}

func createScreen() [][]string {
	screen := make([][]string, 6)
	for i := range screen {
		screen[i] = make([]string, 40)
		for j := range screen[i] {
			screen[i][j] = "."
		}
	}
	return screen
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return 0 - n
}

func inSpriteRange(i int, spriteCentre int) bool {
	return abs(spriteCentre-i) <= 1
}

func displaySprite(cycleXs []int, screen [][]string) [][]string {
	currentCycle := 0
	for i := range screen {
		for j := range screen[i] {
			if inSpriteRange(j, cycleXs[currentCycle]) {
				screen[i][j] = "#"
			}
			currentCycle++
		}
	}
	return screen
}

func displayCRT(cycleXs []int) {
	screen := createScreen()
	screen = displaySprite(cycleXs, screen)
	for _, row := range screen {
		fmt.Println(row)
	}
}

func main() {
	start := time.Now()
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	instructions := parseInstructions(scanner)
	interestingSignalSum := sumInterestingSignalStrengths(instructions)
	cycleXs := runCycles(instructions, 240)
	displayCRT(cycleXs)

	elapsed := time.Since(start)
	fmt.Println(interestingSignalSum)
	log.Printf("Time taken: %s", elapsed)
}
