package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const PART_1_MAX_TIME = 24
const PART_2_MAX_TIME = 32

type Pair[T, U any] struct {
	First  T
	Second U
}

type Blueprint struct {
	id              int
	oreBotCost      int
	clayBotCost     int
	obsidianBotCost Pair[int, int]
	geodeBotCost    Pair[int, int]
	maxOreBuy       int
}

type State struct {
	time         int
	ore          int
	clay         int
	obsidian     int
	geode        int
	oreBots      int
	clayBots     int
	obsidianBots int
	geodeBots    int
}

func max(xs []int) int {
	max := 0
	for i := 0; i < len(xs); i++ {
		if xs[i] > max {
			max = xs[i]
		}
	}
	return max
}

func parseBlueprint(line string) Blueprint {
	reg, _ := regexp.Compile(`\d+`)
	matched := reg.FindAllString(line, 7)
	id, _ := strconv.Atoi(matched[0])
	oreBotCost, _ := strconv.Atoi(matched[1])
	clayBotCost, _ := strconv.Atoi(matched[2])
	obsidianBotOreCost, _ := strconv.Atoi(matched[3])
	obsidianBotClayCost, _ := strconv.Atoi(matched[4])
	geodeBotOreCost, _ := strconv.Atoi(matched[5])
	geodeBotClayCost, _ := strconv.Atoi(matched[6])
	return Blueprint{
		id:              id,
		oreBotCost:      oreBotCost,
		clayBotCost:     clayBotCost,
		obsidianBotCost: Pair[int, int]{First: obsidianBotOreCost, Second: obsidianBotClayCost},
		geodeBotCost:    Pair[int, int]{First: geodeBotOreCost, Second: geodeBotClayCost},
		maxOreBuy:       max([]int{oreBotCost, clayBotCost, obsidianBotOreCost, geodeBotOreCost}),
	}
}

func parseBlueprints(scanner *bufio.Scanner) []Blueprint {
	blueprints := make([]Blueprint, 0)
	for scanner.Scan() {
		blueprints = append(blueprints, parseBlueprint(scanner.Text()))
	}
	return blueprints
}

func nextState(state State, timeTaken int) State {
	return State{
		time:         state.time + timeTaken,
		ore:          state.ore + timeTaken*state.oreBots,
		clay:         state.clay + timeTaken*state.clayBots,
		obsidian:     state.obsidian + timeTaken*state.obsidianBots,
		geode:        state.geode + timeTaken*state.geodeBots,
		oreBots:      state.oreBots,
		clayBots:     state.clayBots,
		obsidianBots: state.obsidianBots,
		geodeBots:    state.geodeBots,
	}
}

func shouldBuyGeodeBot(maxTime int, blueprint Blueprint, state State) bool {
	return state.obsidianBots > 0 &&
		state.time < maxTime-1
}

func newGeodeBotState(blueprint Blueprint, state State, timeTaken int) State {
	newState := nextState(state, timeTaken)
	newState.ore -= blueprint.geodeBotCost.First
	newState.obsidian -= blueprint.geodeBotCost.Second
	newState.geodeBots++
	return newState
}

func shouldBuyObsidianBot(maxTime int, blueprint Blueprint, state State) bool {
	return state.obsidianBots < blueprint.geodeBotCost.Second &&
		state.time < maxTime-1 &&
		state.clayBots > 0
}

func newObsidianBotState(blueprint Blueprint, state State, timeTaken int) State {
	newState := nextState(state, timeTaken)
	newState.ore -= blueprint.obsidianBotCost.First
	newState.clay -= blueprint.obsidianBotCost.Second
	newState.obsidianBots++
	return newState
}

func shouldBuyClayBot(maxTime int, blueprint Blueprint, state State) bool {
	return state.clayBots < blueprint.obsidianBotCost.Second &&
		state.time < maxTime-1
}

func newClayBotState(blueprint Blueprint, state State, timeTaken int) State {
	newState := nextState(state, timeTaken)
	newState.ore -= blueprint.clayBotCost
	newState.clayBots++
	return newState
}

func shouldBuyOreBot(maxTime int, blueprint Blueprint, state State) bool {
	return state.oreBots < blueprint.maxOreBuy &&
		state.time < maxTime-1
}

func newOreBotState(blueprint Blueprint, state State, timeTaken int) State {
	newState := nextState(state, timeTaken)
	newState.ore -= blueprint.oreBotCost
	newState.oreBots++
	return newState
}

func calcTimeTaken(cost int, stock int, bots int) int {
	if stock >= cost {
		return 1
	}
	return ((cost - stock + bots - 1) / bots) + 1
}

func findMaxGeodesOfBlueprint(maxTime int, blueprint Blueprint, state State) int {
	if state.time == maxTime {
		return state.geode
	}

	pathResults := make([]int, 0)
	if state.geodeBots > 0 {
		pathResults = append(pathResults, findMaxGeodesOfBlueprint(maxTime, blueprint, nextState(state, maxTime-state.time)))
	}
	if shouldBuyOreBot(maxTime, blueprint, state) {
		timeTaken := calcTimeTaken(blueprint.oreBotCost, state.ore, state.oreBots)
		if state.time+timeTaken <= maxTime {
			pathResults = append(
				pathResults,
				findMaxGeodesOfBlueprint(maxTime, blueprint, newOreBotState(blueprint, state, timeTaken)),
			)
		}
	}
	if shouldBuyClayBot(maxTime, blueprint, state) {
		timeTaken := calcTimeTaken(blueprint.clayBotCost, state.ore, state.oreBots)
		if state.time+timeTaken <= maxTime {
			pathResults = append(
				pathResults,
				findMaxGeodesOfBlueprint(maxTime, blueprint, newClayBotState(blueprint, state, timeTaken)),
			)
		}
	}
	if shouldBuyObsidianBot(maxTime, blueprint, state) {
		timeTaken := max([]int{
			calcTimeTaken(blueprint.obsidianBotCost.First, state.ore, state.oreBots),
			calcTimeTaken(blueprint.obsidianBotCost.Second, state.clay, state.clayBots),
		})
		if state.time+timeTaken <= maxTime {
			pathResults = append(
				pathResults,
				findMaxGeodesOfBlueprint(maxTime, blueprint, newObsidianBotState(blueprint, state, timeTaken)),
			)
		}
	}
	if shouldBuyGeodeBot(maxTime, blueprint, state) {
		timeTaken := max([]int{
			calcTimeTaken(blueprint.geodeBotCost.First, state.ore, state.oreBots),
			calcTimeTaken(blueprint.geodeBotCost.Second, state.obsidian, state.obsidianBots),
		})
		if state.time+timeTaken <= maxTime {
			pathResults = append(
				pathResults,
				findMaxGeodesOfBlueprint(maxTime, blueprint, newGeodeBotState(blueprint, state, timeTaken)),
			)
		}
	}
	return max(pathResults)
}

func calcQualityLevelSum(blueprints []Blueprint) int {
	totalQualityLevel := 0
	for i, blueprint := range blueprints {
		maxGeodes := findMaxGeodesOfBlueprint(PART_1_MAX_TIME, blueprint, State{oreBots: 1})
		totalQualityLevel += (i + 1) * maxGeodes
	}
	return totalQualityLevel
}

func calcMaxGeodeProduct(blueprints []Blueprint) int {
	geodeProduct := 1
	maxGeodes := make(chan int, len(blueprints))
	for _, blueprint := range blueprints {
		go func(bp Blueprint) {
			maxGeodes <- findMaxGeodesOfBlueprint(PART_2_MAX_TIME, bp, State{oreBots: 1})
		}(blueprint)
	}
	maxGeodesFound := 0
	for maxGeodesFound < len(blueprints) {
		maxGeode := <-maxGeodes
		maxGeodesFound++
		geodeProduct = geodeProduct * maxGeode
	}
	return geodeProduct
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
	blueprints := parseBlueprints(scanner)

	qualityLevelSum := calcQualityLevelSum(blueprints)
	fmt.Println("Quality level sum:", qualityLevelSum)
	elapsed := time.Since(start)
	log.Printf("Time taken: %s", elapsed)

	maxGeodeProduct := calcMaxGeodeProduct(blueprints[:3])
	fmt.Println("First 3 blueprints max geode product:", maxGeodeProduct)
	elapsed = time.Since(start)
	log.Printf("Time taken: %s", elapsed)
}
