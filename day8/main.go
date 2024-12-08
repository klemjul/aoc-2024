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
	part1(gameMap)
	part2(gameMap)
}

func part1(gameMap [][]string) {
	antennas := findAntennas(gameMap)
	mapLimits := calcMapLimits(gameMap)
	allAntinodePositions := []Position{}
	for _, antennaPositions := range antennas {
		antinodePositions := calcAntennaAntinodes(mapLimits, antennaPositions)
		allAntinodePositions = append(allAntinodePositions, antinodePositions...)
	}
	antinodePositionsWithoutDuplicates := slices.DeleteFunc(
		allAntinodePositions,
		func(pos Position) bool {
			found := false
			for _, otherAntinodePosition := range allAntinodePositions {
				if otherAntinodePosition == pos {
					if found {
						return true
					}
					found = true
					continue
				}
			}
			return false
		},
	)
	drawMap(gameMap, antinodePositionsWithoutDuplicates)
	fmt.Printf("Number of antinode positions: %d\n", len(antinodePositionsWithoutDuplicates))

}

func part2(gameMap [][]string) {}

func calcAntennaAntinodes(mapLimits Position, antennaPositions []Position) []Position {
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

func drawMap(gameMap [][]string, antinodePositions []Position) {
	for _, node := range antinodePositions {
		gameMap[node.Row][node.Col] = "#"
	}
	for _, row := range gameMap {
		fmt.Println(strings.Join(row, ""))
	}
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
