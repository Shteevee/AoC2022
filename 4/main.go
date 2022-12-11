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

type SectionAssignment struct {
	lower int
	upper int
}

type Pair[T, U any] struct {
	first  T
	second U
}

func parseSectionAssignment(assignment string) SectionAssignment {
	assignmentSplit := strings.Split(assignment, "-")
	lower, _ := strconv.Atoi(assignmentSplit[0])
	upper, _ := strconv.Atoi(assignmentSplit[1])
	return SectionAssignment{
		lower: lower,
		upper: upper,
	}
}

func parseSectionAssignmentPairs(scanner *bufio.Scanner) []Pair[SectionAssignment, SectionAssignment] {
	sectionAssignmentPairs := make([]Pair[SectionAssignment, SectionAssignment], 0)
	for scanner.Scan() {
		line := scanner.Text()
		splitAssignments := strings.Split(line, ",")
		sectionAssignmentPairs = append(sectionAssignmentPairs, Pair[SectionAssignment, SectionAssignment]{
			first:  parseSectionAssignment(splitAssignments[0]),
			second: parseSectionAssignment(splitAssignments[1]),
		})
	}
	return sectionAssignmentPairs
}

func isOverlapping(a SectionAssignment, b SectionAssignment) bool {
	return (a.upper >= b.lower) && (a.lower <= b.lower)
}

func countOverlapTotal(pairs []Pair[SectionAssignment, SectionAssignment]) int {
	total := 0
	for _, pair := range pairs {
		if isOverlapping(pair.first, pair.second) {
			total += 1
		} else if isOverlapping(pair.second, pair.first) {
			total += 1
		}
	}

	return total
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
	sectionAssignmentPairs := parseSectionAssignmentPairs(scanner)
	subRangeCount := countOverlapTotal(sectionAssignmentPairs)

	elapsed := time.Since(start)
	fmt.Println(subRangeCount)
	log.Printf("Time taken: %s", elapsed)
}
