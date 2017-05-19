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

    const qSize = 5000001

    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    caseNum := scanInt(scanner) // max 20
    queue := make([]int, qSize)

    for ; caseNum > 0; caseNum -- {

        queue[0] = 0
        head := 0
        tail := 1

        K := scanInt(scanner)
        N := scanInt(scanner)

        A := 1983
        sum := 0
        result := 0

        for i :=0; i <=N; i++ {
            // get sig
            sig := A % 10000 + 1 // 1 ~ 10,000
            sum += sig

            for ; head < tail; head ++ {
                rangeValue := sum - queue[head % qSize]
                if rangeValue < K {
                    break
                } else if rangeValue == K {
                    result ++
                }
            }

            queue[tail % qSize] = sum
            tail ++

            // calc next A
            A = (A * 214013 + 2531011) % (1<<32)

        }

        StdoutWriter.WriteString(strconv.Itoa(result))
        StdoutWriter.WriteByte('\n')

    }

}

func scanInt(scanner *bufio.Scanner) int {
    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
    return result
}
