package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type coord struct {
    x int
    y int
}

func add(a coord, b coord) coord {
    return coord {
        a.x + b.x,
        a.y + b.y,
    }
}

func sub(a coord, b coord) coord {
    return coord {
        a.x - b.x,
        a.y - b.y,
    }
}

func gteq(a coord, val int) bool {
    return a.x >= val &&
            a.y >= val
}

func smallerThan(a coord, b coord) bool {
    return a.x < b.x &&
            a.y < b.y
}

func addUnique(val coord, into map[coord]bool) {
    // fmt.Printf("Adding unique at %d\n", val)
    _, ok := into[val]
    if !ok {
        into[val] = true
    }
}

func findAntinodes(antennas map[byte][]coord, boundary coord) map[coord]bool {

    antinodes := make(map[coord]bool)
    for _, coords := range antennas {
        for i := 0; i < len(coords); i++ {
            for j := i + 1; j < len(coords); j++ {
                addUnique(coords[i], antinodes)
                addUnique(coords[j], antinodes)
                diff := sub(coords[i], coords[j])
                over := add(coords[i], diff)
                for gteq(over, 0) && smallerThan(over, boundary) {
                    addUnique(over, antinodes)
                    over = add(over, diff)
                } 
                under := sub(coords[j], diff)
                for gteq(under, 0) && smallerThan(under, boundary) {
                    addUnique(under, antinodes)
                    under = sub(under, diff)
                }
            }
        }
    }
    return antinodes
}


func parseAntennaMap(input *os.File) (map[byte][]coord, coord) {
    result := make(map[byte][]coord)
    reader := bufio.NewScanner(input)

    var y int = 0
    var xmax int = 0
    for reader.Scan() {
        for x, char := range reader.Bytes() {
            if x > xmax { xmax = x }
            if char == '.' {
                continue
            }
            result[char] = append(result[char], coord { x, y })
        }
        y++
    }

    return result, coord { xmax + 1, y }
}


func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    antennaMap, boundary := parseAntennaMap(input)
    // fmt.Printf("Parsed map: %d and boundary %d\n", antennaMap, boundary)

    antinodes := findAntinodes(antennaMap, boundary)

    // fmt.Printf("Antinodes %d\n", antinodes)

    fmt.Printf("Number of antinodes: %d\n", len(antinodes))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
