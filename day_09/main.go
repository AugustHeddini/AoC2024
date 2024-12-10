package main

import (
    "fmt"
    "time"
    "os"
    "bufio"
)

type File struct {
    id int
    len int
}

type Filesystem struct {
    files []File
    spaces []int
    filesystem []int
}

func checksum(fs Filesystem) (total int) {
    for idx, i := range fs.filesystem {
        total += idx * i
    }
    return
}

func format(fs *Filesystem) {
    for len(fs.spaces) != 0 {
        for range fs.files[0].len {
            fs.filesystem = append(fs.filesystem, fs.files[0].id)
        }
        fs.files = fs.files[1:]
        if len(fs.files) == 1 {
            for range fs.files[0].len {
                fs.filesystem = append(fs.filesystem, fs.files[0].id)
            }
            break
        }
        space := fs.spaces[0]
        file := &fs.files[len(fs.files)-1]
        for {
            if space == 0 {
                fs.spaces = fs.spaces[1:]
                break
            }
            fs.filesystem = append(fs.filesystem, file.id)
            file.len--
            if file.len == 0 {
                fs.files = fs.files[:len(fs.files)-1]
                file = &fs.files[len(fs.files)-1]
            }
            space--
        }
    }    
}

func fill(fs *Filesystem) []int {
    spaceIndices := []int{}
    for idx, file := range fs.files {
        for range file.len {
           fs.filesystem = append(fs.filesystem, file.id)
        }
        if idx < len(fs.spaces) {
            spaceIndices = append(spaceIndices, len(fs.filesystem))
            for range fs.spaces[idx] {
                fs.filesystem = append(fs.filesystem, 0)
            }
        }
    }
    return spaceIndices
}

func defragment(fs *Filesystem, spaceIndices []int) {
    for idx := len(fs.files) - 1; idx >= 1; idx-- {
        for sIdx, space := range fs.spaces {
            if spaceIndices[sIdx] > spaceIndices[idx - 1] {
                break
            }
            if fs.files[idx].len <= space {
                for insert := spaceIndices[sIdx]; insert < spaceIndices[sIdx] + fs.files[idx].len; insert++ {
                    fs.filesystem[insert] = fs.files[idx].id
                }
                insert := spaceIndices[idx - 1] + fs.spaces[idx - 1]
                for i := range fs.files[idx].len {
                    fs.filesystem[insert + i] = 0
                }
                fs.spaces[idx - 1] += fs.files[idx].len
                fs.spaces[sIdx] = space - fs.files[idx].len
                spaceIndices[sIdx] += fs.files[idx].len
                break
            }
        }
    }
}

func parseInput(input *os.File) Filesystem {
    reader := bufio.NewReader(input)

    fs := Filesystem {
        []File{},
        []int{},
        []int{},
    }
    id := 0
    readingFile := true
    for {
        b, err := reader.ReadByte()
        if err != nil || b == '\n' {
            break
        }

        i := int(b - '0')
        if readingFile {
            fs.files = append(fs.files, File { id, i })
            id++
        } else {
            fs.spaces = append(fs.spaces, i) 
        }
        readingFile = !readingFile
    }
    
    return fs
}


func main() {
    start := time.Now()

    input, err := os.Open("input")
    if err != nil {
        panic(err)
    }
    defer input.Close()
    filesystem := parseInput(input)
    // format(&filesystem)
    spaceIndices := fill(&filesystem)
    defragment(&filesystem, spaceIndices)

    fmt.Printf("Checksum: %d\n", checksum(filesystem))

    fmt.Printf("Elapsed: %s\n", time.Since(start))
}
