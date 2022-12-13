package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Pair[T, U any] struct {
	first  T
	second U
}

type ElfNumber struct {
	numType  string
	value    int
	children []ElfNumber
}

func splitElfNumber(line string) []string {
	currentBrackets := 0
	indexesToSplit := make([]int, 0)
	line = line[1 : len(line)-1]
	for i, char := range line {
		if char == '[' {
			currentBrackets++
		} else if char == ']' {
			currentBrackets--
		} else if char == ',' && currentBrackets == 0 {
			indexesToSplit = append(indexesToSplit, i)
		}
	}
	indexesToSplit = append(indexesToSplit, len(line))
	previousSplit := 0
	children := make([]string, 0)
	for _, i := range indexesToSplit {
		children = append(children, line[previousSplit:i])
		previousSplit = i + 1
	}
	return children
}

func parseElfNumber(line string) ElfNumber {
	if line == "[]" {
		return ElfNumber{
			numType: "list",
		}
	}
	if !strings.Contains(line, "[") {
		value, _ := strconv.Atoi(line)
		return ElfNumber{
			numType: "int",
			value:   value,
		}
	}
	children := make([]ElfNumber, 0)
	childElfNumbers := splitElfNumber(line)
	for _, child := range childElfNumbers {
		children = append(children, parseElfNumber(child))
	}
	return ElfNumber{
		numType:  "list",
		children: children,
	}
}

func parseElfNumberPairs(scanner *bufio.Scanner) []Pair[ElfNumber, ElfNumber] {
	elfNumberPairs := make([]Pair[ElfNumber, ElfNumber], 0)
	for scanner.Scan() {
		number1 := parseElfNumber(scanner.Text())
		scanner.Scan()
		number2 := parseElfNumber(scanner.Text())
		scanner.Scan()
		elfNumberPairs = append(
			elfNumberPairs,
			Pair[ElfNumber, ElfNumber]{first: number1, second: number2},
		)
	}
	return elfNumberPairs
}

func compareElfNum(left ElfNumber, right ElfNumber) int {
	if left.numType == "int" && right.numType == "int" {
		return left.value - right.value
	}
	if left.numType == "int" && right.numType == "list" {
		return compareElfNum(
			ElfNumber{numType: "list", children: []ElfNumber{left}},
			right,
		)
	}
	if right.numType == "int" && left.numType == "list" {
		return compareElfNum(
			left,
			ElfNumber{numType: "list", children: []ElfNumber{right}},
		)
	}
	for i := 0; i < len(left.children) && i < len(right.children); i++ {
		comparison := compareElfNum(left.children[i], right.children[i])
		if comparison != 0 {
			return comparison
		}
	}
	return len(left.children) - len(right.children)
}

func findPairsInCorrectOrder(elfNumberPairs []Pair[ElfNumber, ElfNumber]) int {
	correctIndexSum := 0
	for i, pair := range elfNumberPairs {
		if compareElfNum(pair.first, pair.second) < 0 {
			correctIndexSum += i + 1
		}
	}
	return correctIndexSum
}

func breakPairsAndAddDividerPackets(elfNumberPairs []Pair[ElfNumber, ElfNumber]) []ElfNumber {
	elfNumbers := make([]ElfNumber, 0)
	for _, pair := range elfNumberPairs {
		elfNumbers = append(elfNumbers, pair.first)
		elfNumbers = append(elfNumbers, pair.second)
	}
	elfNumbers = append(elfNumbers, parseElfNumber("[[2]]"))
	elfNumbers = append(elfNumbers, parseElfNumber("[[6]]"))
	return elfNumbers
}

// assumes that there are no packets matching divider
func calculateDividerIndexProduct(elfNumbers []ElfNumber) int {
	product := 1
	divider1 := parseElfNumber("[[2]]")
	divider2 := parseElfNumber("[[6]]")
	for i, num := range elfNumbers {
		if compareElfNum(num, divider1) == 0 {
			product *= i + 1
		}
		if compareElfNum(num, divider2) == 0 {
			product *= i + 1
		}
	}
	return product
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
	elfNumberPairs := parseElfNumberPairs(scanner)
	correctIndexSum := findPairsInCorrectOrder(elfNumberPairs)
	elfNumbers := breakPairsAndAddDividerPackets(elfNumberPairs)
	sort.Slice(elfNumbers, func(i, j int) bool {
		return compareElfNum(elfNumbers[i], elfNumbers[j]) < 0
	})
	dividerIndexProduct := calculateDividerIndexProduct(elfNumbers)

	elapsed := time.Since(start)
	fmt.Println("Correct index sum:", correctIndexSum)
	fmt.Println("Divider index product:", dividerIndexProduct)
	log.Printf("Time taken: %s", elapsed)
}
