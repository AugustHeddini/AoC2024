package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const TIME_STEP = 100
const X_BOUND = 101
const Y_BOUND = 103

type pair struct {
    x int
    y int
}

type robot struct {
    pos pair
    v pair
}

func (p *pair) addMut(b pair) {
    p.x += b.x
    p.y += b.y
}

func (r *robot) step() {
    r.pos.addMut(r.v)
}

func (r *robot) maintainBounds() {
    r.pos.x = absMod(r.pos.x, X_BOUND)
    r.pos.y = absMod(r.pos.y, Y_BOUND)
}

func absMod(a int, bound int) int {
    return (a + bound) % bound
}

func hasRobot(bots []*robot, pos pair) bool {
    for _, bot := range bots {
        if bot.pos == pos {
            return true
        }
    }
    return false
}

func printRobots(bots []*robot) {
    for y := range Y_BOUND {
        for x := range X_BOUND {
            if hasRobot(bots, pair { x, y }) {
                fmt.Print("#")
            } else {
                fmt.Print(".")
            }
        }
        fmt.Print("\n")
    }
}

func countQuadrants(bots []*robot) (q1, q2, q3, q4 int) {
    for _, bot := range bots {
        if bot.pos.x == X_BOUND / 2 ||
            bot.pos.y == Y_BOUND / 2 {
            continue
        }
        x_quad := bot.pos.x < X_BOUND / 2
        y_quad := bot.pos.y < Y_BOUND / 2

        if x_quad && y_quad {
            q1++
        } else if !x_quad && y_quad {
            q2++
        } else if x_quad && !y_quad {
            q3++
        } else {
            q4++
        }
    }
    return
}

func simulateRobots(bots []*robot) {
    for _, bot := range bots {
        for range TIME_STEP {
            bot.step()
            bot.maintainBounds()
        }
    }
}

func searchForTree(bots []*robot) {
    const threshold = 250
    for i := range 100*TIME_STEP {
        for _, bot := range bots {
            bot.step()
            bot.maintainBounds()
        }
        q1, q2, q3, q4 := countQuadrants(bots)
        if q1 > threshold || q2 > threshold || q3 > threshold || q4 > threshold {
            printRobots(bots)
            fmt.Println("Timestep: ", i)
        }
    } 
}

func parseInput(input *os.File) []*robot {
    reader := bufio.NewScanner(input)

    robots := []*robot{}
    for reader.Scan() {
        line := reader.Text()
        if len(line) == 0 {
            continue
        }
        segments := strings.Split(line, " ")

        posStrs := strings.Split(strings.TrimPrefix(segments[0], "p="), ",")
        velStrs := strings.Split(strings.TrimPrefix(segments[1], "v="), ",")

        var pos pair
        var vel pair
        x, xErr := strconv.Atoi(posStrs[0])
        y, yErr := strconv.Atoi(posStrs[1])
        if xErr == nil && yErr == nil {
            pos = pair { x, y }
        }

        x, xErr = strconv.Atoi(velStrs[0])
        y, yErr = strconv.Atoi(velStrs[1])
        if xErr == nil && yErr == nil {
            vel = pair { x, y }
        }

        robots = append(robots, &robot { pos, vel })
    }
    return robots
}

func part1(bots []*robot) {
    simulateRobots(bots)
    q1, q2, q3, q4 := countQuadrants(bots)
    fmt.Printf("Found safety factor %d\n", q1 * q2 * q3 * q4)
}

func part2(bots []*robot) {
    searchForTree(bots)
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    robots := parseInput(input)
    // part1(robots)
    part2(robots)

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
