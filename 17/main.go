package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"time"
)

type Piece []image.Point
type RockWindow map[image.Point]struct{}
type State struct {
	rockType int
	move     int
}
type CycleInfo struct {
	rocksFallen int
	height      int
}

const START_HEIGHT_OFFSET = 4
const ROCK_TYPES_NUM = 5
const X_UPPER_BOUND = 6
const X_LOWER_BOUND = 0

func nextPiece(pieceNum int, height int) Piece {
	var piece Piece
	switch pieceNum % ROCK_TYPES_NUM {
	case 0:
		piece = []image.Point{
			{X: 2, Y: height + START_HEIGHT_OFFSET},
			{X: 3, Y: height + START_HEIGHT_OFFSET},
			{X: 4, Y: height + START_HEIGHT_OFFSET},
			{X: 5, Y: height + START_HEIGHT_OFFSET},
		}
	case 1:
		piece = []image.Point{
			{X: 2, Y: height + START_HEIGHT_OFFSET + 1},
			{X: 3, Y: height + START_HEIGHT_OFFSET},
			{X: 3, Y: height + START_HEIGHT_OFFSET + 1},
			{X: 3, Y: height + START_HEIGHT_OFFSET + 2},
			{X: 4, Y: height + START_HEIGHT_OFFSET + 1},
		}
	case 2:
		piece = []image.Point{
			{X: 2, Y: height + START_HEIGHT_OFFSET},
			{X: 3, Y: height + START_HEIGHT_OFFSET},
			{X: 4, Y: height + START_HEIGHT_OFFSET},
			{X: 4, Y: height + START_HEIGHT_OFFSET + 1},
			{X: 4, Y: height + START_HEIGHT_OFFSET + 2},
		}
	case 3:
		piece = []image.Point{
			{X: 2, Y: height + START_HEIGHT_OFFSET},
			{X: 2, Y: height + START_HEIGHT_OFFSET + 1},
			{X: 2, Y: height + START_HEIGHT_OFFSET + 2},
			{X: 2, Y: height + START_HEIGHT_OFFSET + 3},
		}
	case 4:
		piece = []image.Point{
			{X: 2, Y: height + START_HEIGHT_OFFSET},
			{X: 3, Y: height + START_HEIGHT_OFFSET},
			{X: 2, Y: height + START_HEIGHT_OFFSET + 1},
			{X: 3, Y: height + START_HEIGHT_OFFSET + 1},
		}
	}
	return piece
}

func parseMoves(scanner *bufio.Scanner) []rune {
	scanner.Scan()
	return []rune(scanner.Text())
}

func max(xs []int) int {
	currentMax := xs[0]
	for _, value := range xs[1:] {
		if value > currentMax {
			currentMax = value
		}
	}
	return currentMax
}

func createHeights() []int {
	return make([]int, 7)
}

func canMoveLeft(piece Piece, rockWindow RockWindow) bool {
	canMove := true
	for _, pos := range piece {
		_, occupied := rockWindow[image.Point{X: pos.X - 1, Y: pos.Y}]
		canMove = canMove && (pos.X-1 >= X_LOWER_BOUND && !occupied)
	}
	return canMove
}

func canMoveRight(piece Piece, rockWindow RockWindow) bool {
	canMove := true
	for _, pos := range piece {
		_, occupied := rockWindow[image.Point{X: pos.X + 1, Y: pos.Y}]
		canMove = canMove && (pos.X+1 <= X_UPPER_BOUND && !occupied)
	}
	return canMove
}

func canMoveDown(piece Piece, rockWindow RockWindow) bool {
	canMove := true
	for _, pos := range piece {
		_, occupied := rockWindow[image.Point{X: pos.X, Y: pos.Y - 1}]
		canMove = canMove && (!occupied && pos.Y-1 > 0)
	}
	return canMove
}

func movePieceX(piece Piece, rockWindow RockWindow, move rune) Piece {
	if move == '<' && canMoveLeft(piece, rockWindow) {
		for i := range piece {
			piece[i].X--
		}
	}
	if move == '>' && canMoveRight(piece, rockWindow) {
		for i := range piece {
			piece[i].X++
		}
	}
	return piece
}

func adjustHeights(heights []int, piece Piece) []int {
	for _, pos := range piece {
		if pos.Y > heights[pos.X] {
			heights[pos.X] = pos.Y
		}
	}
	return heights
}

func adjustRocks(rockWindow RockWindow, piece Piece) {
	for _, p := range piece {
		rockWindow[p] = struct{}{}
	}
}

func performMoves(numOfRocks int, jets []rune) []int {
	jet := 0
	currentHeights := createHeights()
	rockWindow := make(RockWindow)
	currentPiece := nextPiece(0, 0)
	for i := 0; i < numOfRocks; i++ {
		for {
			currentPiece = movePieceX(currentPiece, rockWindow, jets[jet])
			jet = (jet + 1) % len(jets)
			if canMoveDown(currentPiece, rockWindow) {
				for j := range currentPiece {
					currentPiece[j].Y--
				}
			} else {
				currentHeights = adjustHeights(currentHeights, currentPiece)
				adjustRocks(rockWindow, currentPiece)
				currentPiece = nextPiece(i+1, max(currentHeights))
				break
			}
		}
	}
	return currentHeights
}

func findPatternAndCalcHeight(numOfRocks int, jets []rune) int {
	jet := 0
	currentHeights := createHeights()
	rockWindow := make(RockWindow)
	currentPiece := nextPiece(0, 0)
	cycleCache := make(map[State]CycleInfo)
	for i := 0; i < numOfRocks; i++ {
		// looking for cycles by caching on rock type and jet index
		// (thanks to whoever described this nicely on the subreddit)
		currentState := State{rockType: i % ROCK_TYPES_NUM, move: jet}
		if cycleInfo, exists := cycleCache[currentState]; exists {
			// check that the cycle fits into the remaining number of rocks left
			if remaining, diff := numOfRocks-i, i-cycleInfo.rocksFallen; remaining%diff == 0 {
				return max(currentHeights) + remaining/diff*(max(currentHeights)-cycleInfo.height)
			}
		}
		cycleCache[currentState] = CycleInfo{rocksFallen: i, height: max(currentHeights)}
		for {
			// move the piece with jet
			currentPiece = movePieceX(currentPiece, rockWindow, jets[jet])
			jet = (jet + 1) % len(jets)
			if canMoveDown(currentPiece, rockWindow) {
				// the rock falls
				for j := range currentPiece {
					currentPiece[j].Y--
				}
			} else {
				// adjust the heights and move to next piece
				currentHeights = adjustHeights(currentHeights, currentPiece)
				adjustRocks(rockWindow, currentPiece)
				currentPiece = nextPiece(i+1, max(currentHeights))
				break
			}
		}
	}
	return -1
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
	heights := performMoves(2022, moves)
	calcedHeight := findPatternAndCalcHeight(1000000000000, moves)

	elapsed := time.Since(start)
	fmt.Println(max(heights))
	fmt.Println(calcedHeight)
	log.Printf("Time taken: %s", elapsed)
}
