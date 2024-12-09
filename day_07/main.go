package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

type Op func(int, int) int

func add(a, b int) int {
    return a + b
}

func mul(a, b int) int {
    return a * b
}

func conc(a, b int) int {
    concString := fmt.Sprintf("%d%d", a, b)
    if res, err := strconv.Atoi(concString); err == nil {
        return res
    } else {
        panic(err)
    }
}

func operate(target int, nums []int, acc int, op Op) bool {
    if acc > target {
        return false
    }
    if len(nums) == 0 {
        return acc == target
    }
    return operate(target, nums[1:], op(acc, nums[0]), add) ||
        operate(target, nums[1:], op(acc, nums[0]), mul) ||
        operate(target, nums[1:], op(acc, nums[0]), conc)
}


func validateEquations(equations [][]int) (validSum int) {
    for _, equation := range equations {
        if operate(equation[0], equation[1:], 0, add) {
            validSum += equation[0]
        }
    }
    return
}

func parseInput(inputFile *os.File) [][]int {
    result := [][]int{}
    reader := bufio.NewScanner(inputFile)

    for reader.Scan() {
        line := reader.Text()
        if line == "" {
            continue
        }

        line = strings.ReplaceAll(line, ":", "")

        nums := []int{}
        for _, numString := range strings.Split(line, " ") {
            if num, err := strconv.Atoi(numString); err == nil {
                nums = append(nums, num)
            }
        }

        result = append(result, nums)
    }

    return result
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    parsedInput := parseInput(input)
    validCount := validateEquations(parsedInput)

    fmt.Printf("Sum of valid lines: %d\n", validCount)

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
