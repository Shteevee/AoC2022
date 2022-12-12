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

type Tile struct {
	height  int
	parent  *Tile
	visited bool
}

type TileMap struct {
	start Point
	end   Point
	tiles [][]Tile
}

type Queue []Point

func (queue *Queue) enqueue(point Point) {
	*queue = append(*queue, point)
}

func (queue *Queue) dequeue() Point {
	point := (*queue)[0]
	*queue = (*queue)[1:]
	return point
}

func parseMap(scanner *bufio.Scanner) TileMap {
	var start Point
	var end Point
	tiles := make([][]Tile, 0)
	i := 0
	for scanner.Scan() {
		tileRow := make([]Tile, 0)
		for j, c := range scanner.Text() {
			if c == 'S' {
				start = Point{x: j, y: i}
				c = 'a'
			} else if c == 'E' {
				end = Point{x: j, y: i}
				c = 'z'
			}
			tileRow = append(tileRow, Tile{height: int(c) - 97, visited: false})
		}
		tiles = append(tiles, tileRow)
		i++
	}
	return TileMap{start: start, end: end, tiles: tiles}
}

func withinClimbingRange(current Point, candidate Point, tiles [][]Tile) bool {
	return tiles[candidate.y][candidate.x].height-tiles[current.y][current.x].height <= 1
}

func inBounds(candidate Point, tiles [][]Tile) bool {
	return candidate.x >= 0 &&
		candidate.x < len(tiles[0]) &&
		candidate.y >= 0 &&
		candidate.y < len(tiles)
}

func findNextPoints(current Point, tileMap TileMap) []Point {
	candidates := []Point{
		{x: current.x - 1, y: current.y},
		{x: current.x + 1, y: current.y},
		{x: current.x, y: current.y - 1},
		{x: current.x, y: current.y + 1},
	}
	finalists := make([]Point, 0)
	for _, candidate := range candidates {
		if inBounds(candidate, tileMap.tiles) &&
			withinClimbingRange(current, candidate, tileMap.tiles) &&
			!hasBeenVisited(candidate, tileMap) {
			finalists = append(finalists, candidate)
		}
	}
	return finalists
}

func markPointAsVisited(point Point, tileMap *TileMap) {
	(*tileMap).tiles[point.y][point.x].visited = true
}

func hasBeenVisited(point Point, tileMap TileMap) bool {
	return tileMap.tiles[point.y][point.x].visited
}

func markTileParent(current Point, candidate Point, tileMap *TileMap) {
	(*tileMap).tiles[candidate.y][candidate.x].parent = &(*tileMap).tiles[current.y][current.x]
}

func createShortestPath(tileMap *TileMap) {
	queue := Queue{tileMap.start}
	markPointAsVisited(tileMap.start, tileMap)
	for len(queue) > 0 {
		current := queue.dequeue()
		if current == tileMap.end {
			break
		}
		for _, candidate := range findNextPoints(current, *tileMap) {
			markPointAsVisited(candidate, tileMap)
			markTileParent(current, candidate, tileMap)
			queue.enqueue(candidate)
		}
	}
}

func getShortestPath(point Point, tileMap TileMap) []Tile {
	path := make([]Tile, 0)
	current := tileMap.tiles[point.y][point.x]
	for current.parent != nil {
		path = append(path, current)
		current = *current.parent
	}
	return path
}

func cleanMap(tileMap *TileMap) {
	for i := range tileMap.tiles {
		for j := range tileMap.tiles[i] {
			tileMap.tiles[i][j].visited = false
			tileMap.tiles[i][j].parent = nil
		}
	}
}

func withinFloorHikeClimbingRange(current Point, candidate Point, tiles [][]Tile) bool {
	return tiles[current.y][current.x].height-tiles[candidate.y][candidate.x].height <= 1
}

func findFloorHikeNextPoints(current Point, tileMap *TileMap) []Point {
	candidates := []Point{
		{x: current.x - 1, y: current.y},
		{x: current.x + 1, y: current.y},
		{x: current.x, y: current.y - 1},
		{x: current.x, y: current.y + 1},
	}
	finalists := make([]Point, 0)
	for _, candidate := range candidates {
		if inBounds(candidate, tileMap.tiles) &&
			withinFloorHikeClimbingRange(current, candidate, tileMap.tiles) &&
			!hasBeenVisited(candidate, *tileMap) {
			finalists = append(finalists, candidate)
		}
	}
	return finalists
}

func findClosestFloorFromEnd(tileMap *TileMap) Point {
	var floorPoint Point
	queue := Queue{tileMap.end}
	markPointAsVisited(tileMap.end, tileMap)
	for len(queue) > 0 {
		current := queue.dequeue()
		if tileMap.tiles[current.y][current.x].height == 0 {
			floorPoint = current
			break
		}
		for _, candidate := range findFloorHikeNextPoints(current, tileMap) {
			markPointAsVisited(candidate, tileMap)
			markTileParent(current, candidate, tileMap)
			queue.enqueue(candidate)
		}
	}
	return floorPoint
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
	tileMap := parseMap(scanner)
	createShortestPath(&tileMap)
	path := getShortestPath(tileMap.end, tileMap)

	//the map should really be immutable instead
	cleanMap(&tileMap)
	firstFloorPoint := findClosestFloorFromEnd(&tileMap)
	floorPath := getShortestPath(firstFloorPoint, tileMap)

	elapsed := time.Since(start)
	fmt.Println(len(path))
	fmt.Println(len(floorPath))
	log.Printf("Time taken: %s", elapsed)
}
