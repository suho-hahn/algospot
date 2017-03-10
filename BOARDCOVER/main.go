package main

import (
    "log"
    "bufio"
    "os"
    "strconv"
    "fmt"
)

const (
    MAX_WIDTH = 20
    MAX_HEIGHT = 20
)

type Rect struct {
    iStart int
    iLen   int
    jStart int
    jLen   int
}

var cache = make(map[string]int)

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    if ! scanner.Scan() {
        panic(scanner.Err())
    }
    caseNum, _ := strconv.Atoi(scanner.Text())

    //log.Println(caseNum)

    for ; caseNum > 0 ; caseNum -- {
        if ! scanner.Scan() {
            panic(scanner.Err())
        }
        height, _ := strconv.Atoi(scanner.Text())
        if ! scanner.Scan() {
            panic(scanner.Err())
        }
        width, _ := strconv.Atoi(scanner.Text())

        matrix := make([][]int8, height)

        for i:=0; i< height; i++ {
            if ! scanner.Scan() {
                panic(scanner.Err())
            }
            bSlice := scanner.Bytes()
            matrix[i] = make([]int8, width)
            matrix[i][0] = 1
            for j, b := range bSlice {
                if b == '#' {
                    matrix[i][j] = 1
                } else {
                    matrix[i][j] = 0
                }
            }
        }

        //log.Println()
        //log.Println("============================================")
        //log.Println("final result:", solve(copyAndReformMatrix(matrix), 0))
        fmt.Println(solve(copyAndReformMatrix(matrix), 0))


    }
}

var depth = 0
func solve(matrix [][]int8, startCol int) (result int) {

    if result, ok := getCache(matrix); ok {
        return result
    }

    defer func(){
        putCache(matrix, result)
    }()

    depth ++
    defer func(){
        depth --
    }()

    //log.Println("depth", depth)
    //printMatrix(matrix)

    if len(matrix) == 0 {
        return 1
    }

    height := len(matrix)
    width := len(matrix[0])

    if height == 1 {
        for i:=0; i<width; i++ {
            if matrix[height - 1][i] == 0 {
                return 0
            }
        }

        return 1

    }

    resultCount := 0

    i := 0
    for j:=startCol ; j<width ; j++{

        if matrix[i][j] == 1 {
            continue
        }

        fillable := false

        if j > 0 &&
            matrix[i+1][j-1] == 0  &&
            matrix[i][j] == 0 &&
            matrix[i+1][j] == 0 {

            fillable = true

            matrix[i+1][j-1] = 1
            matrix[i][j] = 1
            matrix[i+1][j] = 1

            resultCount += solve(matrix, j)

            matrix[i+1][j-1] = 0
            matrix[i][j] = 0
            matrix[i+1][j] = 0

        }

        if j < width - 1 &&
            matrix[i][j] == 0 &&
            matrix[i][j+1] == 0 &&
            matrix[i+1][j+1] == 0 {

            fillable = true

            matrix[i][j] = 1
            matrix[i][j+1] = 1
            matrix[i+1][j+1] = 1

            resultCount += solve(matrix, j)

            matrix[i][j] = 0
            matrix[i][j+1] = 0
            matrix[i+1][j+1] = 0

        }

        if j < width - 1 &&
            matrix[i][j] == 0 &&
            matrix[i+1][j] == 0  &&
            matrix[i+1][j+1] == 0 {

            fillable = true

            matrix[i][j] = 1
            matrix[i+1][j] = 1
            matrix[i+1][j+1] = 1

            resultCount += solve(matrix, j)

            matrix[i][j] = 0
            matrix[i+1][j] = 0
            matrix[i+1][j+1] = 0

        }

        if j < width - 1 &&
            matrix[i][j] == 0 &&
            matrix[i+1][j] == 0  &&
            matrix[i][j+1] == 0 {

            fillable = true

            matrix[i][j] = 1
            matrix[i+1][j] = 1
            matrix[i][j+1] = 1

            resultCount += solve(matrix, j)

            matrix[i][j] = 0
            matrix[i+1][j] = 0
            matrix[i][j+1] = 0

        }

        if ! fillable {
            return 0
        }

        return resultCount
    }

    return solve(matrix[1:], 0)

}

func copyAndReformMatrix(matrix [][]int8) [][]int8 {

    if len(matrix) == 0 || len(matrix[0]) == 0{
        return [][]int8{}
    }

    copy := make([][]int8, len(matrix))
    for i:=0; i<len(matrix); i++ {
        copy[i] = matrix[i]
    }
    matrix = copy

    height := len(matrix)
    width := len(matrix[0])
    //cut side

    cutTop := 0
    cutBottom := height
    cutLeft := 0
    cutRight := width
    //cut upper side
    for i:=0; i<height; i++ {
        isBreak := false
        for j:=0; j<width; j++ {
            if matrix[i][j] == 0 {
                isBreak = true
            }
        }

        if isBreak {
            break
        }

        cutTop = i + 1

    }

    for i:= height - 1; i> cutTop; i-- {
        isBreak := false
        for j:=0; j<width; j++ {
            if matrix[i][j] == 0 {
                isBreak = true
            }
        }

        if isBreak {
            break
        }

        cutBottom = i

    }

    for j:=0; j<width;j++ {
        isBreak := false
        for i:=0; i<height; i++ {
            if matrix[i][j] == 0 {
                isBreak = true
            }
        }

        if isBreak {
            break
        }

        cutLeft = j + 1

    }

    for j:=width - 1; j> cutLeft;j-- {
        isBreak := false
        for i:=0; i<height; i++ {
            if matrix[i][j] == 0 {
                isBreak = true
            }
        }

        if isBreak {
            break
        }

        cutRight = j

    }

    //fmt.Println(cutTop, cutBottom, cutLeft, cutRight)

    matrix = matrix[cutTop:cutBottom]
    for i := range matrix {
        matrix[i] = matrix[i][cutLeft:cutRight]
    }

    //printMatrix(matrix)


    return matrix

}

func printMatrix(matrix [][]int8) {
    for _, sli := range matrix {
        log.Println(sli)
    }
}

func putCache(matrix [][]int8, value int){
    //todo

}

func getCache(matrix [][]int8) (int, bool) {

    result, ok := cache[CacheKey(matrix)]
    //todo
    return result, ok

}

func CacheKey(matrix [][]int8) string {
    //todo
    return ""
}
