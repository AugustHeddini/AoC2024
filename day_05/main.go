package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "slices"
)

func sum(toSum []int) (total int) {
    total = 0
    for _, val := range toSum {
        total += val 
    }
    return 
}

func contains(container []int, val int) bool {
    for _, elem := range container {
        if elem == val {
            return true
        }
    }
    return false
}

func isValidUpdate(rules map[int][]int, update []int) bool {
    for i := 1; i < len(update); i++ {
        for _, prev := range update[:i] {
            if contains(rules[update[i]], prev) {
                return false
            }
        }
    }
    return true
}

func getSumOfMiddles(lists [][]int) (total int) {
    total = 0
    for _, list := range lists {
        total += list[len(list)/2]
    }
    return
}

func correctInvalidUpdates(rules map[int][]int, invalidUpdates [][]int) (corrected [][]int) {

    for _, invalid := range invalidUpdates {
        fmt.Printf("Unsorted list %d\n", invalid)
        slices.SortFunc(invalid, func(a, b int) int {
            if contains(rules[a], b) {
                return -1
            } else if contains(rules[b], a) {
                return 1
            }
            return 0
        })
        fmt.Printf("Sorted list %d\n", invalid)
        corrected = append(corrected, invalid)
    }

    return
}

func splitUpdates(rules map[int][]int, updates [][]int) (valid [][]int, invalid [][]int) {
    for _, update := range updates {
        if isValidUpdate(rules, update) {
            valid = append(valid, update)
        } else {
            invalid = append(invalid, update)
        }
    }
    return
}

func parsePrintQueue(input *os.File) (map[int][]int, [][]int) {
    scanner := bufio.NewScanner(input)

    printingRules := make(map[int][]int)
    updates := [][]int{}
    readingRules := true
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" { 
            readingRules = false
            continue
        }

        if readingRules {
            nums := []int{}
            for _, val := range strings.Split(line, "|") {
                if intVal, err := strconv.Atoi(val); err == nil {
                    nums = append(nums, intVal)
                }
            }
            printingRules[nums[0]] = append(printingRules[nums[0]], nums[1])
        }

        if !readingRules {
            updateSteps := []int{}
            for _, val := range strings.Split(line, ",") {
                if intVal, err := strconv.Atoi(val); err == nil {
                    updateSteps = append(updateSteps, intVal)
                }
            }
            updates = append(updates, updateSteps)
        }
    }

    return printingRules, updates
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    printingRules, updates := parsePrintQueue(input)

    valid, invalid := splitUpdates(printingRules, updates)
    fmt.Printf("Valid %d\n Invalid %d\n", valid, invalid)

    fmt.Printf("Sum of valid middles: %d\n", getSumOfMiddles(valid))

    corrected := correctInvalidUpdates(printingRules, invalid)

    fmt.Printf("Sum of corrected middles: %d\n", getSumOfMiddles(corrected))

    fmt.Printf("Elapsed: %s", time.Since(start))
}
