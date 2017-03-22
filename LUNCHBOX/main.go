package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "sort"
)

type Box struct {
    HeatTime int
    EatTime int
}

type SortedBoxList []Box

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

func scanCase(scanner *bufio.Scanner) {

    //log.Println("----------------")

    boxCount := scanAsInt(scanner)
    boxList := make([]Box, boxCount)

    for i:=0; i<boxCount; i++ {
        boxList[i].HeatTime = scanAsInt(scanner)
    }

    for i:=0; i<boxCount; i++ {
        boxList[i].EatTime = scanAsInt(scanner)
    }

    // TODO

    sort.Sort(SortedBoxList(boxList))
    maxTime := 0
    curTime := 0
    for i:=0; i<len(boxList); i++ {
        box := boxList[i]
        //log.Println(box)

        maxTime = max(maxTime, curTime + box.HeatTime + box.EatTime)
        curTime = curTime + box.HeatTime

    }

    StdoutWriter.WriteString(strconv.Itoa(maxTime))
    StdoutWriter.WriteByte('\n')

}

func (a SortedBoxList) Len() int           { return len(a) }
func (a SortedBoxList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedBoxList) Less(i, j int) bool {
    if a[i].EatTime > a[j].EatTime {
        return true
    } else if a[i].EatTime < a[j].EatTime {
        return false
    }

    return a[i].HeatTime < a[j].HeatTime
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