package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	TestValue int
	Operands  []int                // 10 16
	Operators []func(int, int) int // add()
}

func main() {
	equations, err := parseEquations("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing equations:", err)
		os.Exit(1)
	}
	part1(equations)
	part2(equations)
}

func part1(equations []Equation) {
	totalValidEqs := 0
	for _, eq := range equations {
		possibilities := generatePossibilities(len(eq.Operands))
		for _, possibility := range possibilities {
			eq.Operators = possibility
			res, _ := testEquation(eq)
			if res {
				totalValidEqs += eq.TestValue
				break
			}
		}
	}

	fmt.Printf("Number of valid equations: %d\n", totalValidEqs)
}

func part2(equations []Equation) {

}

func concat(a int, b int) int {
	res, _ := parseNumber(strconv.Itoa(a) + strconv.Itoa(b))
	return res
}

func addition(a int, b int) int {
	return a + b
}

func multiplication(a int, b int) int {
	return a * b
}

func generatePossibilities(operandSize int) [][]func(int, int) int {
	possibilities := [][]func(int, int) int{}

	generatePossibilitiesRecursive(operandSize-1, []func(int, int) int{}, &possibilities)

	return possibilities
}

func generatePossibilitiesRecursive(remainingOperations int, currentPossibility []func(int, int) int, possibilities *[][]func(int, int) int) {
	if remainingOperations == 0 {
		*possibilities = append(*possibilities, currentPossibility)
		return
	}

	newPossibility := append([]func(int, int) int{}, currentPossibility...)
	newPossibility = append(newPossibility, addition)
	generatePossibilitiesRecursive(remainingOperations-1, newPossibility, possibilities)

	newPossibility = append([]func(int, int) int{}, currentPossibility...)
	newPossibility = append(newPossibility, multiplication)
	generatePossibilitiesRecursive(remainingOperations-1, newPossibility, possibilities)

	newPossibility = append([]func(int, int) int{}, currentPossibility...)
	newPossibility = append(newPossibility, concat)
	generatePossibilitiesRecursive(remainingOperations-1, newPossibility, possibilities)
}

func testEquation(eq Equation) (bool, error) {
	if len(eq.Operands)-1 != len(eq.Operators) {
		return false, fmt.Errorf("invalid equation: %v", eq)
	}
	if len(eq.Operands) < 2 {
		return false, fmt.Errorf("invalid equation: %v", eq)
	}
	res := eq.Operands[0]
	for iOperand := range eq.Operands {
		if iOperand == len(eq.Operands)-1 {
			break
		}
		res = eq.Operators[iOperand](res, eq.Operands[iOperand+1])
	}
	return eq.TestValue == res, nil
}

func parseEquations(filename string) ([]Equation, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	equations := []Equation{}
	for scanner.Scan() {
		line := scanner.Text()
		rawEq := strings.Split(line, ":")
		if len(rawEq) != 2 {
			return nil, fmt.Errorf("invalid equation: %s", line)
		}
		eqTest, err := parseNumber(rawEq[0])
		if err != nil {
			return nil, err
		}

		eqNumbers := []int{}
		for _, s := range strings.Split(strings.TrimSpace(rawEq[1]), " ") {
			i, err := parseNumber(s)
			if err != nil {
				return nil, err
			}
			eqNumbers = append(eqNumbers, i)
		}
		equations = append(equations, Equation{
			TestValue: eqTest,
			Operands:  eqNumbers,
		})
	}
	return equations, nil
}

func parseNumber(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("could not parse %q as number: %w", s, err)
	}
	return i, nil
}
