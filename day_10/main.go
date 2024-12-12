package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type pair struct {
    x int
    y int
}

type trail struct {
    pos pair
    bound pair
}

func sum(arr []int) (sum int) {
    for _, val := range arr {
        sum += val
    }
    return
}

func contains(arr *[]pair, val pair) bool {
    for _, arrVal := range *arr {
        if arrVal == val {
            return true
        }
    }
    return false
}

func (self trail) neighbours() (neighbours []trail) {
    if self.pos.x - 1 >= 0 {
        neighbours = append(neighbours, trail { pair { self.pos.x-1, self.pos.y }, self.bound })
    }
    if self.pos.x + 1 < self.bound.x {
        neighbours = append(neighbours, trail { pair { self.pos.x+1, self.pos.y }, self.bound })
    }
    if self.pos.y - 1 >= 0 {
        neighbours = append(neighbours, trail { pair { self.pos.x, self.pos.y-1 }, self.bound })
    }
    if self.pos.y + 1 < self.bound.y {
        neighbours = append(neighbours, trail { pair { self.pos.x, self.pos.y+1 }, self.bound })
    }
    return
}

func evaluatePath(topoMap [][]int, here trail, prev int, visited *[]pair) int {
    height := topoMap[here.pos.y][here.pos.x]
    if height != prev + 1 || contains(visited, here.pos) {
        return 0
    }
    // *visited = append(*visited, here.pos)    // Uncomment for part 1
    if height == 9 {
        return 1
    }
    total := 0
    for _, neighbour := range here.neighbours() {
        total += evaluatePath(topoMap, neighbour, height, visited)
    }
    return total
}

func evaluateTrailhead(topoMap [][]int, x int, y int) (value int) {
    trailhead := pair { x, y }
    bound := pair { len(topoMap), len(topoMap[0]) }
    trail := trail { trailhead, bound }

    visited := []pair{ trailhead }
    for _, neighbour := range trail.neighbours() {
        value += evaluatePath(topoMap, neighbour, 0, &visited)
    }
    return
}

func findTrailheadValues(topoMap [][]int) (trailheads []int) {
    for y, line := range topoMap {
        for x, val := range line {
            if val == 0 {
                trailhead := evaluateTrailhead(topoMap, x, y)
                trailheads = append(trailheads, trailhead)
            }
        }
    }
    return
}

func parseTopography(input *os.File) [][]int {
    reader := bufio.NewReader(input)

    topoMap := [][]int{}
    line := []int{}
    for {
        b, err := reader.ReadByte()
        if err != nil {
            break
        }
        if b == '\n' {
            if len(line) != 0 {
                topoMap = append(topoMap, line)
            }
            line = []int{}
            continue
        }
        
        line = append(line, int(b - '0'))
    }
    return topoMap
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    topoMap := parseTopography(input)

    trailheads := findTrailheadValues(topoMap)

    fmt.Printf("Found trailheads: %d\n", trailheads)
    fmt.Printf("Sum: %d\n", sum(trailheads))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
