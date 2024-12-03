package main

import (
    "fmt"
    "bufio"
    "os"
    "regexp"
    "strconv"
    "strings"
)

func Sum(vals []int) int {
    sum := 0
    for _, val := range vals {
        sum += val
    }
    return sum
}

func splitRegexpMatch(data []byte, atEOF bool) (advance int, token []byte, err error) {
    re := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`)

    match := re.FindIndex(data)
    if match != nil {
        return match[1], data[match[0]:match[1]], nil
    }
    if !atEOF {
        return 0, nil, nil
    }
    return 0, nil, bufio.ErrFinalToken
}

func parseInputMuls(input *os.File) []int {
    scanner := bufio.NewScanner(input)
    scanner.Split(splitRegexpMatch)

    do := true
    mulResults := []int{}
    for scanner.Scan() {
        instruction := scanner.Text()
        if instruction == "do()" {
            do = true
            continue
        }
        if instruction == "don't()" {
            do = false
            continue
        }

        if do {
            numStrings := strings.Split(strings.Trim(instruction, "mul()"), ",")
    
            first, _ := strconv.Atoi(numStrings[0])
            second, _ := strconv.Atoi(numStrings[1])
    
            mulResults = append(mulResults, first * second)
        }
    }    
    return mulResults
}

func main() {
    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }
    defer input.Close()

    listOfMuls := parseInputMuls(input)
    fmt.Printf("Sum of multiplications: %d\n", Sum(listOfMuls))
}
