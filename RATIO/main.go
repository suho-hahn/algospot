package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
)

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
        result := solveCase(scanner)


        StdoutWriter.WriteString(strconv.Itoa(result))
        StdoutWriter.WriteByte('\n')
    }

}

func solveCase(scanner *bufio.Scanner) int {

    //log.Println("----------------")

    const maxRemainedGames = 2 * 1000 * 1000 * 1000
    N := scanInt(scanner)
    M := scanInt(scanner)

    //log.Println(N, M)

    ratio := M * 100 / N
    maxRatio := (M + maxRemainedGames) * 100 / (N + maxRemainedGames)

    if ratio == maxRatio {
        return -1
    }

    low := 0
    high := maxRemainedGames

    for {
        mid := (low + high) / 2

        newRatio := (M + mid) * 100 / (N + mid)

        if newRatio > ratio {
            high = mid
        } else if newRatio <= ratio {
            low = mid
        }


        if high - low <= 1 {
            return high
        }


    }

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