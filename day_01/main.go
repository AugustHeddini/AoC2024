package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"sort"
	"strconv"
)

func check (e error) {
	if e != nil {
		panic(e)
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Sum(list []int) int {
	res := 0
	for i := 0; i < len(list); i++ {
		res += list[i]
	}
	return res
}

func Count(x int, list []int) int {
	res := 0
	for i := 0; i < len(list); i++ {
		if list[i] == x {
			res += 1
		}
	}
	return res
}

func parse_input(input *os.File) (left, right []int) {
	line_scanner := bufio.NewScanner(input)

	left_list := []int{}
	right_list := []int{}

	for line_scanner.Scan() {
		line := line_scanner.Text()

		parts := strings.Fields(line)

		if len(parts) != 2 {
			continue
		}

		if i, err := strconv.Atoi(parts[0]); err == nil {
			left_list = append(left_list, i)
		} else {
			panic(err)
		}

		if i, err := strconv.Atoi(parts[1]); err == nil {
			right_list = append(right_list, i)
		} else {
			panic(err)
		}
	}

	return left_list, right_list
}

func find_difference(left, right []int) {
	result := []int{}

	for i := 0; i < len(left); i++ {
		result = append(result, Abs(left[i] - right[i]))
	}

	fmt.Print("Total difference: ")
	fmt.Println(Sum(result))
}

func find_similarity(left, right []int) {

	left_counts := []int{}

	for i :=0; i < len(left); i++ {
		left_counts = append(left_counts, Count(left[i], right))
	}

	left_mult := []int{}
	for i :=0; i < len(left); i++ {
		left_mult = append(left_mult, left[i] * left_counts[i])
	}

	fmt.Print("Total similarity: ")
	fmt.Println(Sum(left_mult))
}

func main() {
	input, err := os.Open("input")
	check(err)

	left, right := parse_input(input)

	sort.Ints(left)
	sort.Ints(right)

	find_difference(left, right)

	find_similarity(left, right)

	defer input.Close()
}