package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {
	reports:= parseReports("./inputs.txt")
	unsafeReports := 0
	for _, report := range reports {
		if(len(report) < 2) {
			continue;
		}
		isIncreasing := report[0] < report[len(report) - 1]
		for iLevel, level := range report {
			if iLevel == len(report) - 1 {
				continue;
			}
			levelDif := math.Abs(float64(level - report[iLevel + 1]))
			if levelDif == 0 || levelDif > 3 {
				unsafeReports += 1
				break;
			}
			if isIncreasing && level > report[iLevel + 1] {
				unsafeReports += 1
				break;
			}
			if !isIncreasing && level < report[iLevel + 1] {
				unsafeReports += 1
				break;
			}
		}
	}
	fmt.Printf("Safe reports: %d\n", len(reports) - unsafeReports)
}

func part2() {
	checkReportValidity:= func(report []int) int {
		for i := -1; i < len(report); i++ {
			reportCopy := make([]int, len(report))
			copy(reportCopy, report)
			if i != -1 {
				reportCopy = append(reportCopy[:i], reportCopy[i+1:]...)
			}
			isIncreasing := reportCopy[0] < reportCopy[len(reportCopy) - 1]
			unsafeReport := false
			for iLevel, currentLevel := range reportCopy {
				if iLevel == len(reportCopy) - 1 {
					continue;
				}
				levelDif := math.Abs(float64(currentLevel - reportCopy[iLevel + 1]))
				if levelDif == 0 || levelDif > 3 {
					unsafeReport = true
					break;
				}
				if isIncreasing && currentLevel > reportCopy[iLevel + 1] {
					unsafeReport = true
					break;
				}
				if !isIncreasing && currentLevel < reportCopy[iLevel + 1] {
					unsafeReport = true
					break;
				}
			}
			if !unsafeReport {
				return -1
			}
		}
		return -2
	}

	reports:= parseReports("./inputs.txt")
	unsafeReports := 0
	for _, report := range reports {
		if(len(report) < 2) {
			continue;
		}
		reportValidity := checkReportValidity(report)
		if reportValidity == -2 {
			unsafeReports += 1
		}
	}
	fmt.Printf("Safe reports with tolerance: %d\n", len(reports) - unsafeReports)
}


func parseReports(filename string) ([][]int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", filename)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var reports [][]int

	for scanner.Scan() {
		line := scanner.Text()
		levelsRaw := strings.Split(line, " ")
		levels := make([]int, len(levelsRaw))

		for i, str := range levelsRaw {
			intVal, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error parsing string to int:", err)
				continue;
			}
			levels[i] = intVal
		}
		reports = append(reports, levels)
	}

	return reports
}

