package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x int
	y int
}

type Sensor struct {
	pos    Point
	beacon Point
	sRange int
}

type Range struct {
	min int
	max int
}

func abs(n int) int {
	if n < 0 {
		return 0 - n
	}
	return n
}

func manhattanDistance(p1 Point, p2 Point) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

func parsePos(line string, prefix string) Point {
	trimmedLine := strings.TrimPrefix(line, prefix)
	splitCoords := strings.Split(trimmedLine, ", ")
	x, _ := strconv.Atoi(splitCoords[0])
	y, _ := strconv.Atoi(strings.TrimPrefix(splitCoords[1], "y="))
	return Point{x: x, y: y}
}

func parseSensor(line string) Sensor {
	splitLine := strings.Split(line, ":")
	pos := parsePos(splitLine[0], "Sensor at x=")
	beacon := parsePos(splitLine[1], " closest beacon is at x=")
	return Sensor{
		pos:    pos,
		beacon: beacon,
		sRange: manhattanDistance(pos, beacon),
	}
}

func parseSensors(scanner *bufio.Scanner) []Sensor {
	sensors := make([]Sensor, 0)
	for scanner.Scan() {
		sensor := parseSensor(scanner.Text())
		sensors = append(sensors, sensor)
	}
	return sensors
}

func yInSensorRange(targetY int, sensor Sensor) bool {
	return abs(sensor.pos.y-targetY) <= sensor.sRange
}

func findSensorsInRangeOfY(sensors []Sensor, targetY int) []Sensor {
	sensorsInRange := make([]Sensor, 0)
	for _, sensor := range sensors {
		if yInSensorRange(targetY, sensor) {
			sensorsInRange = append(sensorsInRange, sensor)
		}
	}
	return sensorsInRange
}

func calculateSensorWidthAtY(sensor Sensor, distanceFromY int) int {
	return (2*sensor.sRange + 1) - (2 * distanceFromY)
}

func countSensorCoveredPosAtY(sensors []Sensor, targetY int) int {
	maxX := math.MinInt
	minX := math.MaxInt
	for _, sensor := range sensors {
		sensorWidthAtTarget := calculateSensorWidthAtY(sensor, abs(targetY-sensor.pos.y))
		halfWidth := (sensorWidthAtTarget - 1) / 2
		if sensor.pos.x-halfWidth < minX {
			minX = sensor.pos.x - halfWidth
		}
		if sensor.pos.x+halfWidth > maxX {
			maxX = sensor.pos.x + halfWidth
		}
	}
	// add one because counting from 0th x value
	return maxX - minX + 1
}

func beaconsOnTargetY(sensors []Sensor, targetY int) int {
	beaconsOnTargetY := make(map[Point]struct{})
	for _, sensor := range sensors {
		if sensor.beacon.y == targetY {
			beaconsOnTargetY[sensor.beacon] = struct{}{}
		}
	}
	return len(beaconsOnTargetY)
}

func calculateTuningFreq(p Point) int {
	return p.x*4000000 + p.y
}

func rangeBreakPos(ranges []Range) int {
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].min == ranges[j].min {
			return ranges[i].max < ranges[j].max
		}
		return ranges[i].min < ranges[j].min
	})
	maxX := ranges[0].max
	for i := 1; i < len(ranges); i++ {
		if maxX < ranges[i].min {
			return maxX + 1
		} else if maxX < ranges[i].max {
			maxX = ranges[i].max
		}
	}
	return -1
}

func findDistressBeacon(sensors []Sensor, maxY int) Point {
	rangeBreakX := -1
	y := 0
	for ; y <= maxY; y++ {
		xRanges := make([]Range, 0)
		for _, sensor := range sensors {
			sensorWidthAtTarget := calculateSensorWidthAtY(sensor, abs(y-sensor.pos.y))
			if sensorWidthAtTarget > 0 {
				halfWidth := (sensorWidthAtTarget - 1) / 2
				xRange := Range{
					min: sensor.pos.x - halfWidth,
					max: sensor.pos.x + halfWidth,
				}
				xRanges = append(xRanges, xRange)
			}
		}
		rangeBreakX = rangeBreakPos(xRanges)
		if rangeBreakX > -1 {
			break
		}
	}
	return Point{x: rangeBreakX, y: y}
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

	targetY := 2000000
	maxY := 4000000
	scanner := bufio.NewScanner(file)
	sensors := parseSensors(scanner)
	sensorsInRange := findSensorsInRangeOfY(sensors, targetY)
	coveredTargetYPosCount := countSensorCoveredPosAtY(sensorsInRange, targetY)
	beaconsOnTargetY := beaconsOnTargetY(sensors, targetY)
	distressBeacon := findDistressBeacon(sensors, maxY)

	elapsed := time.Since(start)
	fmt.Println("Cannot contain a beacon:", coveredTargetYPosCount-beaconsOnTargetY)
	fmt.Println("Distress beacon tuning frequency:", calculateTuningFreq(distressBeacon))
	log.Printf("Time taken: %s", elapsed)
}
