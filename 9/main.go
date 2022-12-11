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

const ROPE_LENGTH = 10

type Move struct {
	direction string
	distance  int
}

type Point struct {
	x int
	y int
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return 0 - n
}

func parseMoves(scanner *bufio.Scanner) []Move {
	moves := make([]Move, 0)
	for scanner.Scan() {
		line := scanner.Text()
		moveSplit := strings.Split(line, " ")
		distance, _ := strconv.Atoi(moveSplit[1])
		moves = append(moves, Move{direction: moveSplit[0], distance: distance})
	}
	return moves
}

func nextTo(head Point, tail Point) bool {
	return abs(head.x-tail.x) <= 1 && abs(head.y-tail.y) <= 1
}

func moveDiagonally(head Point, tail Point) Point {
	if head.x > tail.x && head.y > tail.y {
		tail.x++
		tail.y++
	} else if head.x > tail.x && head.y < tail.y {
		tail.x++
		tail.y--
	} else if head.x < tail.x && head.y > tail.y {
		tail.x--
		tail.y++
	} else if head.x < tail.x && head.y < tail.y {
		tail.x--
		tail.y--
	}
	return tail
}

func moveVertically(head Point, tail Point) Point {
	if head.y > tail.y {
		tail.y++
	} else if head.y < tail.y {
		tail.y--
	}
	return tail
}

func moveHorizontally(head Point, tail Point) Point {
	if head.x > tail.x {
		tail.x++
	} else if head.x < tail.x {
		tail.x--
	}
	return tail
}

func moveKnot(
	head Point,
	tail Point,
	traversedPoints map[Point]struct{},
	isTail bool,
) (Point, map[Point]struct{}) {
	if !nextTo(head, tail) {
		if head.x != tail.x && head.y != tail.y {
			tail = moveDiagonally(head, tail)
		} else if head.x == tail.x && head.y != tail.y {
			tail = moveVertically(head, tail)
		} else if head.x != tail.x && head.y == tail.y {
			tail = moveHorizontally(head, tail)
		}
		if isTail {
			traversedPoints[tail] = struct{}{}
		}
	}
	return tail, traversedPoints
}

func moveRope(rope []Point, traversedPoints map[Point]struct{}) ([]Point, map[Point]struct{}) {
	tailIndex := len(rope) - 1
	for !nextTo(rope[0], rope[1]) {
		for i := 1; i < len(rope); i++ {
			rope[i], traversedPoints = moveKnot(rope[i-1], rope[i], traversedPoints, i == tailIndex)
		}
	}

	return rope, traversedPoints
}

func moveHead(head Point, move Move) Point {
	switch move.direction {
	case "U":
		head.y += move.distance
	case "D":
		head.y -= move.distance
	case "R":
		head.x += move.distance
	case "L":
		head.x -= move.distance
	}
	return head
}

func findTailTraversedPoints(moves []Move) map[Point]struct{} {
	traversedPoints := make(map[Point]struct{})
	rope := make([]Point, ROPE_LENGTH)
	tail := len(rope) - 1
	for i := range rope {
		rope[i] = Point{x: 0, y: 0}
	}
	traversedPoints[rope[tail]] = struct{}{}
	for _, move := range moves {
		rope[0] = moveHead(rope[0], move)
		rope, traversedPoints = moveRope(rope, traversedPoints)
	}

	return traversedPoints
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
	moves := parseMoves(scanner)
	traversedPoints := findTailTraversedPoints(moves)

	elapsed := time.Since(start)
	fmt.Println(len(traversedPoints))
	log.Printf("Time taken: %s", elapsed)
}
