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

type Monkey struct {
	items          []int
	operation      func(old int) int
	test           func(n int) int
	itemsInspected int
}

func (monkey Monkey) performOperation(testDivisorProduct int) Monkey {
	monkey.items[0] = monkey.operation(monkey.items[0]) % testDivisorProduct
	monkey.itemsInspected++
	return monkey
}

func (monkey Monkey) performTest() int {
	return monkey.test(monkey.items[0])
}

func parseItems(line string) []int {
	itemsSplit := strings.Split(strings.TrimPrefix(line, "  Starting items: "), ", ")
	items := make([]int, 0)
	for _, itemStr := range itemsSplit {
		item, _ := strconv.Atoi(itemStr)
		items = append(items, item)
	}
	return items
}

func parseOperation(line string) func(old int) int {
	operationStr := strings.TrimPrefix(line, "  Operation: new = ")
	if strings.Contains(operationStr, "*") {
		multiplierStr := strings.Split(operationStr, " * ")[1]
		if multiplierStr == "old" {
			return func(old int) int { return old * old }
		}
		multiplier, _ := strconv.Atoi(multiplierStr)
		return func(old int) int { return old * multiplier }
	}
	addStr := strings.Split(operationStr, " + ")[1]
	if addStr == "old" {
		return func(old int) int { return old + old }
	}
	add, _ := strconv.Atoi(addStr)
	return func(old int) int { return old + add }
}

func parseTest(scanner *bufio.Scanner) (func(n int) int, int) {
	testStr := strings.TrimPrefix(scanner.Text(), "  Test: divisible by ")
	test, _ := strconv.Atoi(testStr)
	scanner.Scan()
	trueCaseStr := strings.TrimPrefix(scanner.Text(), "    If true: throw to monkey ")
	trueCase, _ := strconv.Atoi(trueCaseStr)
	scanner.Scan()
	falseCaseStr := strings.TrimPrefix(scanner.Text(), "    If false: throw to monkey ")
	falseCase, _ := strconv.Atoi(falseCaseStr)
	return func(n int) int {
			if n%test == 0 {
				return trueCase
			}
			return falseCase
		},
		test
}

func parseMonkeys(scanner *bufio.Scanner) ([]Monkey, int) {
	monkeys := make([]Monkey, 0)
	testDivisorProduct := 1
	for scanner.Scan() {
		scanner.Scan()
		items := parseItems(scanner.Text())
		scanner.Scan()
		operation := parseOperation(scanner.Text())
		scanner.Scan()
		test, testDivisor := parseTest(scanner)
		scanner.Scan()
		testDivisorProduct *= testDivisor
		monkeys = append(
			monkeys,
			Monkey{
				items:          items,
				operation:      operation,
				test:           test,
				itemsInspected: 0,
			},
		)
	}
	return monkeys, testDivisorProduct
}

func cutIndex(i int, items []int) (int, []int) {
	item := items[i]
	return item, append(items[:i], items[i+1:]...)
}

func throwToMonkey(
	monkeys []Monkey,
	fromMonkeyIndex int,
	toMonkeyIndex int,
) []Monkey {
	item := 0
	item, monkeys[fromMonkeyIndex].items = cutIndex(0, monkeys[fromMonkeyIndex].items)
	monkeys[toMonkeyIndex].items = append(monkeys[toMonkeyIndex].items, item)
	return monkeys
}

func performRounds(monkeys []Monkey, testDivisorProduct int, rounds int) []Monkey {
	for round := 0; round < rounds; round++ {
		for m := range monkeys {
			for len(monkeys[m].items) > 0 {
				monkeys[m] = monkeys[m].performOperation(testDivisorProduct)
				throwTo := monkeys[m].performTest()
				monkeys = throwToMonkey(monkeys, m, throwTo)
			}
		}
	}
	return monkeys
}

func calculateMonkeyBusinesLevel(monkeys []Monkey) int {
	monkeyItemsInspected := make([]int, 0)
	for _, monkey := range monkeys {
		monkeyItemsInspected = append(monkeyItemsInspected, monkey.itemsInspected)
	}
	sort.Ints(monkeyItemsInspected)
	return monkeyItemsInspected[len(monkeyItemsInspected)-1] * monkeyItemsInspected[len(monkeyItemsInspected)-2]
}

// way of getting test divisor sum has room
// room for improvement but it was low effort
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
	monkeys, testDivisorProduct := parseMonkeys(scanner)
	monkeys = performRounds(monkeys, testDivisorProduct, 10000)
	monkeyBusinessLevel := calculateMonkeyBusinesLevel(monkeys)

	elapsed := time.Since(start)
	fmt.Println(monkeyBusinessLevel)
	log.Printf("Time taken: %s", elapsed)
}
