package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main(){
	part1()
	part2()
}

func part1() {
	instructions := parseInstructions("./inputs.txt")

	numbers, _ := extractMul(instructions)

	result := 0
	for _, numberPair := range numbers {
		result += numberPair[0] * numberPair[1]
	}

	fmt.Printf("Result after cleaning memory : %v\n", result) 
}

func part2() {
	instructions := parseInstructions("./inputs.txt")
	pattern := `(mul\(([0-9]+),([0-9]+)\)|do\(\)|don't\(\))`

	r, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex pattern:", err)
		return
	}

	matchIndex := r.FindAllStringSubmatchIndex(instructions, -1)

	shouldMultiply := true
	numbers := [][]int{}
	result := 0
	for _, m := range matchIndex {
		start := m[0]
		end := m[1]
		value := instructions[start:end]
		if value == "do()" {
			shouldMultiply = true
		}else if value == "don't()" {
			shouldMultiply = false
		} else {
			if(shouldMultiply){
				newNumbers, _ := extractMul(value)
				numbers = append(numbers,newNumbers...)
			}
		}
	}
	
	for _, numberPair := range numbers {
		result += numberPair[0] * numberPair[1]
	}

	fmt.Printf("Result after cleaning memory with do() / don't() instructions : %v\n", result)
}

func extractMul(mulString string)  ( [][]int, error) {
	numbers := [][]int{}
	pattern := `mul\(([0-9]+),([0-9]+)\)`

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex pattern: %w", err)
	}

	match := r.FindAllStringSubmatch(mulString, -1)
	for _, match := range match {
		x, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("error converting left number: %w", err)
		}
		y, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, fmt.Errorf("error converting right number: %w", err)
		}
		numbers = append(numbers, []int{x, y})
	}
	return numbers,nil
}

func parseInstructions(filename string) string{
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	return string(data)
}