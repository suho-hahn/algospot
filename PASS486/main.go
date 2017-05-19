package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    //"time"
    "math"
    //"sort"
)

type Factor struct {
    Number int32
    Count  uint16
}

var StdoutWriter *bufio.Writer

const maxDivisor = 1000 * 1000 * 10
var factorCount = make([]uint16, maxDivisor + 1)
//var factorSet = make(map[uint16][]int)

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)

    // init divisor count

    //start := time.Now()
    for i:=1; i<=maxDivisor/2; i++ {
        for j:=i; j<=maxDivisor; j += i {
            if factorCount[j] < math.MaxUint16{
                factorCount[j] ++
            }
        }

        //factorSet[factorCount[i]] = append(factorSet[factorCount[i]], i)
    }

    for i:=maxDivisor/2 + 1; i<=maxDivisor; i++ {
        factorCount[i] ++
        //factorSet[factorCount[i]] = append(factorSet[factorCount[i]], i)
    }

    //for k := range factorSet {
    //    slice := factorSet[k]
    //    sort.Ints(slice)
    //
    //    //log.Println(k, factorSet[k])
    //}

    //end := time.Now()
    //log.Println("precalc elapse:", (end.UnixNano() - start.UnixNano()) / 1000 / 1000 )

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

    num := uint16(scanInt(scanner)) // < 400
    low := scanInt(scanner) //
    high := scanInt(scanner) //

    result := 0
    for i:=low; i<=high; i++ {
        if factorCount[i] == num {
            result ++
        }
    }

    return result

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

func max(a, b uint16) uint16{
    if a > b {
        return a
    }
    return b
}