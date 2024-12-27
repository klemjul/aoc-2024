// Thanks to https://github.com/derailed-dash/Advent-of-Code/blob/master/src/AoC_2024/Dazbo's_Advent_of_Code_2024.ipynb

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	gardenMap, err := parseGardenMap("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing garden map:", err)
		os.Exit(1)
	}
	part1(gardenMap)
	part2(gardenMap)
}

func part1(gardenMap [][]string) {
	seen := make(map[Position]bool)
	fencesPrice := 0
	for row := range gardenMap {
		for col := range gardenMap[row] {
			if _, ok := seen[Position{Row: row, Col: col}]; ok {
				continue
			}
			regionPlots, regionPerimeter, _ := findRegionFromOrigin(Position{Row: row, Col: col}, gardenMap)
			for _, pos := range regionPlots {
				seen[pos] = true
			}
			fencesPrice += len(regionPlots) * regionPerimeter
		}
	}
	fmt.Println("Part 1:", fencesPrice)
}
func part2(gardenMap [][]string) {
	seen := make(map[Position]bool)
	fencesPrice := 0
	for row := range gardenMap {
		for col := range gardenMap[row] {
			if _, ok := seen[Position{Row: row, Col: col}]; ok {
				continue
			}
			regionPlots, _, regionEdges := findRegionFromOrigin(Position{Row: row, Col: col}, gardenMap)
			for _, pos := range regionPlots {
				seen[pos] = true
			}
			sides := calculcateSides(regionEdges)
			fencesPrice += len(regionPlots) * sides
		}
	}
	fmt.Println("Part 2:", fencesPrice)
}

// using BFS flood fill (BFS is a general graph traversal algorithm and BFS flood fill is a variant that works for regions)
func findRegionFromOrigin(origin Position, garden [][]string) (group []Position, perimeter int, perimeterEdges map[Position][]Position) {
	plant := atPosition(garden, origin)

	seen := make(map[Position]bool)
	seen[origin] = true

	plots := make(map[Position]bool)
	perimeter = 0

	queue := []Position{origin}
	current := Position{}

	perimeterEdges = make(map[Position][]Position)
	for _, dirn := range DirectionCoords {
		perimeterEdges[dirn] = []Position{}
	}

	for len(queue) > 0 {
		current, queue = pop(queue)
		currentPlant := atPosition(garden, current)

		if currentPlant == plant { // This plot is in our region
			plots[current] = true
		}

		for _, dirn := range DirectionCoords { // Get the neighbours, one direction at a time
			neighbour := Position{Row: current.Row + dirn.Row, Col: current.Col + dirn.Col}

			// if the neighbour in grid and same plant, it's in the same region so queue it
			if inGrid(garden, neighbour) && atPosition(garden, neighbour) == plant {
				if !seen[neighbour] {
					queue = append(queue, neighbour)
					seen[neighbour] = true
				}
			} else {
				/* this neighbour represents a perimeter
				 ++++
				+AAAA+
				 ++++
				*/
				perimeter++
				perimeterEdges[dirn] = append(perimeterEdges[dirn], neighbour)
			}
		}
	}

	for plot := range plots {
		group = append(group, plot)
	}

	return group, perimeter, perimeterEdges
}

func calculcateSides(edges map[Position][]Position) int {
	/* visual representation of side identification using perimeter edges
	NORTH        EAST         SOUTH        WEST
	...........  ...........  ...........  ...........
	.RRRR......  .***R......  .****......  .R***......
	.****.RRR..  .***R.**R..  .RR**.**R..  .R***.R**..
	...**R**...  ...****R...  ...****R...  ...R****...
	...****....  ...***R....  ...*RRR....  ...R***....
	...*.......  ...R.......  ...R.......  ...R.......
	...........  ...........  ...........  ...........
	*/
	sides := 0
	for _, edges := range edges {
		seen := make(map[Position]bool)
		for _, edge := range edges { // all edges facing the same direction
			if !seen[edge] {
				sides += 1

				queue := []Position{edge}
				current := Position{}
				for len(queue) > 0 {
					current, queue = pop(queue)
					if seen[current] {
						continue
					}

					seen[current] = true
					for _, dirn := range DirectionCoords {
						neighbour := Position{Row: current.Row + dirn.Row, Col: current.Col + dirn.Col}
						if contains(edges, neighbour) {
							queue = append(queue, neighbour) // adjacent, part of the same side
						}
					}

				}
			}
		}
	}
	return sides
}
func atPosition(garden [][]string, pos Position) string {
	return garden[pos.Row][pos.Col]
}

func inGrid(grid [][]string, pos Position) bool {
	return pos.Row >= 0 && pos.Row < len(grid) && pos.Col >= 0 && pos.Col < len(grid[0])
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func pop[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T
		return zero, s
	}
	return s[len(s)-1], s[:len(s)-1]
}

func parseGardenMap(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := [][]string{}
	for scanner.Scan() {
		lines = append(lines, strings.Split(scanner.Text(), ""))
	}
	return lines, nil
}

type Position struct {
	Row int
	Col int
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
