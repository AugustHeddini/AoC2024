package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const OFFSET int = 10000000000000
const LIMIT int = 100
const A_COST int = 3
const B_COST int = 1

type tuple struct {
    x int
    y int
}

type machine struct {
    a tuple
    b tuple
    prize tuple
}

func (a tuple) mulScalar(multiplier int) tuple {
    return tuple { a.x * multiplier, a.y * multiplier }
}

func (a tuple) div(b tuple) tuple {
    return tuple { a.x / b.x, a.y / b.y }
}

func (a tuple) add(val int) tuple {
    return tuple { a.x + val, a.y + val }
}

func (a tuple) sub(b tuple) tuple {
    return tuple {a.x - b.x, a.y - b.y }
}

func (a tuple) mod(b tuple) tuple {
    return tuple { a.x % b.x, a.y % b.y }
}

func (a tuple) isZero() bool {
    return a.x == 0 && a.y == 0
}

func determinant(left tuple, right tuple) float64 {
    return float64(left.x * right.y - left.y * right.x)
}

func sum(lst []int) (total int) {
    for _, val := range lst {
        total += val
    }
    return
}

func findCostCramer(game machine) (int, error) {
    gameDeterminant := determinant(game.a, game.b)
    if gameDeterminant == 0 {
        return 0, errors.New("Could not find integer solution for coefficients")
    }

    aMul := determinant(game.prize, game.b) / gameDeterminant
    bMul := determinant(game.a, game.prize) / gameDeterminant

    _, adecimals := math.Modf(aMul)
    _, bdecimals := math.Modf(bMul)
    if adecimals != 0 || bdecimals != 0|| 
        aMul < 0 || bMul < 0 {
        return 0, errors.New("Could not find integer solution for coefficients")
    }

    return A_COST * int(math.Round(aMul)) + B_COST * int(math.Round(bMul)), nil
}

func findCost(game machine) (int, error) {
    minCost := math.MaxInt
    for i := range LIMIT {
        removeABatch := game.prize.sub(game.a.mulScalar(i)) 
        if removeABatch.mod(game.b).isZero() {
            coeff := removeABatch.div(game.b)
            if coeff.x != coeff.y {
                continue
            }
            cost := A_COST * i + B_COST * coeff.x
            if cost < minCost {
                minCost = cost
            }
        }
    }
    if minCost == math.MaxInt {
        return 0, errors.New("Could not find valid coefficients")
    }

    return minCost, nil
}

func findMinCost(games []machine) []int {
    costs := []int{}
    for _, game := range games {
        if cost, err := findCostCramer(game); err == nil {
            costs = append(costs, cost)
        }
    }
    return costs
}

func parseInput(input *os.File) []machine {
    reader := bufio.NewScanner(input)

    arcade := []machine{}
    for reader.Scan() {
        line := reader.Text()
        if line == "\n" {
            continue
        }
        var game machine
        for i := range 3 {
            nums := strings.FieldsFunc(line, func(c rune) bool { return !unicode.IsNumber(c) })
            if x, err1 := strconv.Atoi(nums[0]); err1 == nil {
                if y, err2 := strconv.Atoi(nums[1]);  err2 == nil {
                    switch i {
                    case 0:
                        game.a = tuple { x, y }
                    case 1:
                        game.b = tuple { x, y }
                    case 2:
                        // game.prize = tuple { x, y }
                        game.prize = tuple { x, y }.add(OFFSET)
                    }
                }
            }
            reader.Scan()
            line = reader.Text()
        }
        arcade = append(arcade, game)
    }

    return arcade
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    arcade := parseInput(input)
    costs := findMinCost(arcade)
    fmt.Printf("Found total min cost %d for %d prizes\n", sum(costs), len(costs))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
