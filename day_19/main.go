package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func fitsTowel(pattern []byte, idx int, towel []byte) bool {
    if idx + len(towel) >= len(pattern) {
        return false
    }
    for i, b := range towel {
        if pattern[idx + i] != b {
            return false
        }
    }
    return true
}

func testTowels(towels [][]byte, pattern []byte, idx int) bool {
    if idx == len(pattern) - 1 {
        return true
    }
    for _, towel := range towels {
        if fitsTowel(pattern, idx, towel) {
            if testTowels(towels, pattern, idx + len(towel)) {
                return true
            }
        } 
    }
    return false
}

func validateArrangements(towels [][]byte, patterns [][]byte) (valid int) {
    for _, pattern := range patterns {
        if testTowels(towels, pattern, 0) {
            valid++
        }
    }
    return
}

func validateArrangementsParallel(towels [][]byte, patterns [][]byte) (valid int) {
    var wg sync.WaitGroup
    ch := make(chan bool)
    for _, pattern := range patterns {
        wg.Add(1)
        go func() {
            defer wg.Done()
            ch <- testTowels(towels, pattern, 0)
        }()
    }
    go func() {
        wg.Wait()
        close(ch)
    }()

    for possible := range ch {
        if possible {
            valid++
        }
    }
    return
}

func parseInput(input *os.File) (towels [][]byte, patterns [][]byte) {
    reader := bufio.NewReader(input)

    towelLine, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
    }
    towelLine = strings.ReplaceAll(towelLine, ",", "")
    for _, towel := range strings.Split(towelLine, " ") {
        towel = strings.Trim(towel, "\n")
        towels = append(towels, []byte(towel))
    }

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        line = strings.ReplaceAll(line, "\n", "")
        if len(line) == 0 {
            continue
        }
        patterns = append(patterns, []byte(line))
    }

    return
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    towels, patterns := parseInput(input)
    // fmt.Printf("Towels %c\nPatterns %c\n", towels, patterns)
    validPatterns := validateArrangements(towels, patterns)

    fmt.Printf("Found %d valid patterns\n", validPatterns)

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
