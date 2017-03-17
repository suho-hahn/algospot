package main

import (
    "log"
    "bufio"
    "os"
    "strconv"
    "strings"
    "bytes"
    "math"
    //"runtime/pprof"
    //"flag"
)

type Seq struct {
    Prefix    int
    Next      *Seq
    Length    int
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

var stdoutWriter *bufio.Writer

func main() {

    /*
    max cases: 50
    max strings per case: [1, 15]
    max length of string: [1, 40]
     */

    stdoutWriter = bufio.NewWriter(os.Stdout)
    defer stdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)

    caseNum := scanLineAsInt(scanner)

    for ; caseNum > 0; caseNum -- {

        strNum := scanLineAsInt(scanner)
        strList := make([]string, strNum)

        for i:=0; i<strNum; i++ {
            str := scanLineAsString(scanner)
            strList[i] = str
        }

        solve(strList)

    }


}

type JoinStr struct {

    strList []string
    seqMap  []*Seq
    overlap [][]int

}



func solve(strList []string) {

    //log.Println("---------------------")

    // overlap[left][right] = overlapping length
    joinStr := JoinStr{
        strList: strList,
        //skipList: make([]bool, len(strList)),
        seqMap: make([]*Seq, 1<<uint(len(strList))), // seqList[조합 bitmap][seq 조합 개수만큼]
        overlap: make([][]int, len(strList)),
    }


    for i := 0; i<len(strList); i++ {

        joinStr.overlap[i] = make([]int, len(strList))
        for j := 0; j<len(strList); j++ {
            if i == j {
                continue
            }
            joinStr.overlap[i][j] = calcOverlap(strList[i], strList[j])
        }

        joinStr.seqMap[1<<uint(i)] = &Seq{
            Prefix: i,
            Next: nil,
            Length: len(strList[i]),
        }

    }

    //log.Println(joinStr.seqListMap)


    for i:=1; i<len(joinStr.seqMap); i++ {
        joinStr.FillMapAt(i)
    }

    //for i, seqList := range joinStr.seqListMap {
    //    log.Println("bitmask", i)
    //    for j, seq := range seqList {
    //        result := fmt.Sprint(" result", j, ":")
    //        for seq != nil {
    //            result = fmt.Sprintf("%s %d(%p)", result, seq.Prefix, seq)
    //            seq = seq.Next
    //        }
    //        log.Println(result)
    //
    //    }
    //}

    seq := joinStr.seqMap[len(joinStr.seqMap) - 1]
    stdoutWriter.Write([]byte(strList[seq.Prefix]))

    exPrefix := seq.Prefix
    seq = seq.Next
    for seq != nil {

        curPrefix := seq.Prefix
        overlapLen := joinStr.overlap[exPrefix][curPrefix]

        stdoutWriter.Write([]byte(strList[curPrefix])[overlapLen:])

        exPrefix = curPrefix
        seq = seq.Next
    }

    stdoutWriter.WriteByte('\n')

}

func (join *JoinStr) FillMapAt(mapKey int) {

    if join.seqMap[mapKey] != nil {
        return
    }

    //log.Println("mapKey:", mapKey)

    var resultSeq *Seq
    resultLength := math.MaxInt32

    for prefixStrIndex, prefixBitmask := 0, 1 ;
        prefixBitmask < mapKey ;
        prefixStrIndex, prefixBitmask = prefixStrIndex + 1, prefixBitmask * 2 {

        //log.Println("try", i)

        if mapKey & prefixBitmask == 0 {
            continue
        }

        suffixSeqIndex := mapKey - prefixBitmask
        suffixSeq := join.seqMap[suffixSeqIndex]
        suffixStrIndex := suffixSeq.Prefix

        prefixLen := len(join.strList[prefixStrIndex])
        suffixLen := suffixSeq.Length
        overlapLen := join.overlap[prefixStrIndex][suffixStrIndex]
        concatLen := prefixLen + suffixLen - overlapLen

        if concatLen > resultLength {
            continue
        }

        var newSeq *Seq
        if overlapLen == prefixLen {
            newSeq = suffixSeq
        } else if overlapLen == suffixLen {
            newSeq = &Seq{
                Prefix: prefixStrIndex,
                Next: suffixSeq.Next,
                Length: concatLen,
            }
        } else {
            newSeq = &Seq{
                Prefix: prefixStrIndex,
                Next: suffixSeq,
                Length: concatLen,
            }
        }


        if concatLen < resultLength {
            resultLength = concatLen
            resultSeq = newSeq
        } else if join.CompareSeq(newSeq, resultSeq) == -1 { // ==
            resultSeq = newSeq
        }

    }

    join.seqMap[mapKey] = resultSeq

}

func (join *JoinStr) CompareSeq(seqA, seqB *Seq) int {

    strA := join.strList[seqA.Prefix]
    strB := join.strList[seqB.Prefix]

    for {

        minLen := min(len(strA), len(strB))

        cmp := bytes.Compare([]byte(strA)[:minLen], []byte(strB)[:minLen])
        if cmp != 0 {
            return cmp
        }

        if len(strA) == len(strB) {
            seqA = seqA.Next
            seqB = seqB.Next

            if seqA == nil && seqB == nil {
                return 0
            } else if seqA == nil {
                return -1
            } else {
                return 1
            }

        } else if len(strA) < len(strB) {
            seqA = seqA.Next
            if seqA == nil {
                return -1
            }
            strA = join.strList[seqA.Prefix]
        } else {
            seqB = seqB.Next
            if seqB == nil {
                return 1
            }
            strB = join.strList[seqB.Prefix]
        }

    }


}

func calcOverlap(a, b string) int { // overlap len

    if strings.Index(b, a) != -1 || strings.Index(b, a) != -1  {
        return min(len(a), len(b))
    }

    for i := min(len(a), len(b)) - 1; i >= 1; i-- {
        if bytes.Equal([]byte(a[len(a) - i:]), []byte(b[:i])) {
            return i
        }
    }
    return 0
}

func scanLineAsString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanLineAsInt(scanner *bufio.Scanner) int {
    result, _ := strconv.Atoi(scanLineAsString(scanner))
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