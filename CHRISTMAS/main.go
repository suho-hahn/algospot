package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "math"
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
        scanCase(scanner) // <=60
    }

}

func scanCase(scanner *bufio.Scanner) {

    //log.Println("----------------")

    const MOD = 20091101

    N := scanInt(scanner) //<=100,000
    K := scanInt(scanner) //<=100,000

    psum := make([]int, N) // index: sum length from index 0
    psum[0] = scanInt(scanner) % K

    modCount := make([]int, K) // [mod value]{ count, count }
    modLastIndex := make([]int, K) // [mod value]{ index, index }
    memsetRepeat(modLastIndex, math.MinInt64)

    modCount[0] = 1
    modLastIndex[0] = -1

    prevManyBuyIndex := -1
    manyBuyCount := 0

    modCount[psum[0]] ++
    if prevManyBuyIndex <= modLastIndex[psum[0]] {
        prevManyBuyIndex = 0
        manyBuyCount ++
    }
    modLastIndex[psum[0]] = 0

    for i:=1; i<N; i++ {
        psum[i] = (psum[i-1] + scanInt(scanner)) % K

        modCount[psum[i]] ++

        if prevManyBuyIndex <= modLastIndex[psum[i]] {
            prevManyBuyIndex = i
            manyBuyCount ++
        }
        modLastIndex[psum[i]] = i

    }

    buyCombiCount := 0
    for i:=0; i<K; i++ {
        if modCount[i] >= 2 {
            buyCombiCount += (modCount[i] - 1) * modCount[i] / 2
            buyCombiCount %= MOD
        }
    }

    StdoutWriter.WriteString(strconv.Itoa(buyCombiCount))
    StdoutWriter.WriteByte(' ')
    StdoutWriter.WriteString(strconv.Itoa(manyBuyCount))
    StdoutWriter.WriteByte('\n')

}

func memsetRepeat(a []int, v int) {
    if len(a) == 0 {
        return
    }
    a[0] = v
    for bp := 1; bp < len(a); bp *= 2 {
        copy(a[bp:], a[:bp])
    }
}

func scanInt(scanner *bufio.Scanner) int {

    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
    return result
}
