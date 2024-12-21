package main

import (
    "fmt"
    "time"
    "os"
    "bufio"
)

type pair struct {
    x int
    y int
}

type robot struct {
    pos pair
}

type warehouse struct {
    floor [][]byte
    robot
}

func (a *pair) add(b pair) {
    a.x += b.x
    a.y += b.y
}

func (wh *warehouse) get(loc pair) byte {
    return wh.floor[loc.y][loc.x]
}

func (wh *warehouse) set(loc pair, val byte) {
    wh.floor[loc.y][loc.x] = val
}

func (wh *warehouse) swap(src pair, dst pair) {
    temp := wh.get(dst)
    wh.set(dst, wh.get(src))
    wh.set(src, temp)
}

func (wh *warehouse) getBoxCoords(loc pair) (left pair, right pair) {
    if wh.get(loc) == '[' {
        left = loc
        right = add(loc, pair{ 1, 0 })
    } else {
        left = add(loc, pair{ -1, 0 })
        right = loc
    }
    return 
}

func (wh *warehouse) pushBox(loc pair, dir pair, depth int) (map[int][]func(), bool) {
    leftHalf, rightHalf := wh.getBoxCoords(loc)
    nextLeft := add(leftHalf, dir)
    nextRight := add(rightHalf, dir)

    blocked := false
    callbacks := map[int][]func(){}
    if wh.get(nextLeft) == '#' || wh.get(nextRight) == '#' {
        return nil, true
    }
    if nextLeft != rightHalf && (wh.get(nextLeft) == '[' || wh.get(nextLeft) == ']') {
        leftCallback, leftBlocked := wh.pushBox(nextLeft, dir, depth + 1)
        for k, fn := range leftCallback {
            callbacks[k] = append(callbacks[k], fn...)
        }
        blocked = blocked || leftBlocked    
    }
    if nextRight != leftHalf && wh.get(nextRight) == '[' {
        rightCallback, rightBlocked := wh.pushBox(nextRight, dir, depth + 1)
        for k, fn := range rightCallback {
            callbacks[k] = append(callbacks[k], fn...)
        }
        blocked = blocked || rightBlocked    
    }
    if blocked {
        return nil, true
    }

    doMove := func() {
        wh.set(leftHalf, '.')
        wh.set(rightHalf, '.')
        wh.set(nextLeft, '[')
        wh.set(nextRight, ']')
    }
    callbacks[depth] = append(callbacks[depth], doMove)
    return callbacks, false
}

func (wh *warehouse) step(dir pair) {
    inFrontOfRobot := wh.robot.pos
    inFrontOfRobot.add(dir)
    switch wh.get(inFrontOfRobot) {
    case '#':
        return
    case 'O':
        cur := inFrontOfRobot
        for wh.get(cur) == 'O' {
            cur.add(dir)
        }
        if wh.get(cur) == '#' {
            return
        }
        wh.swap(inFrontOfRobot, cur)
    }
    wh.robot.pos.add(dir)
}

func (wh *warehouse) stepWide(dir pair) {
    inFrontOfRobot := wh.robot.pos
    inFrontOfRobot.add(dir)
    switch wh.get(inFrontOfRobot) {
    case '#':
        return
    case '[':
        fallthrough
    case ']':
        callback, blocked := wh.pushBox(inFrontOfRobot, dir, 0)
        if blocked {
            return
        }
        for i := len(callback) - 1; i >= 0; i-- {
            for _, fn := range callback[i] {
                fn()
            }
        }
    }
    wh.robot.pos.add(dir)
}

func add(a pair, b pair) pair {
    return pair{ a.x + b.x, a.y + b.y }
}

func printMap(wh warehouse) {
    fmt.Printf("Warehouse map:\n")
    for y, row := range wh.floor {
        if wh.robot.pos.y == y {
            row[wh.robot.pos.x] = '@'
        }   
        fmt.Printf("%c\n", row)
    }
}

func parseDir(dirChar byte) pair {
    switch dirChar {
    case '>':
        return pair { 1, 0 }
    case '<':
        return pair { -1, 0 }
    case '^':
        return pair { 0, -1 }
    default:
        return pair { 0, 1 }
    }
}

func locateRobot(line []byte) int {
    for i, b := range line {
        if b == '@' {
            return i
        }
    }
    return -1
}

func iterateRobot(moves []pair, wh warehouse) {
    for _, move := range moves {
        // printMap(wh)
        // wh.step(move)
        wh.stepWide(move)
    }
}

func countGPS(wh warehouse) (score int) {
    for y, row := range wh.floor {
        for x, b := range row {
            if b == 'O' || b == '[' {
                score += 100 * y + x
            }
        }
    }
    return
}

func parseWarehouse(reader *bufio.Reader) (wh warehouse) {
    for {
        if line, err := reader.ReadBytes('\n'); err == nil {
            if len(line) == 1 {
                break
            }
            if rb := locateRobot(line); rb > -1 {
                wh.robot = robot { pair { rb, len(wh.floor) } }
            }
            wh.floor = append(wh.floor, line[:len(line) - 1])  // Drop the trailing '\n'
        }
    }
    wh.set(wh.robot.pos, '.')
    return
}

func parseWarehouseWide(reader *bufio.Reader) (wh warehouse) {
    for {
        if line, err := reader.ReadBytes('\n'); err == nil {
            if len(line) == 1 {
                break
            }
            wideLine := []byte{}
            for _, b := range line {
                switch b {
                case '#':
                    wideLine = append(wideLine, '#')
                    wideLine = append(wideLine, '#')
                case '.':
                    wideLine = append(wideLine, '.')
                    wideLine = append(wideLine, '.')
                case 'O':
                    wideLine = append(wideLine, '[')
                    wideLine = append(wideLine, ']')
                case '@':
                    wh.robot = robot { pair { len(wideLine), len(wh.floor) } }
                    wideLine = append(wideLine, '.')
                    wideLine = append(wideLine, '.')
                }
            }
            wh.floor = append(wh.floor, wideLine)
        }
    }
    return
}

func parseMoves(reader *bufio.Reader) (moves []pair) {
    for {
        b, err := reader.ReadByte()
        if err != nil {
            break
        }
        if b == '\n' {
            continue
        }
        moves = append(moves, parseDir(b))
    }

    return 
}

func parseInput(input *os.File) (warehouse, []pair) {
    reader := bufio.NewReader(input)
    // wh := parseWarehouse(reader)
    wh := parseWarehouseWide(reader)
    moves := parseMoves(reader)
    return wh, moves
}

func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }

    warehouse, moves := parseInput(input)
    iterateRobot(moves, warehouse)

    printMap(warehouse)
    fmt.Printf("Found GPS score %d\n", countGPS(warehouse))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
