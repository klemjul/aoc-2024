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
	count := 0
	cache := make(map[string]int)
	for stone := range initialStones {
		count += countStonesForBlinks(initialStones[stone], blinkCount, cache)
	}
	fmt.Println("Part1: Number of stones after 25 blinks:", count)
}

func part2(initialStones []int) {
	blinkCount := 75
	count := 0
	cache := make(map[string]int)
	for stone := range initialStones {
		count += countStonesForBlinks(initialStones[stone], blinkCount, cache)
	}
	fmt.Println("Part2: Number of stones after 75 blinks:", count)
}

func blinkStone(stone int) []int {
	if stone == 0 {
		return []int{1}
	}

	stoneString := strconv.Itoa(stone)
	if len(stoneString)%2 == 0 {
		middleIndex := len(stoneString) / 2
		leftNumber, _ := parseNumber(stoneString[:middleIndex])
		rightNumber, _ := parseNumber(stoneString[middleIndex:])

		return []int{leftNumber, rightNumber}
	}

	return []int{stone * 2024}
}

func countStonesForBlinks(stone int, blinkCount int, cache map[string]int) int {
	cacheKey := fmt.Sprintf("%d_%d", blinkCount, stone)
	if val, ok := cache[cacheKey]; ok {
		return val
	}
	if blinkCount == 0 {
		return 1
	}
	stoneCount := 0
	newStones := blinkStone(stone)
	for _, newStone := range newStones {
		stoneCount += countStonesForBlinks(newStone, blinkCount-1, cache)
	}

	cache[cacheKey] = stoneCount
	return stoneCount
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
