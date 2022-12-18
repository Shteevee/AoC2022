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

type Pair[T, U any] struct {
	first  T
	second U
}

type Valve struct {
	name           string
	open           bool
	flowRate       int
	tunnels        []*Valve
	valveDistances map[*Valve]int
	visited        bool
}

type Queue []*Valve

func (queue *Queue) enqueue(valve *Valve) {
	*queue = append(*queue, valve)
}

func (queue *Queue) dequeue() *Valve {
	valve := (*queue)[0]
	*queue = (*queue)[1:]
	return valve
}

func parseValve(line string) Pair[*Valve, []string] {
	trimmedLine := strings.TrimPrefix(line, "Valve ")
	name := trimmedLine[:2]
	trimmedLine = strings.TrimPrefix(trimmedLine[2:], " has flow rate=")
	splitLine := strings.Split(trimmedLine, "; ")
	flowRate, _ := strconv.Atoi(splitLine[0])
	tunnelsString := strings.TrimPrefix(splitLine[1], "tunnels lead to valves ")
	tunnelsString = strings.TrimPrefix(tunnelsString, "tunnel leads to valve ")
	return Pair[*Valve, []string]{
		first: &Valve{
			name:     name,
			flowRate: flowRate,
			open:     false,
			visited:  false,
		},
		second: strings.Split(tunnelsString, ", "),
	}
}

func parseValvesInfo(scanner *bufio.Scanner) []Pair[*Valve, []string] {
	valvesInfo := make([]Pair[*Valve, []string], 0)
	for scanner.Scan() {
		valveInfo := parseValve(scanner.Text())
		valvesInfo = append(valvesInfo, valveInfo)
	}
	return valvesInfo
}

func findValveInfo(valvesInfo []Pair[*Valve, []string], name string) *Valve {
	var valve *Valve
	for _, valveInfo := range valvesInfo {
		if valveInfo.first.name == name {
			valve = valveInfo.first
		}
	}
	return valve
}

func parseValves(scanner *bufio.Scanner) []*Valve {
	valvesInfo := parseValvesInfo(scanner)
	valves := make([]*Valve, 0)
	for _, valveInfo := range valvesInfo {
		for _, name := range valveInfo.second {
			valveInfo.first.tunnels = append(
				valveInfo.first.tunnels,
				findValveInfo(valvesInfo, name),
			)
		}
		valves = append(valves, valveInfo.first)
	}
	return valves
}

func findValve(valves []*Valve, name string) *Valve {
	var foundValve *Valve
	for _, valve := range valves {
		if valve.name == name {
			foundValve = valve
		}
	}
	return foundValve
}

func resetVisited(valves []*Valve) {
	for _, valve := range valves {
		valve.visited = false
	}
}

func findDistanceBetween(start *Valve, targetValve *Valve) int {
	prev := make(map[*Valve]*Valve)
	queue := Queue{start}
	start.visited = true
	for len(queue) > 0 {
		current := queue.dequeue()
		if current == targetValve {
			break
		}
		for _, tunnel := range current.tunnels {
			if !tunnel.visited {
				tunnel.visited = true
				prev[tunnel] = current
				queue.enqueue(tunnel)
			}
		}
	}
	distance := 1
	prevValve := prev[targetValve]
	for prevValve != start {
		distance++
		prevValve = prev[prevValve]
	}
	return distance
}

func findShortestPathToWorkingValves(valve *Valve, valves []*Valve) map[*Valve]int {
	valveDistances := make(map[*Valve]int)
	for _, targetValve := range valves {
		if valve != targetValve && targetValve.flowRate > 0 {
			valveDistances[targetValve] = findDistanceBetween(valve, targetValve)
			resetVisited(valves)
		}
	}
	return valveDistances
}

func contructPathDistances(valves []*Valve) {
	for _, valve := range valves {
		// only want to find distances for valves
		// that increase flow (and the start)
		if valve.flowRate > 0 || valve.name == "AA" {
			valveDistances := findShortestPathToWorkingValves(valve, valves)
			valve.valveDistances = valveDistances
		}
	}
}

func copyMap(m map[*Valve]bool) map[*Valve]bool {
	newM := make(map[*Valve]bool)
	for k, v := range m {
		newM[k] = v
	}
	return newM
}

func enoughRemainingTimeToTravel(valves map[*Valve]int, openValves map[*Valve]bool, time int) bool {
	enoughTime := false
	for valve, travelTime := range valves {
		if !openValves[valve] {
			enoughTime = enoughTime || time >= travelTime
		}
	}
	return enoughTime
}

func calcOptimalPressureRelease(
	currentValve *Valve,
	time int,
	flowRate int,
	totalFlow int,
	openValves map[*Valve]bool,
) int {
	if time == 0 {
		return totalFlow
	}
	if len(openValves) == len(currentValve.valveDistances)+1 {
		return totalFlow + flowRate*time
	}
	if !enoughRemainingTimeToTravel(currentValve.valveDistances, openValves, time) {
		return totalFlow + flowRate*time
	}
	localMaxFlowRate := totalFlow
	for nextValve := range currentValve.valveDistances {
		if !openValves[nextValve] {
			nextOpenValves := copyMap(openValves)
			nextOpenValves[nextValve] = true
			timeStep := currentValve.valveDistances[nextValve] + 1
			if time-timeStep >= 0 {
				nextFlowRate := calcOptimalPressureRelease(
					nextValve,
					time-timeStep,
					flowRate+nextValve.flowRate,
					totalFlow+timeStep*flowRate,
					nextOpenValves,
				)
				if nextFlowRate > localMaxFlowRate {
					localMaxFlowRate = nextFlowRate
				}
			}
		}
	}
	return localMaxFlowRate
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
	valves := parseValves(scanner)
	contructPathDistances(valves)
	optimalPressureRelease := calcOptimalPressureRelease(findValve(valves, "AA"), 30, 0, 0, make(map[*Valve]bool))

	elapsed := time.Since(start)
	fmt.Println(valves)
	fmt.Println(optimalPressureRelease)
	log.Printf("Time taken: %s", elapsed)
}
