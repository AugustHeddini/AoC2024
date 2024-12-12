package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func length(a int) int {
    x, cnt := 10, 1
    for x <= a {
        x *= 10
        cnt++
    }
    return cnt
}

func pow10(a int) int {
    res := 1
    for range a {
        res *= 10
    }
    return res
}

const BLINKS int = 25
func blinkAtStones(stones []int) []int {
    for range BLINKS {
        stones = iterateStones(stones)
    }
    return stones
}

func iterateStones(stones []int) []int {
    newStones := []int{}
    for _, stone := range stones {
        if stone == 0 {
            newStones = append(newStones, 1)
        } else if nrDigits := length(stone); nrDigits % 2 == 0 {
            divider := pow10(nrDigits/2) 
            front := stone / divider
            back := stone % divider

            newStones = append(newStones, front)
            newStones = append(newStones, back)
        } else {
            newStones = append(newStones, stone * 2024)
        }
    }
    return newStones
}

func parseStones(input *os.File) []int {
    reader := bufio.NewReader(input)

    stones := []int{}
    line, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
    }
    for _, stone := range strings.Split(strings.Trim(line, "\n"), " ") {
        if num, err := strconv.Atoi(stone); err == nil {
            stones = append(stones, num)
        }
    }
    return stones
}


func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    stones := parseStones(input)
    stones = blinkAtStones(stones)

    fmt.Printf("Result %d\n", len(stones))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
