package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1(){
	leftNumbers, rightNumbers := parseLists("./inputs.txt")
	if(len(leftNumbers) != len(rightNumbers)) {
		fmt.Printf("Lists are not the same length, left: %d, right: %d\n", len(leftNumbers), len(rightNumbers))
		return
	}
	sort.Ints(leftNumbers)
	sort.Ints(rightNumbers)

	totalDistance := 0
	for i := range leftNumbers {
		totalDistance += int(math.Abs(float64((leftNumbers[i] - rightNumbers[i]))))
	}

	fmt.Printf("Total distance: %d\n", totalDistance)
}

func part2(){
	leftNumbers, rightNumbers := parseLists("./inputs.txt")
	if(len(leftNumbers) != len(rightNumbers)) {
		fmt.Printf("Lists are not the same length, left: %d, right: %d\n", len(leftNumbers), len(rightNumbers))
		return
	}

	similarityScore := 0
	for li := range leftNumbers {
		leftNumberOccurence := 0
		for ri := range rightNumbers {
			if leftNumbers[li] == rightNumbers[ri] {
				leftNumberOccurence++
			}
		}
		similarityScore += leftNumbers[li] * leftNumberOccurence
	}

	fmt.Printf("Similarity score: %d\n", similarityScore)
}



func parseLists(filename string) ([]int, []int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var leftNumbers []int
	var rightNumbers []int

	for scanner.Scan() {
		line := scanner.Text()
		substrings := strings.Split(line, "   ")
		if len(substrings) < 2 {
			fmt.Printf("Failed to parse line: %s\n", line)
			continue;
		}
		leftNumber, err := strconv.Atoi(substrings[0])
		if err != nil {
			fmt.Printf("Failed to parse left number: %s\n", substrings[0])
			continue;
		}
		leftNumbers = append(leftNumbers, leftNumber)

		rightNumber, err := strconv.Atoi(substrings[1])
		if err != nil {
			fmt.Printf("Failed to parse right number: %s\n", substrings[1])
			continue;
		}
		rightNumbers = append(rightNumbers, rightNumber)
	}

	return leftNumbers, rightNumbers
}