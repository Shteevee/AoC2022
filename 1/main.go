package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

func parseElfCalories(scanner *bufio.Scanner) [][]int {
	allElfCalories := make([][]int, 0)
	elfCalories := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			allElfCalories = append(allElfCalories, elfCalories)
			elfCalories = make([]int, 0)
		} else {
			calories, _ := strconv.Atoi(line)
			elfCalories = append(elfCalories, calories)
		}
	}

	return allElfCalories
}

func sumElfCalories(allElfCalories [][]int) []int {
	totals := make([]int, 0)
	for _, elf := range allElfCalories {
		total := 0
		for _, calorie := range elf {
			total += calorie
		}
		totals = append(totals, total)
	}

	return totals
}

func findTopThree(elfCalorieTotals []int) []int {
	sort.Ints(elfCalorieTotals)
	return elfCalorieTotals[len(elfCalorieTotals)-3:]
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
	allElfCalories := parseElfCalories(scanner)
	elfCalorieTotals := sumElfCalories(allElfCalories)
	topThreeCalorieTotals := findTopThree(elfCalorieTotals)

	elapsed := time.Since(start)
	fmt.Println(topThreeCalorieTotals)
	fmt.Println(topThreeCalorieTotals[0] + topThreeCalorieTotals[1] + topThreeCalorieTotals[2])

	log.Printf("Time taken: %s", elapsed)
}
