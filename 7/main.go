package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const TOTAL_DISK_SPACE = 70000000
const TARGET_SPACE = 30000000

func parseCommands(scanner *bufio.Scanner) []string {
	commands := make([]string, 0)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	return commands
}

func isFile(command string) bool {
	return !strings.HasPrefix(command, "$") && !strings.HasPrefix(command, "dir")
}

func isDirChange(command string) bool {
	return strings.HasPrefix(command, "$ cd")
}

func isGoPreviousDir(command string) bool {
	return command == "$ cd .."
}

func parseFileBytes(command string) int {
	commandSplit := strings.Split(command, " ")
	byteValue, _ := strconv.Atoi(commandSplit[0])
	return byteValue
}

func countDirBytes(commands []string) map[string]int {
	currentPath := ""
	currentDirs := []string{"/"}
	dirBytes := make(map[string]int)
	for _, command := range commands[1:] {
		if isGoPreviousDir(command) {
			lastDirIndex := strings.LastIndex(currentPath, "/")
			currentPath = currentPath[:lastDirIndex]
			currentDirs = currentDirs[:len(currentDirs)-1]
		} else if isDirChange(command) {
			dir := strings.TrimPrefix(command, "$ cd ")
			currentPath += "/" + dir
			currentDirs = append(currentDirs, currentPath)
		} else if isFile(command) {
			bytes := parseFileBytes(command)
			for _, dir := range currentDirs {
				dirBytes[dir] += bytes
			}
		}
	}
	return dirBytes
}

func calculateTargetByteSum(dirByteCounts map[string]int) int {
	total := 0
	for _, byteCount := range dirByteCounts {
		if byteCount <= 100000 {
			total += byteCount
		}
	}
	return total
}

func findDirToDelete(dirByteCounts map[string]int) int {
	remainingSpace := TOTAL_DISK_SPACE - dirByteCounts["/"]
	targetSpace := TARGET_SPACE - remainingSpace
	currentMin := math.MaxInt
	bytesDeleted := 0
	for _, bytes := range dirByteCounts {
		if bytes >= targetSpace && bytes-targetSpace < currentMin {
			currentMin = bytes - targetSpace
			bytesDeleted = bytes
		}
	}
	return bytesDeleted
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
	commands := parseCommands(scanner)
	dirByteCounts := countDirBytes(commands)
	targetSum := calculateTargetByteSum(dirByteCounts)
	spaceToDelete := findDirToDelete(dirByteCounts)

	elapsed := time.Since(start)
	fmt.Println(targetSum)
	fmt.Println(spaceToDelete)
	log.Printf("Time taken: %s", elapsed)
}
