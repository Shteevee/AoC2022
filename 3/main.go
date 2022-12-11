package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func parseRucksackGroups(scanner *bufio.Scanner) [][]string {
	rucksackGroups := make([][]string, 0)
	rucksackGroup := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		rucksackGroup = append(rucksackGroup, line)
		if len(rucksackGroup) == 3 {
			rucksackGroups = append(rucksackGroups, rucksackGroup)
			rucksackGroup = make([]string, 0)
		}
	}
	return rucksackGroups
}

func createRucksackSet(rucksack string) map[rune]struct{} {
	ruckSackSet := make(map[rune]struct{})
	for _, char := range rucksack {
		ruckSackSet[char] = struct{}{}
	}
	return ruckSackSet
}

// assumes there's only one occurence of a char that
// matches all three
func findGroupCommonType(rucksackGroup []string) rune {
	commonChar := ' '
	firstSackSet := createRucksackSet(rucksackGroup[0])
	secondSackSet := createRucksackSet(rucksackGroup[1])
	for _, char := range rucksackGroup[2] {
		_, firstFound := firstSackSet[char]
		_, secondFound := secondSackSet[char]
		if firstFound && secondFound {
			commonChar = char
		}
	}
	return commonChar
}

func findCommonTypes(rucksackGroups [][]string) []rune {
	commonTypes := make([]rune, 0)
	for _, rucksackGroup := range rucksackGroups {
		commonTypes = append(commonTypes, findGroupCommonType(rucksackGroup))
	}
	return commonTypes
}

func commonTypesToPriorities(commonTypes []rune) []int {
	priorities := make([]int, 0)
	for _, commonType := range commonTypes {
		if commonType >= 97 {
			priorities = append(priorities, int(commonType)-96)
		} else {
			priorities = append(priorities, int(commonType)-38)
		}
	}
	return priorities
}

func sumCommonTypePriorites(commonTypes []rune) int {
	sum := 0
	priorities := commonTypesToPriorities(commonTypes)
	for _, priority := range priorities {
		sum += priority
	}
	return sum
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
	rucksacksGroups := parseRucksackGroups(scanner)
	commonTypes := findCommonTypes(rucksacksGroups)
	prioritySum := sumCommonTypePriorites(commonTypes)

	elapsed := time.Since(start)
	fmt.Println(prioritySum)
	log.Printf("Time taken: %s", elapsed)
}
