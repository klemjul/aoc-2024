package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {
	rules, allPages, err := ParsePrinterInstructions("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing printing rules:", err)
		debug.PrintStack()
		return
	}
	allPagesSorted := GetAllPagesSorted(allPages, rules)
	middleSortedPageNumberCount := 0
	for iPages := range allPagesSorted {
		if !reflect.DeepEqual(allPages[iPages], allPagesSorted[iPages]) {
			continue
		}
		middleIndex := len(allPages[iPages]) / 2
		middlePage := allPages[iPages][middleIndex]
		middleSortedPageNumberCount += middlePage
	}
	fmt.Println("Correctly ordered sorted page number count:", middleSortedPageNumberCount)
}

func part2() {
	rules, allPages, err := ParsePrinterInstructions("./inputs.txt")
	if err != nil {
		fmt.Println("Error parsing printing rules:", err)
		debug.PrintStack()
		return
	}
	allPagesSorted := GetAllPagesSorted(allPages, rules)
	middleSortedPageNumberCount := 0
	for iPages := range allPagesSorted {
		if reflect.DeepEqual(allPages[iPages], allPagesSorted[iPages]) {
			continue
		}
		middleIndex := len(allPagesSorted[iPages]) / 2
		middlePage := allPagesSorted[iPages][middleIndex]
		middleSortedPageNumberCount += middlePage
	}
	fmt.Printf("Incorrectly ordered sorted page number count: %d\n", middleSortedPageNumberCount)
}

type PrintingOrderingRule struct {
	PageBefore int
	PageAfter  int
}

func ParsePrinterInstructions(filename string) ([]PrintingOrderingRule, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	parsePagesToProduce := false

	orderingRules := []PrintingOrderingRule{}
	pagesToProduce := [][]int{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line == "" {
			parsePagesToProduce = true
			continue
		}
		if !parsePagesToProduce {
			orderingRuleRaw := strings.Split(line, "|")
			pageBefore, err := strconv.Atoi(orderingRuleRaw[0])
			if err != nil {
				return nil, nil, err
			}

			pageAfter, err := strconv.Atoi(orderingRuleRaw[1])
			if err != nil {
				return nil, nil, err
			}
			orderingRules = append(orderingRules, PrintingOrderingRule{PageBefore: pageBefore, PageAfter: pageAfter})
		} else {
			pagesUpdateRaw := strings.Split(line, ",")
			pagesUpdate := []int{}
			for _, page := range pagesUpdateRaw {
				pageNumber, err := strconv.Atoi(page)
				if err != nil {
					return nil, nil, err
				}
				pagesUpdate = append(pagesUpdate, pageNumber)
			}
			pagesToProduce = append(pagesToProduce, pagesUpdate)
		}
	}
	return orderingRules, pagesToProduce, nil
}

func GetAllPagesSorted(allPages [][]int, rules []PrintingOrderingRule) [][]int {
	allPagesCopy := make([][]int, len(allPages))
	for i, pages := range allPages {
		pagesCopy := make([]int, len(pages))
		copy(pagesCopy, pages)
		sort.Slice(pagesCopy, func(i, j int) bool {
			for _, rule := range rules {
				if rule.PageBefore == pagesCopy[i] && rule.PageAfter == pagesCopy[j] {
					return true
				}
			}
			return false
		})

		allPagesCopy[i] = pagesCopy
	}
	return allPagesCopy
}
