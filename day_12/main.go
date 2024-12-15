package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

type pos struct {
    x int
    y int
}

type Garden struct {
    plots []plot
}

type plot struct {
    plant byte
    area int
    perimeter int
    perimMap map[pos][]pos
}

func (self pos) checkBounds(bounds pos) bool {
    if self.x < 0 || self.x >= bounds.x {
        return false
    }
    if self.y < 0 || self.y >= bounds.y {
        return false
    }
    return true
}

func (self pos) neighbours() (neighbours []pos) {
    return []pos {
        { self.x - 1, self.y },
        { self.x + 1, self.y },
        { self.x, self.y - 1 },
        { self.x, self.y + 1},
    }
}

func (a pos) add(b pos) pos {
    return pos{ a.x + b.x, a.y + b.y }
}

func (a pos) sub(b pos) pos {
    return pos{ a.x - b.x, a.y - b.y }
}

func (self Garden) calculateFenceCost() (total int) {
    for _, plot := range self.plots {
        total += plot.area * plot.perimeter
    }
    return
}

func (g Garden) calculateBulkFenceCost() (total int) {
    for _, plot := range g.plots {
        sides := plot.findSides()
        total += plot.area * sides 
    }
    return
}

func (self plot) String() string {
    return fmt.Sprintf("{%c: a=%d, p=%d}", self.plant, self.area, self.perimeter)
}

func (p *plot) findPerimeterLine(coord pos, facing pos) (line []pos) {
    dir := pos { 0, 0 }
    if facing.x == 0 {
        dir.x++
    } else {
        dir.y++
    }
    
    stepper := coord.add(dir)
    for {
        if stepAdj, ok := p.perimMap[stepper]; ok {
            expectedBorder := stepper.sub(facing)
            if contains(stepAdj, expectedBorder) {
                line = append(line, stepper)
            } else { break }
            stepper = stepper.add(dir)
        } else {
            break
        }
    }

    stepper = coord
    for {
        if stepAdj, ok := p.perimMap[stepper]; ok {
            expectedBorder := stepper.sub(facing)
            if contains(stepAdj, expectedBorder) {
                line = append(line, stepper)
            } else { break }
            stepper = stepper.sub(dir)
        } else {
            break
        }
    }
    return
}

func (p *plot) findPerimeterNeighbours(coord pos) (neighbours []pos) {
    for _, n := range coord.neighbours() {
        if nAdj, ok := p.perimMap[n]; ok {
            if ok := hasCommonNeighbours(p.perimMap[coord], nAdj); ok {
                neighbours = append(neighbours, n)
            }
        }
    }
    return
}

func (p *plot) findSides() (sides int) {
    visited := map[pos][]pos{}
    for coord, adjacent := range p.perimMap {
        for _, adj := range adjacent {
            facing := coord.sub(adj)

            if contains(visited[coord], facing) {
                continue
            }
            line := p.findPerimeterLine(coord, facing)
            for _, perimPoint := range line {
                visited[perimPoint] = append(visited[perimPoint], facing)
            }
            sides++
        }
    }
    return
}

func (p *plot) sortedPerimKeys() []pos {
    sorted := []pos{}
    for key := range p.perimMap {
        sorted = append(sorted, key)
    }
    slices.SortFunc(sorted, func(a, b pos) int {
        return len(p.perimMap[b]) - len(p.perimMap[a])
    })
    return sorted
}

func hasCommonNeighbours(a []pos, b []pos) bool {
    for _, adj := range a {
        for _, aN := range adj.neighbours() {
            if contains(b, aN) {
                return true
            }
        }
    }
    return false
}

func contains(lst []pos, val pos) bool {
    for _, item := range lst {
        if item == val {
            return true
        }
    }
    return false
}

func parsePlot(coord pos, garden *[][]byte, visited *map[pos]bool) plot {
    bounds := pos{ len(*garden), len((*garden)[0]) }

    plant := (*garden)[coord.x][coord.y]
    area := 0
    perimeter := 0

    perimMap, inspecting := map[pos][]pos{}, []pos{ coord }
    for len(inspecting) > 0 {
        spot := inspecting[0]
        inspecting = inspecting[1:]

        if (*visited)[spot] {
            continue
        }
        (*visited)[spot] = true
        
        bordering := 0
        for _, neighbour := range spot.neighbours() {
            if neighbour.checkBounds(bounds) &&
                (*garden)[neighbour.x][neighbour.y] == plant {
                bordering++
                if !(*visited)[neighbour] {
                    inspecting = append(inspecting, neighbour)
                }
            } else {
                perimMap[neighbour] = append(perimMap[neighbour], spot)
            }
        }
        area++
        perimeter += 4 - bordering
    }
    return plot{ plant, area, perimeter, perimMap }
}

func parseGarden(input *os.File) Garden {
    scanner := bufio.NewScanner(input)

    gardenChars := [][]byte{}
    for scanner.Scan() {
        bytes := scanner.Bytes()
        if len(bytes) == 0 {
            continue
        }
        row := []byte{}
        for _, b := range bytes {
            row = append(row, b)
        }
        gardenChars = append(gardenChars, row) 
    }

    visited := map[pos]bool{}
    garden := Garden { []plot{} }
    for x, line := range gardenChars {
        for y, _ := range line {
            coord := pos { x, y }
            if visited[coord] {
                continue
            }

            plot := parsePlot(coord, &gardenChars, &visited)
            garden.plots = append(garden.plots, plot)
        }
    }
    
    return garden
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }
    garden := parseGarden(input)
    fmt.Printf("Total cost: %d\n", garden.calculateFenceCost())
    fmt.Printf("Total bulk cost: %d\n", garden.calculateBulkFenceCost())

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
