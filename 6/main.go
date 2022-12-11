package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const BUFFER_OFFSET = 14

func parseSignal(scanner *bufio.Scanner) string {
	scanner.Scan()
	line := scanner.Text()
	return line
}

func isPacketStart(bufferMap map[rune]int) bool {
	packetStart := true
	for _, count := range bufferMap {
		packetStart = packetStart && count <= 1
	}
	return packetStart
}

func findPacketStart(signal string) int {
	start := -1
	bufferMap := make(map[rune]int)
	for _, char := range signal[:BUFFER_OFFSET] {
		bufferMap[char] += 1
	}
	for i, char := range signal[BUFFER_OFFSET:] {
		if isPacketStart(bufferMap) {
			start = i + BUFFER_OFFSET
			break
		}
		bufferMap[char] += 1
		bufferMap[rune(signal[i])] -= 1
	}
	return start
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
	signal := parseSignal(scanner)
	packetStart := findPacketStart(signal)

	elapsed := time.Since(start)
	fmt.Println(packetStart)
	log.Printf("Time taken: %s", elapsed)
}
