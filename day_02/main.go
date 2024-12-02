package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"sync"
	"slices"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func parse_reports(input *os.File) ([][]int) {
	lines := bufio.NewScanner(input)

	reports := [][]int{}

	for lines.Scan() {
		report := lines.Text()
		values := []int{}

		if report == "" {
			continue
		}

		for _, val_str := range strings.Fields(report) {
			if val, err := strconv.Atoi(val_str); err == nil {
				values = append(values, val)
			} else {
				panic(err)
			}
		}

		reports = append(reports, values)
	}
	return reports
}

func check_values(a int, b int, increasing bool) bool {
	if (increasing && b < a) || (!increasing && a < b) {
		return false
	}
	return a != b && abs(a-b) < 4
}

func is_safe(report []int, dampener bool) bool {
	increasing := report[0] < report[1]

	for i := 0; i < len(report) - 1; i++ {
		if !check_values(report[i], report[i+1], increasing) {
			if dampener {
				if i != 0 {
					sans_prev := slices.Concat(report[:i-1], report[i:])
					if is_safe(sans_prev, false) {
						return true
					}
				}

				sans_this := slices.Concat(report[:i], report[i+1:])
				if is_safe(sans_this, false) {
					return true
				}


				sans_next := slices.Concat(report[:i+1], report[i+2:])
				if is_safe(sans_next, false) {
					return true
				}
			}
			return false
		}
	}
	return true
}

func check_reports(reports [][]int) int {
	var wg sync.WaitGroup
	ch := make(chan bool)
	for _, report := range reports {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- is_safe(report, true)
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	safe_count := 0
	for safe := range ch {
		if safe {
			safe_count++
		}
	}

	return safe_count
}

func main() {
	input, err := os.Open("input")
	check(err)

	reports := parse_reports(input)

	safe_report_count := check_reports(reports)

	fmt.Print("Safe count: ")
	fmt.Println(safe_report_count)
}
