package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Position struct {
	Row int
	Col int
}
type TopograficMap struct {
	gameMap            [][]int
	trailHeadPositions []Position
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var DirectionCoords = map[Direction]Position{
	Right: {0, 1},
	Down:  {1, 0},
	Left:  {0, -1},
	Up:    {-1, 0},
}

func main() {
	fmt.Println("day10")
	gameMap, err := parseMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing map:", err)
		os.Exit(1)
	}
	part1(*gameMap)
	part2(*gameMap)
}

func part1(topograficMap TopograficMap) {
	mapScore := 0
	for _, trailHeadPosition := range topograficMap.trailHeadPositions {
		mapScore += calcTrailheadScore(trailHeadPosition, topograficMap.gameMap)
	}

	fmt.Printf("Map score: %d\n", mapScore)
}
func part2(topograficMap TopograficMap) {}

func calcTrailheadScore(start Position, gameMap [][]int) int {
	queue := []Position{start}
	currentPos := Position{}
	finalPositions := []Position{}
	for len(queue) > 0 {
		currentPos, queue = pop(queue)

		for _, direction := range DirectionCoords {
			nextPos := Position{Row: currentPos.Row + direction.Row, Col: currentPos.Col + direction.Col}
			if nextPos.Row < 0 || nextPos.Row >= len(gameMap) || nextPos.Col < 0 || nextPos.Col >= len(gameMap[0]) {
				continue
			}

			if gameMap[nextPos.Row][nextPos.Col] != gameMap[currentPos.Row][currentPos.Col]+1 {
				continue
			}

			if gameMap[nextPos.Row][nextPos.Col] == 9 {
				if slices.Contains(finalPositions, nextPos) {
					continue
				}
				finalPositions = append(finalPositions, nextPos)
				continue
			}

			queue = append(queue, nextPos)
		}
	}

	return len(finalPositions)
}

func parseMap(filename string) (*TopograficMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := [][]int{}
	trailHeadPositions := []Position{}
	for scanner.Scan() {
		line := []int{}
		for _, col := range strings.Split(scanner.Text(), "") {
			colNo, err := parseNumber(col)
			if err != nil {
				return nil, err
			}
			if colNo == 0 {
				trailHeadPositions = append(trailHeadPositions, Position{Row: len(lines), Col: len(line)})
			}
			line = append(line, colNo)
		}
		lines = append(lines, line)
	}

	return &TopograficMap{gameMap: lines, trailHeadPositions: trailHeadPositions}, nil
}

func parseNumber(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("could not parse %q as number: %w", s, err)
	}
	return i, nil
}

func pop[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T
		return zero, s
	}
	return s[len(s)-1], s[:len(s)-1]
}
