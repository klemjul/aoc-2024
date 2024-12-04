package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {	
	wordSearchGrid, err := parseWordSearch("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}
	fmt.Printf("Number of occurences of XMAS: %v\n", len(findWord(wordSearchGrid, "XMAS")))
}

func part2() {
	wordSearchGrid, err := parseWordSearch("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	fmt.Printf("Number of occurences of X-MAS: %v\n", len(findXmas(wordSearchGrid)))
}

func parseWordSearch(filename string) ( [][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, strings.Split(line, ""))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
 
type Coord struct {
	x int
	y int
}

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
	DownRight
	DownLeft
	UpRight
	UpLeft
)

var DirectionCoords = map[Direction]Coord{
	Right:    {0, 1},
	Down:     {1, 0},
	Left:     {0, -1},
	Up:       {-1, 0},
	DownRight: {1, 1},
	DownLeft: {1, -1},
	UpRight:  {-1, 1},
	UpLeft:   {-1, -1},
}

func oppositeDirection(dir Coord) Coord {
	return Coord{-dir.x, -dir.y}
}

func findXmas(grid [][]string) []Coord {
	occurrences := make([]Coord, 0)
	// iterate over grid
	for iRow, cols := range grid {
		for iCol := range cols {
			currentLetter := grid[iRow][iCol]
			if currentLetter == "A" {
				masCount := 0
				// check all corners of the X to find two MAS
				for _, direction := range []Coord{
					DirectionCoords[DownRight],
					DirectionCoords[DownLeft],
					DirectionCoords[UpRight],
					DirectionCoords[UpLeft],
				} {
					origin := Coord{iRow, iCol}
					startPos := Coord{origin.x + direction.x * 1,origin.y + direction.y * 1}
					if checkWordInDirection(grid, "MAS", startPos, oppositeDirection(direction)) {
						masCount++
						if masCount == 2 {
							occurrences = append(occurrences, origin)
							break;
						}
					}
				}
			}
		}
	}
	return occurrences
}

func findWord(grid [][]string, word string) []Coord {
	occurrences := make([]Coord, 0)
	// iterate over grid
	for iRow, cols := range grid {
		for iCol := range cols {
			// check all 8 directions to find word
			for _, direction := range DirectionCoords {
				if checkWordInDirection(grid, word, Coord{iRow, iCol}, direction) {
					occurrences = append(occurrences, Coord{iRow, iCol})
				}
			}
		}
	}
	return occurrences
}

func checkWordInDirection(grid [][]string, word string, origin Coord, direction Coord) bool {
	for iLetter := 0; iLetter < len(word); iLetter++ {
		nextX := origin.x + direction.x * iLetter
		nextY := origin.y + direction.y * iLetter
		// ensure next letter is in grid
		if nextX < 0 || nextX >= len(grid) || nextY < 0 || nextY >= len(grid[nextX]) {
			return false
		}
		nextWordLetter := string([]rune(word)[iLetter])
		nextGridLetter := string([]rune(grid[nextX][nextY])[0])
		if nextWordLetter != nextGridLetter {
			return false
		}
	}
	return true
}
