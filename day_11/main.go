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

func sumValues(m map[int]int) (total int) {
	for _, val := range m {
		total += val
	}
	return
}

func blinkAtStones(stones []int) (total int) {
	var wg sync.WaitGroup
	ch := make(chan int)
	for _, stone := range stones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- iterateStoneRecursive(stone, 0)
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for val := range ch {
		total += val
	}
	return
}

func blinkAtLanternStones(stones map[int]int) map[int]int {
	for range BLINKS {
		stones = iterateLanternStones(stones)
	}
	return stones
}

func iterateLanternStones(stones map[int]int) map[int]int {
	nextPairs := map[int]int{}
	for stone, nr := range stones {
		nextStones := iterateStone(stone)
		for _, nStone := range nextStones {
			nextPairs[nStone] += nr
		}
	}
	return nextPairs
}

func iterateStone(stone int) []int {
	newStones := []int{}
	if stone == 0 {
		newStones = append(newStones, 1)
	} else if nrDigits := length(stone); nrDigits%2 == 0 {
		divider := pow10(nrDigits / 2)
		front := stone / divider
		back := stone % divider
		newStones = append(newStones, front)
		newStones = append(newStones, back)
	} else {
		newStones = append(newStones, stone*2024)
	}
	return newStones
}

func iterateStoneRecursive(stone int, iteration int) int {
	if iteration == BLINKS {
		return 1
	}
	if stone == 0 {
		return iterateStoneRecursive(1, iteration+1)
	} else if nrDigits := length(stone); nrDigits%2 == 0 {
		divider := pow10(nrDigits / 2)
		front := stone / divider
		back := stone % divider

		return iterateStoneRecursive(front, iteration+1) + iterateStoneRecursive(back, iteration+1)
	} else {
		return iterateStoneRecursive(stone*2024, iteration+1)
	}
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

func parseLanternStones(input *os.File) map[int]int {
	reader := bufio.NewReader(input)

	stones := map[int]int{}
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	for _, stone := range strings.Split(strings.Trim(line, "\n"), " ") {
		if num, err := strconv.Atoi(stone); err == nil {
			stones[num] = 1
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

	// stones := parseStones(input)
	// totalStones := blinkAtStones(stones)

	stoneMap := parseLanternStones(input)
	finalStoneMap := blinkAtLanternStones(stoneMap)
	totalStones := sumValues(finalStoneMap)

	fmt.Printf("Result %d\n", totalStones)

	fmt.Printf("Elapsed: %s\n", time.Since(start))
}
