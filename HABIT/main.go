package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "sort"
)

type StringList []string

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
        scanCase(scanner)
    }

}

var result int

func scanCase(scanner *bufio.Scanner) {

    //log.Println("----------------")
    result = 0
    minGroupCount := scanInt(scanner)
    str := scanString(scanner)
    substrList := make([]string, len(str))

    for i:=0; i<len(str); i++ {
        substrList[i] = str[i:]
    }

    sort.Strings(substrList)
    StringList(substrList).GroupByPrefix(minGroupCount, 0)
    //log.Println("Result:", result)
    StdoutWriter.WriteString(strconv.Itoa(result))
    StdoutWriter.WriteByte('\n')

}

func (strList StringList) GroupByPrefix(minGroupCount int, prefixLen int) {

    //log.Println(Indent(depth), strList)

    prefixByte := byte(' ')
    groupStart := 0
    groupedCount := 0

    var i int = 0
    for i=0; i<len(strList); i++ {
        if len(strList[i]) > prefixLen {
            break
        }
    }

    for ; i<len(strList); i++ {

        str := strList[i]

        if prefixByte == str[prefixLen] {
            groupedCount ++
            continue
        }

        if groupedCount >= minGroupCount {
            result = max(result, prefixLen + 1)
            StringList(strList[groupStart:i]).GroupByPrefix(minGroupCount, prefixLen + 1)
        }

        prefixByte = str[prefixLen]
        groupedCount = 1
        groupStart = i

    }

    if groupedCount >= minGroupCount {
        result = max(result, prefixLen + 1)
        StringList(strList[groupStart:]).GroupByPrefix(minGroupCount, prefixLen + 1)
    }

}

func Indent(depth int) string {
    if depth == 0 {
        return ">"
    }
    return " " + Indent(depth - 1)
}

func (a StringList) Len() int           { return len(a) }
func (a StringList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StringList) Less(i, j int) bool {
    return a[i][0] < a[j][0]
}

func scanString(scanner *bufio.Scanner) string {
    //if ! scanner.Scan() {
    //    log.Println(scanner.Err())
    //    panic(scanner.Err())
    //}

    scanner.Scan()
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