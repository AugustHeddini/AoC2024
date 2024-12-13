package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const BLINKS int = 75
const INTERVAL int = 10

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

func blinkAtStones(stones []int) []int {
    for range BLINKS {
        stones = iterateStones(stones)
    }
    return stones
}

func blinkAtStonesDynamic(stones []int) int {

    for range 5 {
        stones = iterateStones(stones)
    }
    
    var wg sync.WaitGroup
    ch := make(chan int)

    var seen sync.Map
    for _, stone := range stones {
        wg.Add(1)
        go func() {
            defer wg.Done()
            ch <- iterateStoneDynamic(stone, 5, &seen)
        }()
    }
    go func() {
        wg.Wait()
        close(ch)
    }()

    total := 0
    for val := range ch {
        total += val
    }

    return total
}

func iterateStoneDynamic(stone int, blink int, seen *sync.Map) int {
    if blink >= BLINKS {
        return 1
    }
    if dynVals, ok := seen.Load(stone); ok {
        dynStones := dynVals.([]int)
        if blink + INTERVAL >= BLINKS {
            return len(dynStones)
        }
        // fmt.Printf("Loaded stone %d with val %d!\n", stone, dynStones)
        res := 0 
        for _, nextStone := range dynStones{
            res += iterateStoneDynamic(nextStone, blink + INTERVAL, seen)
        }
        return res
    }

    stones := []int{stone}
    for range INTERVAL {
        stones = iterateStones(stones)
    }
    seen.LoadOrStore(stone, stones)
    res := 0
    for _, nextStone := range stones {
        res += iterateStoneDynamic(nextStone, blink + INTERVAL, seen)
    }
    return res
}

func iterateStoneRecursive(stone int, iteration int) int {
    if iteration == BLINKS {
        return 1;
    }
    if stone == 0 {
        return iterateStoneRecursive(1, iteration + 1)
    } else if nrDigits := length(stone); nrDigits % 2 == 0 {
        divider := pow10(nrDigits/2) 
        front := stone / divider
        back := stone % divider

        return iterateStoneRecursive(front, iteration + 1) + iterateStoneRecursive(back, iteration + 1)
    } else {
        return iterateStoneRecursive(stone * 2024, iteration + 1)
    }
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
    // stones = blinkAtStones(stones)
    totalStones := blinkAtStonesDynamic(stones)

    // fmt.Printf("Result %d\n", len(stones))
    fmt.Printf("Result %d\n", totalStones)

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
