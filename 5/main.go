package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const BOX_WIDTH = 4

type Instruction struct {
	quantity    int
	origin      int
	destination int
}

type stack []string

func (s stack) Push(v string) stack {
	return append(s, v)
}

func (s stack) PushMultiple(v []string) stack {
	return append(s, v...)
}

// don't care if stack is empty
func (s stack) Pop() (stack, string) {
	length := len(s)
	return s[:length-1], s[length-1]
}

func (s stack) PopMultiple(quantity int) (stack, []string) {
	length := len(s)
	return s[:length-quantity], s[length-quantity:]
}

func (s stack) isEmpty() bool {
	return len(s) == 0
}

func reverseStacks(stacks []stack) []stack {
	reversedStacks := make([]stack, len(stacks))
	for i, stack := range stacks {
		for !stack.isEmpty() {
			newStack, value := stack.Pop()
			stack = newStack
			reversedStacks[i] = reversedStacks[i].Push(value)
		}
	}
	return reversedStacks
}

func parseBoxStacks(scanner *bufio.Scanner) ([]stack, []Instruction) {
	scanner.Scan()
	line := scanner.Text()
	numOfStacks := ((len(line) - 3) / BOX_WIDTH) + 1
	reversedStacks := make([]stack, numOfStacks)
	for line != "" {
		for i := range reversedStacks {
			value := string(line[i*BOX_WIDTH+1])
			if value != " " && !unicode.IsNumber([]rune(value)[0]) {
				reversedStacks[i] = reversedStacks[i].Push(value)
			}
		}
		scanner.Scan()
		line = scanner.Text()
	}
	instructions := make([]Instruction, 0)
	for scanner.Scan() {
		line = scanner.Text()
		trimmedInstruction := strings.TrimPrefix(line, "move ")
		trimmedInstruction = strings.ReplaceAll(trimmedInstruction, " from ", ",")
		trimmedInstruction = strings.ReplaceAll(trimmedInstruction, " to ", ",")
		splitInstruction := strings.Split(trimmedInstruction, ",")
		quantity, _ := strconv.Atoi(splitInstruction[0])
		origin, _ := strconv.Atoi(splitInstruction[1])
		destination, _ := strconv.Atoi(splitInstruction[2])
		instructions = append(instructions, Instruction{quantity: quantity, origin: origin - 1, destination: destination - 1})
	}
	boxStacks := reverseStacks(reversedStacks)
	return boxStacks, instructions
}

func performInstructions(boxStacks []stack, instructions []Instruction) []stack {
	for _, instruction := range instructions {
		poppedStack, value := boxStacks[instruction.origin].PopMultiple(instruction.quantity)
		pushedStack := boxStacks[instruction.destination].PushMultiple(value)
		boxStacks[instruction.origin] = poppedStack
		boxStacks[instruction.destination] = pushedStack
	}
	return boxStacks
}

func readTopBoxes(stacks []stack) string {
	top := ""
	for _, stack := range stacks {
		_, value := stack.Pop()
		top = top + value
	}
	return top
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
	boxStacks, instructions := parseBoxStacks(scanner)
	boxStacks = performInstructions(boxStacks, instructions)
	topBoxes := readTopBoxes(boxStacks)

	elapsed := time.Since(start)
	fmt.Println(topBoxes)
	log.Printf("Time taken: %s", elapsed)
}
