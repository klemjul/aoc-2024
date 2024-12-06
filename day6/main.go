package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

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
	dictinctPositionsVisitedByGuard := countDistinctPositionsGuardVisited(guardMap)

	fmt.Printf("Guard visit %v disting positions\n", dictinctPositionsVisitedByGuard)
}

func part2() {
}

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

func countDistinctPositionsGuardVisited(guardMap [][]string) int {
	distinctPos := 0
	for i := 0; i < len(guardMap); i++ {
		for j := 0; j < len(guardMap[i]); j++ {
			if guardMap[i][j] == "X" {
				distinctPos += 1
			}
		}
	}
	return distinctPos
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
