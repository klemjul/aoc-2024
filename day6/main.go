package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Coord struct {
	Row  int
	Cell int
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
	gameMap, err := parseMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing map:", err)
		return
	}
	part1(gameMap)
	part2(gameMap)
}

func part1(gameMap [][]string) {
	_, visitedMapCells := runGame(gameMap)

	dictinctPositionsVisitedByGuard := 0
	for r := range visitedMapCells {
		for c := range visitedMapCells[r] {
			if slices.Contains(GuardSymbols, visitedMapCells[r][c]) {
				dictinctPositionsVisitedByGuard += 1
			}
		}
	}

	fmt.Printf("Guard visit %v distinct positions\n", dictinctPositionsVisitedByGuard)
}

func part2(gameMap [][]string) {
	// run without obstacle to identify ceels that can be blocked
	initGuardPos, _ := findGuard(gameMap)
	_, visitedMapCells := runGame(gameMap)
	cellsThatCanBeBlocked := []Coord{}
	for r := range visitedMapCells {
		for c := range visitedMapCells[r] {
			if r == initGuardPos.Row && c == initGuardPos.Cell {
				continue
			}
			if slices.Contains(GuardSymbols, visitedMapCells[r][c]) {
				cellsThatCanBeBlocked = append(cellsThatCanBeBlocked, Coord{Row: r, Cell: c})
			}
		}
	}

	// run multiple game to identify cells that loop the guard
	coordsThatLoopTheGuard := []Coord{}
	for _, cell := range cellsThatCanBeBlocked {
		gameMap[cell.Row][cell.Cell] = "#"
		guardStuckInLoop, _ := runGame(gameMap)
		gameMap[cell.Row][cell.Cell] = "."
		if guardStuckInLoop {
			coordsThatLoopTheGuard = append(coordsThatLoopTheGuard, cell)
		}
	}

	fmt.Printf("Coords that loop the guard: %v\n", len(coordsThatLoopTheGuard))
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
				return Coord{Row: i, Cell: j}, guardMap[i][j]
			}
		}
	}
	return Coord{-1, -1}, ""
}

func nextGuardPosition(guardMap [][]string, guardPos Coord, guardSymbol string) (Coord, string) {
	guardDirection := GuardSymbolsDirections[guardSymbol]
	nextPos := Coord{guardPos.Row + DirectionCoords[guardDirection].Row, guardPos.Cell + DirectionCoords[guardDirection].Cell}
	if nextPos.Row < 0 || nextPos.Row >= len(guardMap) || nextPos.Cell < 0 || nextPos.Cell >= len(guardMap[0]) {
		return Coord{Row: -1, Cell: -1}, guardSymbol
	}
	if guardMap[nextPos.Row][nextPos.Cell] == "#" {
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
		if nextGuardPos.Row < 0 || nextGuardPos.Cell < 0 {
			guardInMap = false
			break
		}
		if visitedMapCells[nextGuardPos.Row][nextGuardPos.Cell] == nextGuardSymbol {
			guardStuckInLoop = true
			break
		}
		if guardSymbol == nextGuardSymbol {
			guardPos = nextGuardPos
			visitedMapCells[nextGuardPos.Row][nextGuardPos.Cell] = nextGuardSymbol
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
		return nil, err
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
