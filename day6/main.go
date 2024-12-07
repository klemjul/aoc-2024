package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Coord struct {
	x int
	y int
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var DirectionCoords = map[Direction]Coord{
	Right: {0, 1},
	Down:  {1, 0},
	Left:  {0, -1},
	Up:    {-1, 0},
}

var GuardSymbols = []string{"^", ">", "v", "<"}
var GuardSymbolsDirections = map[string]Direction{"^": Up, ">": Right, "v": Down, "<": Left}
var GuardDirectionSymbols = map[Direction]string{
	Up:    "^",
	Right: ">",
	Down:  "v",
	Left:  "<",
}

func main() {
	part1()
	part2()
}

func part1() {
	guardMap, err := parseMap("./inputs.txt")
	guardInMap := true
	if err != nil {
		fmt.Println("Error parsing map:", err)
		return
	}

	for guardInMap {
		initGuardPos, initGuardSymbol := findGuard(guardMap)
		nextGuardPos, nextGuardSymbol := nextGuardPosition(guardMap, initGuardPos, initGuardSymbol)
		if nextGuardPos.x < 0 || nextGuardPos.y < 0 {
			guardInMap = false
			guardMap[initGuardPos.x][initGuardPos.y] = "X"
			break
		}
		if initGuardSymbol == nextGuardSymbol {
			guardMap[initGuardPos.x][initGuardPos.y] = "X"
			guardMap[nextGuardPos.x][nextGuardPos.y] = nextGuardSymbol
		} else {
			guardMap[initGuardPos.x][initGuardPos.y] = nextGuardSymbol
		}
	}

	dictinctPositionsVisitedByGuard := 0
	for i := 0; i < len(guardMap); i++ {
		for j := 0; j < len(guardMap[i]); j++ {
			if guardMap[i][j] == "X" {
				dictinctPositionsVisitedByGuard += 1
			}
		}
	}

	fmt.Printf("Guard visit %v disting positions\n", dictinctPositionsVisitedByGuard)
}

func part2() {
	gameMap, err := parseMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing map:", err)
		return
	}

	// run without obstacle to identify ceels that can be blocked
	initGuardPos, _ := findGuard(gameMap)
	_, visitedMapCells := runGame(gameMap)
	cellsThatCanBeBlocked := []Coord{}
	for i := 0; i < len(visitedMapCells); i++ {
		for j := 0; j < len(visitedMapCells[i]); j++ {
			if i == initGuardPos.x && j == initGuardPos.y {
				continue
			}
			if slices.Contains(GuardSymbols, visitedMapCells[i][j]) {
				cellsThatCanBeBlocked = append(cellsThatCanBeBlocked, Coord{x: i, y: j})
			}
		}
	}

	// run multiple game to identify cells that loop the guard
	cellsThatLoopTheGuard := []Coord{}
	for _, cell := range cellsThatCanBeBlocked {
		gameMap[cell.x][cell.y] = "#"
		guardStuckInLoop, _ := runGame(gameMap)
		gameMap[cell.x][cell.y] = "."
		if guardStuckInLoop {
			cellsThatLoopTheGuard = append(cellsThatLoopTheGuard, cell)
		}
	}

	fmt.Printf("Cells that loop the guard: %v\n", len(cellsThatLoopTheGuard))
}

func nextDirection(current Direction, angle int) Direction {
	turns := (angle / 90) % 4

	newDirection := (int(current) + turns) % 4
	return Direction(newDirection)
}

func findGuard(guardMap [][]string) (Coord, string) {
	for i := 0; i < len(guardMap); i++ {
		for j := 0; j < len(guardMap[i]); j++ {
			if slices.Contains(GuardSymbols, guardMap[i][j]) {
				return Coord{x: i, y: j}, guardMap[i][j]
			}
		}
	}
	return Coord{-1, -1}, ""
}

func nextGuardPosition(guardMap [][]string, guardPos Coord, guardSymbol string) (Coord, string) {
	guardDirection := GuardSymbolsDirections[guardSymbol]
	nextPos := Coord{guardPos.x + DirectionCoords[guardDirection].x, guardPos.y + DirectionCoords[guardDirection].y}
	if nextPos.x < 0 || nextPos.x >= len(guardMap) || nextPos.y < 0 || nextPos.y >= len(guardMap[0]) {
		return Coord{x: -1, y: -1}, guardSymbol
	}
	if guardMap[nextPos.x][nextPos.y] == "#" {
		guardDirection = nextDirection(guardDirection, 90)
		return guardPos, GuardDirectionSymbols[guardDirection]
	}
	return nextPos, guardSymbol
}

func runGame(gameMap [][]string) (bool, [][]string) {
	visitedMapCells := copyMap(gameMap)
	guardInMap := true
	guardStuckInLoop := false
	guardPos, guardSymbol := findGuard(gameMap)

	for guardInMap {
		nextGuardPos, nextGuardSymbol := nextGuardPosition(gameMap, guardPos, guardSymbol)
		if nextGuardPos.x < 0 || nextGuardPos.y < 0 {
			guardInMap = false
			break
		}
		if visitedMapCells[nextGuardPos.x][nextGuardPos.y] == nextGuardSymbol {
			guardStuckInLoop = true
			break
		}
		if guardSymbol == nextGuardSymbol {
			guardPos = nextGuardPos
			visitedMapCells[nextGuardPos.x][nextGuardPos.y] = nextGuardSymbol
		} else {
			guardSymbol = nextGuardSymbol
		}
	}
	return guardStuckInLoop, visitedMapCells
}

func copyMap(guardMap [][]string) [][]string {
	newMap := make([][]string, len(guardMap))
	for i := range guardMap {
		newMap[i] = make([]string, len(guardMap[i]))
		copy(newMap[i], guardMap[i])
	}
	return newMap
}

func parseMap(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "")
		lines = append(lines, cols)
	}
	return lines, nil
}
