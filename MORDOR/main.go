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
        solveCase(scanner)
    }

}

type Query struct {
    Start int
    End int
}

func solveCase(scanner *bufio.Scanner) {

    signCount := scanInt(scanner)
    queryCount := scanInt(scanner)

    adjustedSignCount := ceilSize(signCount)
    minSignList := make([]int, adjustedSignCount)
    maxSignList := make([]int, adjustedSignCount)

    memsetRepeat(minSignList, math.MaxInt32)
    memsetRepeat(maxSignList, math.MinInt32)

    for i:=0; i< signCount; i++ {
        sign := scanInt(scanner)
        minSignList[i] = sign
        maxSignList[i] = sign
    }


    treeSize := adjustedSignCount * 2
    minTree := make([]int, treeSize)
    maxTree := make([]int, treeSize)

    makeRangeMin(minTree, minSignList, 1, 0, adjustedSignCount)
    makeRangeMax(maxTree, maxSignList, 1, 0, adjustedSignCount)

    //log.Println(minTree)
    //log.Println(maxTree)
    //log.Println(signCount)

    for i:=0; i<queryCount; i++ {
        start := scanInt(scanner)
        end := scanInt(scanner)

        resultMin := treeMin(minTree, 1, 0, adjustedSignCount - 1, start, end)
        //log.Println("min:", resultMin)

        resultMax := treeMax(maxTree, 1, 0, adjustedSignCount - 1, start, end)
        //log.Println("max:", resultMax)

        //log.Println("result:", resultMax - resultMin)

        StdoutWriter.WriteString(strconv.Itoa(resultMax - resultMin))
        StdoutWriter.WriteByte('\n')


    }




}

func treeMin(tree []int, nodeIdx, nodeStart, nodeEnd, queryStart, queryEnd int) (result int) {

    //log.Println(nodeIdx, nodeStart, nodeEnd, queryStart, queryEnd)

    if nodeStart >= queryStart && nodeEnd <= queryEnd {
        return tree[nodeIdx]
    } else if nodeStart > queryEnd || nodeEnd < queryStart {
        return math.MaxInt32
    } else {
        mid := (nodeStart + nodeEnd) / 2
        leftMin := treeMin(tree, nodeIdx * 2, nodeStart, mid, queryStart, queryEnd)
        rightMin := treeMin(tree, nodeIdx * 2 + 1, mid + 1, nodeEnd, queryStart, queryEnd)
        return min(leftMin, rightMin)
    }

}

func treeMax(tree []int, nodeIdx, nodeStart, nodeEnd, queryStart, queryEnd int) (result int) {

    //log.Println(nodeIdx, nodeStart, nodeEnd, queryStart, queryEnd)


    if nodeStart >= queryStart && nodeEnd <= queryEnd {
        return tree[nodeIdx]
    } else if nodeStart > queryEnd || nodeEnd < queryStart {
        return math.MinInt32
    } else {
        mid := (nodeStart + nodeEnd) / 2
        leftMin := treeMax(tree, nodeIdx * 2, nodeStart, mid, queryStart, queryEnd)
        rightMin := treeMax(tree, nodeIdx * 2 + 1, mid + 1, nodeEnd, queryStart, queryEnd)
        return max(leftMin, rightMin)
    }

}

func makeRangeMin(tree, values []int, nodeIdx, from, len int) int {

    //defer func(){
    //    log.Println(from, len, tree[nodeIdx])
    //}()

    if len == 1 {
        tree[nodeIdx] = values[from]
        return tree[nodeIdx]
    }

    leftMin := makeRangeMin(tree, values, nodeIdx * 2, from, len/2)
    rightMin := makeRangeMin(tree, values, nodeIdx * 2 + 1, from + len/2, len/2)

    tree[nodeIdx] = min(leftMin, rightMin)
    return tree[nodeIdx]

}

func makeRangeMax(tree, values []int, nodeIdx, from, len int) int {

    //defer func(){
    //    log.Println(from, len, tree[nodeIdx])
    //}()


    if len == 1 {
        tree[nodeIdx] = values[from]
        return tree[nodeIdx]
    }

    leftMax := makeRangeMax(tree, values, nodeIdx * 2, from, len/2)
    rightMax := makeRangeMax(tree, values, nodeIdx * 2 + 1, from + len/2, len/2)

    tree[nodeIdx] = max(leftMax, rightMax)
    return tree[nodeIdx]

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

func ceilSize(size int) int {

    properSize := 1
    for properSize < size {
        properSize *= 2
    }

    return properSize

}

func scanInt(scanner *bufio.Scanner) int {
    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
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