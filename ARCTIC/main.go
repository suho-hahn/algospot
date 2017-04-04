package main

import (
    "bufio"
    "log"
    "os"
    "strconv"
    "math"
    "fmt"
)

type Point struct {
    X float64
    Y float64
    Visit bool
}

var StdoutWriter *bufio.Writer

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {

    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    caseNum := scanInt(scanner) // max 50

    for ; caseNum > 0; caseNum -- {
        SolveCase(scanner)
    }

}

func SolveCase(scanner *bufio.Scanner) {

    N := scanInt(scanner)

    pointList := make([]Point, N)

    for i:=0; i<N; i++ {
        pointList[i].X = scanFloat(scanner)
        pointList[i].Y = scanFloat(scanner)
    }

    dist := make([][]float64, N)
    for i:=0; i<N; i++ {
        dist[i] = make([]float64, N)
        for j:=0; j<N; j++ {
            dist[i][j] = pointList[i].Distance(pointList[j])
        }
    }

    low := 0.0
    high := 1420.0

    for {

        mid := (low + high) / 2

        Span(pointList, dist, mid, 0, 0)


        impossible := false
        for i:=0; i<N; i++ {
            if pointList[i].Visit == false {
                impossible = true
            }
            pointList[i].Visit = false
        }


        if ! impossible {
            high = mid
        } else {
            low = mid
        }

        if high - low < 0.000000001 {
            resultStr := fmt.Sprintf("%.02f\n", high)
            StdoutWriter.WriteString(resultStr)
            return
        }

    }


}

func Span(pointList []Point, distCache [][]float64, maxDist float64, curPointIdx int, depth int) {

    //log.Println(depth, len(unreachedMap))

    pointList[curPointIdx].Visit = true

    for i := range pointList {

        if pointList[i].Visit {
            continue
        }

        if distCache[curPointIdx][i] <= maxDist {
            Span(pointList, distCache, maxDist, i, depth + 1)
        }
    }

}

func (p1 Point) Distance(p2 Point) float64{
    X := math.Abs(p1.X - p2.X)
    Y := math.Abs(p1.Y - p2.Y)

    return math.Sqrt(X*X + Y*Y)
}


func scanString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanInt(scanner *bufio.Scanner) int {
    result, _ := strconv.Atoi(scanString(scanner))
    return result
}

func scanFloat(scanner *bufio.Scanner) float64 {
    result, _ := strconv.ParseFloat(scanString(scanner), 64)
    return result
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a > b {
        return b
    }
    return a
}