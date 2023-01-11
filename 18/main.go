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

type PointSet map[Point]struct{}

type PointQueue []Point

type Point struct {
	x int
	y int
	z int
}

type BoundingCube struct {
	minX int
	maxX int
	minY int
	maxY int
	minZ int
	maxZ int
}

func (queue *PointQueue) enqueue(point Point) {
	*queue = append(*queue, point)
}

func (queue *PointQueue) dequeue() Point {
	point := (*queue)[0]
	*queue = (*queue)[1:]
	return point
}

func parseLava(scanner *bufio.Scanner) (map[Point]int, PointSet) {
	sides := 6
	lavaSurfaceArea := make(map[Point]int)
	lavaPoints := make(PointSet)
	for scanner.Scan() {
		pointSplit := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(pointSplit[0])
		y, _ := strconv.Atoi(pointSplit[1])
		z, _ := strconv.Atoi(pointSplit[2])
		point := Point{x: x, y: y, z: z}
		lavaSurfaceArea[point] = sides
		lavaPoints[point] = struct{}{}
	}
	return lavaSurfaceArea, lavaPoints
}

func getAdjPoints(p Point) []Point {
	return []Point{
		{x: p.x + 1, y: p.y, z: p.z},
		{x: p.x - 1, y: p.y, z: p.z},
		{x: p.x, y: p.y + 1, z: p.z},
		{x: p.x, y: p.y - 1, z: p.z},
		{x: p.x, y: p.y, z: p.z + 1},
		{x: p.x, y: p.y, z: p.z - 1},
	}
}

func calcSurfaceArea(lavaSurfaceArea map[Point]int, lavaPoints PointSet) int {
	for point := range lavaPoints {
		adjPoints := getAdjPoints(point)
		for _, adjPoint := range adjPoints {
			lavaSurfaceArea[adjPoint]--
		}
	}
	surfaceArea := 0
	for point := range lavaPoints {
		surfaceArea += lavaSurfaceArea[point]
	}
	return surfaceArea
}

// this is a bit gross
func findBoundingCube(points PointSet) BoundingCube {
	boundingCube := BoundingCube{
		minX: math.MaxInt,
		maxX: math.MinInt,
		minY: math.MaxInt,
		maxY: math.MinInt,
		minZ: math.MaxInt,
		maxZ: math.MinInt,
	}
	for p := range points {
		if p.x > boundingCube.maxX {
			boundingCube.maxX = p.x
		}
		if p.x < boundingCube.minX {
			boundingCube.minX = p.x
		}
		if p.y > boundingCube.maxY {
			boundingCube.maxY = p.y
		}
		if p.y < boundingCube.minY {
			boundingCube.minY = p.y
		}
		if p.z > boundingCube.maxZ {
			boundingCube.maxZ = p.z
		}
		if p.z < boundingCube.minZ {
			boundingCube.minZ = p.z
		}
	}
	// make the box bigger by one in each direction
	// to allow exploration around border
	boundingCube.minX--
	boundingCube.minY--
	boundingCube.minZ--
	boundingCube.maxX++
	boundingCube.maxY++
	boundingCube.maxZ++
	return boundingCube
}

func isInBoundary(point Point, boundingCube BoundingCube) bool {
	return point.x <= boundingCube.maxX &&
		point.x >= boundingCube.minX &&
		point.y <= boundingCube.maxY &&
		point.y >= boundingCube.minY &&
		point.z <= boundingCube.maxZ &&
		point.z >= boundingCube.minZ
}

func createExternalSurfaceAreaMap(lavaPoints PointSet) map[Point]int {
	boundingCube := findBoundingCube(lavaPoints)
	explored := make(PointSet)
	lavaBoundaryPoints := make(map[Point]int)
	pointQueue := make(PointQueue, 0)
	pointQueue.enqueue(Point{x: boundingCube.minX, y: boundingCube.minY, z: boundingCube.minZ})
	for len(pointQueue) > 0 {
		currentPoint := pointQueue.dequeue()
		adjPoints := getAdjPoints(currentPoint)
		for _, adjPoint := range adjPoints {
			_, isLava := lavaPoints[adjPoint]
			_, isExplored := explored[adjPoint]
			if isLava {
				lavaBoundaryPoints[adjPoint]++
			}
			if isInBoundary(adjPoint, boundingCube) && !isLava && !isExplored {
				explored[adjPoint] = struct{}{}
				pointQueue.enqueue(adjPoint)
			}
		}
	}
	return lavaBoundaryPoints
}

func calcExternalSurfaceArea(lavaPoints PointSet) int {
	externalLavaSurfaceAreaMap := createExternalSurfaceAreaMap(lavaPoints)
	externalSurfaceArea := 0
	for _, surfaceArea := range externalLavaSurfaceAreaMap {
		externalSurfaceArea += surfaceArea
	}
	return externalSurfaceArea
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
	lavaSurfaceArea, lavaPoints := parseLava(scanner)
	surfaceArea := calcSurfaceArea(lavaSurfaceArea, lavaPoints)
	externalSurfaceArea := calcExternalSurfaceArea(lavaPoints)

	elapsed := time.Since(start)
	fmt.Println("Surface area:", surfaceArea)
	fmt.Println("External surface area:", externalSurfaceArea)
	log.Printf("Time taken: %s", elapsed)
}
