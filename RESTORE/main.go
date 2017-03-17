package main

import (
    "log"
    "bufio"
    "os"
    "strconv"
    "strings"
    "bytes"
    "fmt"
    "math"
    //"flag"
    //"runtime/pprof"
)

type Seq struct {
    Prefix    int
    Next      *Seq
    Length    int
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

//var memprofile = flag.String("memprofile", "", "write memory profile to this file")
//var profIdx = 0

func main() {
    //flag.Parse()
    /*
    max cases: 50
    max strings per case: [1, 15]
    max length of string: [1, 40]
     */

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
    resultString := make([]byte, 0, seq.Length)
    resultString = append(resultString, []byte(strList[seq.Prefix]) ...)

    exPrefix := seq.Prefix
    seq = seq.Next
    for seq != nil {

        curPrefix := seq.Prefix
        overlapLen := joinStr.overlap[exPrefix][curPrefix]

        resultString = append(
            resultString[:len(resultString) - overlapLen],
            []byte(strList[curPrefix])...)

        exPrefix = curPrefix
        seq = seq.Next
    }

    fmt.Println(string(resultString))

    //if *memprofile != "" {
    //
    //    profIdx ++
    //    f, err := os.Create(fmt.Sprint(*memprofile, "_", profIdx))
    //
    //    if err != nil {
    //        log.Fatal(err)
    //    }
    //
    //    pprof.WriteHeapProfile(f)
    //    f.Close()
    //    return
    //}


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

        if concatLen < resultLength {
            resultLength = concatLen
        } else if concatLen >= resultLength {
            continue
        }

        if overlapLen == prefixLen {
            resultSeq = suffixSeq
        } else if overlapLen == suffixLen {
            newSeq := &Seq{
                Prefix: prefixStrIndex,
                Next: suffixSeq.Next,
                Length: concatLen,
            }
            resultSeq = newSeq
        } else {
            newSeq := &Seq{
                Prefix: prefixStrIndex,
                Next: suffixSeq,
                Length: concatLen,
            }
            resultSeq = newSeq
        }
    }

    join.seqMap[mapKey] = resultSeq

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