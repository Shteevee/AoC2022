package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func parseRpsRounds(scanner *bufio.Scanner) []string {
	rpsRounds := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		rpsRounds = append(rpsRounds, line)
	}
	return rpsRounds
}

func createRoundScoring() map[string]int {
	roundScoring := make(map[string]int)
	roundScoring["A X"] = 3
	roundScoring["A Y"] = 4
	roundScoring["A Z"] = 8
	roundScoring["B X"] = 1
	roundScoring["B Y"] = 5
	roundScoring["B Z"] = 9
	roundScoring["C X"] = 2
	roundScoring["C Y"] = 6
	roundScoring["C Z"] = 7
	return roundScoring
}

func calculateTotalScore(rpsRounds []string, roundScoring map[string]int) int {
	score := 0
	for _, round := range rpsRounds {
		score += roundScoring[round]
	}
	return score
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
	rpsRounds := parseRpsRounds(scanner)
	roundScoring := createRoundScoring()
	totalScore := calculateTotalScore(rpsRounds, roundScoring)

	elapsed := time.Since(start)
	fmt.Println(totalScore)
	log.Printf("Time taken: %s", elapsed)
}
