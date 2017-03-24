package main

import (
    "bufio"
    "log"
    "os"
    "strconv"
    "fmt"
    "sort"
    "math"
)

type Score struct {
    Rank     float64
    Students float64
    V        float64
}

type SortedScore []Score

var StdoutWriter *bufio.Writer

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {

    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    caseNum := scanAsInt(scanner) // max 50

    for ; caseNum > 0; caseNum -- {
        scanCase(scanner)
    }

}

var high float64 = 1.0
var low float64 = 0.0

func scanCase(scanner *bufio.Scanner) {

    //log.Println("--------")

    n := scanAsInt(scanner)
    k := scanAsInt(scanner)

    scoreList := make([]Score, n)
    for i:=0; i<n; i++ {
        scoreList[i].Rank = float64(scanAsInt(scanner))
        scoreList[i].Students = float64(scanAsInt(scanner))
    }

    high = 1.0
    low = 0.0

    solve(scoreList, k, low, high)




}

func solve(scoreList []Score, k int, low, high float64) {

    exMid := -1.0

    for {

        mid := (low + high) / 2

        for i := range scoreList {
            scoreList[i].V = mid * scoreList[i].Students - scoreList[i].Rank
        }

        sort.Sort(SortedScore(scoreList))

        //log.Println(scoreList)

        sumV := 0.0
        for i := 0; i < k; i++ {
            sumV += scoreList[i].V
        }

        //log.Println(mid, sumV)

        //time.Sleep(1 * time.Second)

        if sumV == 0 || math.Abs(mid - exMid) < 0.000000001 {
            StdoutWriter.WriteString(fmt.Sprintf("%.9f\n", mid))
            return
        } else if sumV > 0 {
            high = mid
        } else if sumV < 0 {
            low = mid
        }

        exMid = mid
    }

}

func (a SortedScore) Len() int           { return len(a) }
func (a SortedScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedScore) Less(i, j int) bool {
    return a[i].V > a[j].V
}

func scanAsString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanAsInt(scanner *bufio.Scanner) int {
    result, _ := strconv.Atoi(scanAsString(scanner))
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

func minFloat(a, b float64) float64 {
    if a > b {
        return b
    }
    return a
}