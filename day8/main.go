package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Position struct {
	Row int
	Col int
}

func main() {
	gameMap, err := parseMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing map:", err)
		return
	}
	part1(copyMap((gameMap)))
	part2(copyMap((gameMap)))
}

func part1(gameMap [][]string) {
	antinodePositionsWithoutDuplicates := countUniqueAntinodePositions(gameMap, false)
	fmt.Printf("Part1: Number of antinode positions: %d\n", antinodePositionsWithoutDuplicates)
}

func part2(gameMap [][]string) {
	antinodePositionsWithoutDuplicates := countUniqueAntinodePositions(gameMap, true)
	fmt.Printf("Part2: Number of antinode positions: %d\n", antinodePositionsWithoutDuplicates)
}

func countUniqueAntinodePositions(gameMap [][]string, inLineAntinodes bool) int {
	antennas := findAntennas(gameMap)
	mapLimits := calcMapLimits(gameMap)
	allAntinodePositions := []Position{}

	for _, antennaPositions := range antennas {
		antinodePositions := calcAntennaAntinodes(mapLimits, antennaPositions, inLineAntinodes)
		allAntinodePositions = append(allAntinodePositions, antinodePositions...)
		if inLineAntinodes {
			allAntinodePositions = append(allAntinodePositions, antennaPositions...)
		}
	}

	antinodePositionsWithoutDuplicates := []Position{}
	for _, antinodePosition := range allAntinodePositions {
		if !slices.Contains(antinodePositionsWithoutDuplicates, antinodePosition) {
			antinodePositionsWithoutDuplicates = append(antinodePositionsWithoutDuplicates, antinodePosition)
		}
	}
	return len(antinodePositionsWithoutDuplicates)
}

func calcAntennaAntinodes(mapLimits Position, antennaPositions []Position, inLineAntinodes bool) []Position {
	antinodePositions := []Position{}
	for _, firstAntennaPosition := range antennaPositions {
		for _, secondAntennaPosition := range antennaPositions {
			if firstAntennaPosition == secondAntennaPosition {
				continue
			}
			antinodePos := Position{
				Row: 2*firstAntennaPosition.Row - secondAntennaPosition.Row,
				Col: 2*firstAntennaPosition.Col - secondAntennaPosition.Col,
			}
			if antinodePos.Row < 0 || antinodePos.Row >= mapLimits.Row || antinodePos.Col < 0 || antinodePos.Col >= mapLimits.Col {
				continue
			}
			antinodePositions = append(antinodePositions, antinodePos)

			antinodePosInMap := true
			if inLineAntinodes {
				for antinodePosInMap {
					antinodePos = Position{
						Row: antinodePos.Row + firstAntennaPosition.Row - secondAntennaPosition.Row,
						Col: antinodePos.Col + firstAntennaPosition.Col - secondAntennaPosition.Col,
					}
					if antinodePos.Row < 0 || antinodePos.Row >= mapLimits.Row || antinodePos.Col < 0 || antinodePos.Col >= mapLimits.Col {
						antinodePosInMap = false
					} else {
						antinodePositions = append(antinodePositions, antinodePos)
					}
				}
			}
		}
	}
	return antinodePositions
}

func calcMapLimits(gameMap [][]string) Position {
	return Position{
		Row: len(gameMap),
		Col: len(gameMap[0]),
	}
}

func findAntennas(gameMap [][]string) map[string][]Position {
	foundAntennas := map[string][]Position{}
	for r := range gameMap {
		for c := range gameMap[r] {
			if gameMap[r][c] != "." {
				foundAntennas[gameMap[r][c]] = append(foundAntennas[gameMap[r][c]], Position{Row: r, Col: c})
			}
		}
	}
	return foundAntennas
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
