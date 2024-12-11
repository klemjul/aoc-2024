package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	initialStones, err := parseStones("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing stones:", err)
		os.Exit(1)
	}
	part1(initialStones)
	part2(initialStones)
}

func part1(initialStones []int) {
	blinkCount := 25
	stonesStates := [][]int{initialStones}
	for blinkNo := 0; blinkNo < blinkCount; blinkNo++ {
		currentStones := stonesStates[len(stonesStates)-1]
		stonesStates = append(stonesStates, calcNextStoneState(currentStones))
	}
	fmt.Println("Number of stones after 25 blinks:", len(stonesStates[len(stonesStates)-1]))
}

func part2(initialStones []int) {}

func calcNextStoneState(stones []int) []int {
	nextStoneState := []int{}
	for _, stone := range stones {
		if stone == 0 {
			nextStoneState = append(nextStoneState, 1)
			continue
		}

		stoneString := strconv.Itoa(stone)
		if len(stoneString)%2 == 0 {
			middleIndex := len(stoneString) / 2
			leftNumber, _ := parseNumber(stoneString[:middleIndex])
			rightNumber, _ := parseNumber(stoneString[middleIndex:])

			nextStoneState = append(nextStoneState, leftNumber)
			nextStoneState = append(nextStoneState, rightNumber)
			continue
		}

		nextStoneState = append(nextStoneState, stone*2024)
	}
	return nextStoneState
}

func parseStones(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := []int{}
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), " ")
		for col := range cols {
			colNo, err := parseNumber(cols[col])
			if err != nil {
				return nil, err
			}
			line = append(line, colNo)
		}
	}
	return line, nil
}

func parseNumber(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("could not parse %q as number: %w", s, err)
	}
	return i, nil
}
