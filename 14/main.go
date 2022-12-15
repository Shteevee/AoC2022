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

type Point struct {
	x int
	y int
}

func createPointRange(a Point, b Point) []Point {
	points := make([]Point, 0)
	if a.x < b.x {
		for i := a.x + 1; i <= b.x; i++ {
			points = append(points, Point{x: i, y: a.y})
		}
	}
	if a.x > b.x {
		for i := a.x - 1; i >= b.x; i-- {
			points = append(points, Point{x: i, y: a.y})
		}
	}
	if a.y < b.y {
		for i := a.y + 1; i <= b.y; i++ {
			points = append(points, Point{x: a.x, y: i})
		}
	}
	if a.y > b.y {
		for i := a.y - 1; i >= b.y; i-- {
			points = append(points, Point{x: a.x, y: i})
		}
	}
	return points
}

func parsePoint(coord string) Point {
	splitCoord := strings.Split(coord, ",")
	x, _ := strconv.Atoi(splitCoord[0])
	y, _ := strconv.Atoi(splitCoord[1])
	return Point{x: x, y: y}
}

func parseRockPath(line string) []Point {
	path := make([]Point, 0)
	splitCoords := strings.Split(line, " -> ")
	prevPoint := parsePoint(splitCoords[0])
	path = append(path, prevPoint)
	for _, coord := range splitCoords[1:] {
		point := parsePoint(coord)
		path = append(path, createPointRange(prevPoint, point)...)
		prevPoint = point
	}
	return path
}

func parseRockPaths(scanner *bufio.Scanner) [][]Point {
	rockPaths := make([][]Point, 0)
	for scanner.Scan() {
		rockPath := parseRockPath(scanner.Text())
		rockPaths = append(rockPaths, rockPath)
	}
	return rockPaths
}

func findMaxY(rockPaths [][]Point) int {
	currentMax := 0
	for _, path := range rockPaths {
		for _, point := range path {
			if point.y > currentMax {
				currentMax = point.y
			}
		}
	}
	return currentMax
}

func findMinX(rockPaths [][]Point) int {
	currentMin := rockPaths[0][0].x
	for _, path := range rockPaths {
		for _, point := range path {
			if point.x < currentMin {
				currentMin = point.x
			}
		}
	}
	return currentMin
}

func findMaxX(rockPaths [][]Point) int {
	currentMax := 0
	for _, path := range rockPaths {
		for _, point := range path {
			if point.x > currentMax {
				currentMax = point.x
			}
		}
	}
	return currentMax
}

func createRockMap(rockPaths [][]Point) ([][]string, Point) {
	maxY := findMaxY(rockPaths)
	maxX := findMaxX(rockPaths)
	minX := findMinX(rockPaths)
	xStart := (maxY*2 - (maxX - minX)) / 2
	rockMap := make([][]string, maxY+2)
	for i := range rockMap {
		// max width is twice the height for this triangle
		// (plus 11 and I'm too tired to work out why)
		rockMap[i] = make([]string, maxY*2+11)
		for j := range rockMap[i] {
			rockMap[i][j] = "."
		}
	}
	for _, path := range rockPaths {
		for _, point := range path {
			rockMap[point.y][point.x-minX+xStart] = "#"
		}
	}
	start := Point{x: 500 - minX + xStart, y: 0}
	rockMap[start.y][start.x] = "+"
	return rockMap, start
}

func displayRockMap(rockMap [][]string) {
	for _, row := range rockMap {
		fmt.Println(row)
	}
}

func grainDidNotSpill(rockMap [][]string, start Point) ([][]string, bool) {
	grain := start
	for {
		if grain.y == len(rockMap)-1 {
			return rockMap, false
		} else if rockMap[grain.y+1][grain.x] == "." {
			grain.y++
		} else if rockMap[grain.y+1][grain.x-1] == "." {
			grain.x--
			grain.y++
		} else if rockMap[grain.y+1][grain.x+1] == "." {
			grain.x++
			grain.y++
		} else {
			rockMap[grain.y][grain.x] = "o"
			break
		}
	}
	return rockMap, true
}

func findGrainsUntilSpill(rockMap [][]string, start Point) int {
	grains := 0
	grain := start
	noGrainSpill := true
	rockMap, noGrainSpill = grainDidNotSpill(rockMap, grain)
	for noGrainSpill {
		grains++
		grain = start
		rockMap, noGrainSpill = grainDidNotSpill(rockMap, grain)
	}
	return grains
}

func sourceNotBlocked(rockMap [][]string, start Point) ([][]string, bool) {
	grain := start
	for {
		if grain.y == len(rockMap)-1 {
			rockMap[grain.y][grain.x] = "o"
			return rockMap, true
		} else if rockMap[grain.y+1][grain.x] == "." {
			grain.y++
		} else if rockMap[grain.y+1][grain.x-1] == "." {
			grain.x--
			grain.y++
		} else if rockMap[grain.y+1][grain.x+1] == "." {
			grain.x++
			grain.y++
		} else {
			if grain == start {
				return rockMap, false
			}
			rockMap[grain.y][grain.x] = "o"
			break
		}
	}
	return rockMap, true
}

func findGrainsUntilBlockedSource(rockMap [][]string, start Point) int {
	grains := 0
	grain := start
	sourceUnblocked := true
	rockMap, sourceUnblocked = sourceNotBlocked(rockMap, grain)
	for sourceUnblocked {
		grains++
		grain = start
		rockMap, sourceUnblocked = sourceNotBlocked(rockMap, grain)
	}
	grains++
	return grains
}

func copyRocksMap(rockMap [][]string) [][]string {
	withFloor := make([][]string, len(rockMap))
	for i := range rockMap {
		withFloor[i] = make([]string, len(rockMap[i]))
		copy(withFloor[i], rockMap[i])
	}
	return withFloor
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
	rockPaths := parseRockPaths(scanner)
	rocksMap, source := createRockMap(rockPaths)
	rocksMapWithFloor := copyRocksMap(rocksMap)
	numOfGrainsUntilFall := findGrainsUntilSpill(rocksMap, source)
	numOfGrainsUntilBlock := findGrainsUntilBlockedSource(rocksMapWithFloor, source)

	elapsed := time.Since(start)
	fmt.Println("Grains of sand until fall:", numOfGrainsUntilFall)
	fmt.Println("Grains of sand until source blocked:", numOfGrainsUntilBlock)
	log.Printf("Time taken: %s", elapsed)
}
