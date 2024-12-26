package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"slices"
	"time"
)

var (
    up = coord{ 0, -1 }
    down = coord{ 0, 1 }
    left = coord{ -1, 0 }
    right = coord{ 1, 0 }
)

type coord struct {
    x int
    y int
}

type maze struct {
    obstacles map[coord]bool
    waypoints map[coord]bool
    start coord
    end coord
    bounds coord
}

func add(a *coord, b *coord) coord {
    return coord { a.x + b.x, a.y + b.y }
}

func min(a int, b int) int {
    if a <= b {
        return a
    }
    return b
}

func max(a int, b int) int {
    if a >= b {
        return a
    }
    return b
}

func abs(a int) int {
    if a < 0 {
        return -a
    }
    return a
}

func dist(a coord, b coord) int {
    return abs(a.x - b.x) + abs(a.y - b.y)
}

func isLine(a coord, b coord) bool {
    return a.x == b.x || a.y == b.y
}

func contains(lst []coord, elem coord) bool {
    for _, item := range lst {
        if item == elem {
            return true
        }
    }
    return false
}

func (m *maze) print() {
    for y := range m.bounds.y {
        for x := range m.bounds.x {
            if m.obstacles[coord{ x, y }] {
                fmt.Printf("#")
            } else if m.waypoints[coord{ x, y }] {
                fmt.Printf("O")
            } else {
                fmt.Printf(".")
            }
        }
        fmt.Println()
    }
}

func (m *maze) parseCorners() {
    corners := map[coord]bool{ m.start: true, m.end: true }

    for wall, _ := range m.obstacles {
        nUp := add(&wall, &up)
        nDown := add(&wall, &down)
        // Using neighbours bits to represent neighbour positions
        // where 4 lowest bits are up, down, left, right in that order
        // Such that
        // n = 0b1011
        // represents
        //  .#.
        //  #N#
        //  ...
        neigbours := uint(0)
        if m.obstacles[nUp] {
            neigbours |= 0b1000
        }
        if m.obstacles[nDown] {
            neigbours |= 0b0100
        }
        if m.obstacles[add(&wall, &left)] {
            neigbours |= 0b0010
        }
        if m.obstacles[add(&wall, &right)] {
            neigbours |= 0b0001
        }

        if bits.OnesCount(neigbours) >= 3 || neigbours == 0b1100 || neigbours == 0b0011 {
            continue
        }

        for _, dir := range []coord{ nUp, nDown } {
            if !m.obstacles[dir] {
                if neigbours & 0b0010 == 0 { 
                    diagLeft := add(&dir, &left)
                    if !m.obstacles[diagLeft] {
                        corners[diagLeft] = true
                    }
                }
                if neigbours & 0b0001 == 0 { 
                    diagRight := add(&dir, &right)
                    if !m.obstacles[diagRight] {
                        corners[diagRight] = true
                    }
                }
            }
        }
    }
    m.waypoints = corners
}

func (m *maze) unobstructedLine(from coord, to coord) bool {
    if from.x == to.x {
        for y := min(from.y, to.y); y < max(from.y, to.y); y++ {
            if m.obstacles[coord{ from.x, y }] {
                return false
            }
        }
        return true
    }
    if from.y == to.y {
        for x := min(from.x, to.x); x < max(from.x, to.x); x++ {
            if m.obstacles[coord{ x, from.y }] {
                return false
            }
        }
        return true
    }
    return false
}

func (m *maze) pathfind() []coord {
    distances := map[coord]int{}
    prev := map[coord]coord{m.start: m.start}
    queue := []coord{}
    waypointNeighbours := map[coord][]coord{}
    for wp, _ := range m.waypoints {
        distances[wp] = int(^uint(0)>>1)
        queue = append(queue, wp)
        for n, _ := range m.waypoints {
            if n == wp {
                continue
            }
            if m.unobstructedLine(wp, n) {
                waypointNeighbours[wp] = append(waypointNeighbours[wp], n)
            }
        }
    }
    distances[m.start] = 0

    for len(queue) > 0 {
        slices.SortFunc(queue, func(a, b coord) int {
            if distances[a] < distances[b] {
                return -1
            } else if distances[a] > distances[b] {
                return 1
            }
            return 0
        })
        elem := queue[0]
        if elem == m.end {
            break
        }
        queue = queue[1:]

        for _, n := range waypointNeighbours[elem] {
            newDist := distances[elem] + dist(elem, n)
            if !isLine(prev[elem], n) {
                newDist += 1000
            }
            if newDist < distances[n] {
                distances[n] = newDist 
                prev[n] = elem
            }
        }
    }

    path := []coord{}
    for elem := m.end; elem != m.start; elem = prev[elem] {
        // fmt.Printf("Backtracing elem %d onto path %d\n", elem, path)
        path = append(path, elem)
    }
    path = append(path, m.start)

    fmt.Printf("Found distance to end %d\n", distances[m.end])

    return path
}

func parseInput(input *os.File) maze {
    reader := bufio.NewScanner(input)

    obstacles := map[coord]bool{}
    var start coord
    var end coord

    y := 0
    xMax := 0
    for reader.Scan() {
        for x, b := range reader.Bytes() {
            if y == 0 {
                xMax = x
            }
            if b == '#' {
                obstacles[coord{ x, y }] = true
            } else if b == 'S' {
                start = coord{ x, y }
            } else if b == 'E' {
                end = coord{ x, y }
            }
        }
        y++
    }

    bounds := coord{ xMax + 1, y }
    maze := maze{
        obstacles,
        nil,
        start,
        end,
        bounds,
    }
    return maze
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    maze := parseInput(input)
    maze.parseCorners()
    maze.print()

    fmt.Printf("Found path %d\n", maze.pathfind())

    fmt.Printf("Elapsed %s\n", time.Since(start))
}
