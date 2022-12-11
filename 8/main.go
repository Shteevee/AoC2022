package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Point struct {
	x int
	y int
}

func parseCommands(scanner *bufio.Scanner) [][]int {
	treeRows := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		treeRow := make([]int, 0)
		for _, treeRune := range line {
			treeRow = append(treeRow, int(treeRune-48))
		}
		treeRows = append(treeRows, treeRow)
	}
	return treeRows
}

func indexToPoint(x int, y int) Point {
	return Point{
		x: x,
		y: y,
	}
}

func transposeGrid(treeGrid [][]int) [][]int {
	transposedGrid := make([][]int, len(treeGrid))
	for i := range transposedGrid {
		transposedGrid[i] = make([]int, len(treeGrid))
	}
	for i := range treeGrid {
		for j := range treeGrid[i] {
			transposedGrid[j][i] = treeGrid[i][j]
		}
	}
	return transposedGrid
}

func processRowRight(row []int, rowIndex int, visibleTrees map[Point]bool) map[Point]bool {
	currentTallestValue := row[0]
	for j, tree := range row[1 : len(row)-1] {
		if tree > currentTallestValue {
			visibleTrees[indexToPoint(rowIndex, j+1)] = true
			currentTallestValue = tree
		}
	}
	return visibleTrees
}

func processRowLeft(row []int, rowIndex int, visibleTrees map[Point]bool) map[Point]bool {
	currentTallestValue := row[len(row)-1]
	for j := len(row) - 2; j > 0; j-- {
		if row[j] > currentTallestValue {
			visibleTrees[indexToPoint(rowIndex, j)] = true
			currentTallestValue = row[j]
		}
	}
	return visibleTrees
}

func processColumnDown(column []int, columnIndex int, visibleTrees map[Point]bool) map[Point]bool {
	currentTallestValue := column[0]
	for i, tree := range column[1 : len(column)-1] {
		if tree > currentTallestValue {
			visibleTrees[indexToPoint(i+1, columnIndex)] = true
			currentTallestValue = tree
		}
	}
	return visibleTrees
}

func processColumnUp(column []int, columnIndex int, visibleTrees map[Point]bool) map[Point]bool {
	currentTallestValue := column[len(column)-1]
	for i := len(column) - 2; i > 0; i-- {
		if column[i] > currentTallestValue {
			visibleTrees[indexToPoint(i, columnIndex)] = true
			currentTallestValue = column[i]
		}
	}
	return visibleTrees
}

func findVisibleTrees(treeGrid [][]int) map[Point]bool {
	visibleTrees := make(map[Point]bool)
	for i, row := range treeGrid[1 : len(treeGrid)-1] {
		visibleTrees = processRowRight(row, i+1, visibleTrees)
		visibleTrees = processRowLeft(row, i+1, visibleTrees)
	}
	transposedGrid := transposeGrid(treeGrid)
	for j, column := range transposedGrid[1 : len(transposedGrid)-1] {
		visibleTrees = processColumnDown(column, j+1, visibleTrees)
		visibleTrees = processColumnUp(column, j+1, visibleTrees)
	}
	return visibleTrees
}

func countVisibleTree(visibleTrees map[Point]bool, treeGrid [][]int) int {
	return len(visibleTrees) + len(treeGrid)*2 + 2*(len(treeGrid)-2)
}

func isViewBlocked(start int, nextTree int) bool {
	return start <= nextTree
}

func findLeftViewScore(tree Point, treeGrid [][]int) int {
	view := 0
	for {
		view++
		if tree.y-view <= 0 || isViewBlocked(treeGrid[tree.x][tree.y], treeGrid[tree.x][tree.y-view]) {
			break
		}
	}
	return view
}

func findRightViewScore(tree Point, treeGrid [][]int) int {
	view := 0
	for {
		view++
		if tree.y+view >= len(treeGrid)-1 || isViewBlocked(treeGrid[tree.x][tree.y], treeGrid[tree.x][tree.y+view]) {
			break
		}
	}
	return view
}

func findUpViewScore(tree Point, treeGrid [][]int) int {
	view := 0
	for {
		view++
		if tree.x-view <= 0 || isViewBlocked(treeGrid[tree.x][tree.y], treeGrid[tree.x-view][tree.y]) {
			break
		}
	}
	return view
}

func findDownViewScore(tree Point, treeGrid [][]int) int {
	view := 0
	for {
		view++
		if tree.x+view >= len(treeGrid)-1 || isViewBlocked(treeGrid[tree.x][tree.y], treeGrid[tree.x+view][tree.y]) {
			break
		}
	}
	return view
}

func findScenicScore(tree Point, treeGrid [][]int) int {
	left := findLeftViewScore(tree, treeGrid)
	right := findRightViewScore(tree, treeGrid)
	up := findUpViewScore(tree, treeGrid)
	down := findDownViewScore(tree, treeGrid)
	return left * right * up * down
}

func findHighestScenicScore(treeGrid [][]int) int {
	highestScore := 0
	for i := range treeGrid {
		for j := range treeGrid[i] {
			score := findScenicScore(Point{x: i, y: j}, treeGrid)
			if score > highestScore {
				highestScore = score
			}
		}
	}
	return highestScore
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
	treeGrid := parseCommands(scanner)
	visibleTrees := findVisibleTrees(treeGrid)
	totalVisible := countVisibleTree(visibleTrees, treeGrid)
	highestScore := findHighestScenicScore(treeGrid)

	elapsed := time.Since(start)
	fmt.Println(totalVisible)
	fmt.Println(highestScore)
	log.Printf("Time taken: %s", elapsed)
}
