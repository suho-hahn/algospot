package main

import (
    "log"
    "bufio"
    "os"
    "strconv"
    "math"
    "strings"
    "fmt"
    "bytes"
)

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {

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
        strList := make([]string, 1 << uint(strNum))

        for i:=0; i<strNum; i++ {
            str := scanLineAsString(scanner)
            strList[1 << uint(i)] = str
        }

        solve(strList)

    }

}

func solve(strList []string) {

    //log.Println("---------------------")
    fmt.Println(getStringAt(strList, len(strList) - 1))

}

func getStringAt(strList []string, index int) string {
    if len(strList[index]) > 0 {
        return strList[index]
    }

    //log.Println("index:", index)

    result := ""
    minLength := math.MaxInt32
    for i:= 1; i<index ; i *= 2 {

        if index & i == 0 {
            continue
        }
        j := index - i

        a := getStringAt(strList, i)
        b := getStringAt(strList, j)

        joinResult, ok := join(a, b, minLength)

        if ! ok {
            continue
        }

        //log.Println(index, i, j, joinResult)
        if len(joinResult) < minLength {
            result = joinResult
            minLength = len(result)
        }

    }

    strList[index] = result

    //log.Println(index, ":", result)

    return result

}

func join(a, b string, maxLen int) (string, bool) {

    if strings.Index(b, a) != -1 {
        return b, true
    } else if strings.Index(a, b) != -1 {
        return a, true
    }

    lenSum := len(a) + len(b)

    // always len(a) > len(b)
    for i:= min(len(a), len(b)) - 1; i >= 0; i-- {

        if lenSum - i > maxLen {
            return "", false
        }

        // concat AB
        if bytes.Equal([]byte(a[len(a) - i:]), []byte(b[:i])) {
            concatAB := string(append([]byte(a), b[i:]...))
            return concatAB, true
        }
        // concat BA
        //if bytes.Equal([]byte(b[len(b) - i:]), []byte(a[:i])) {
        //    concatBA := string(append([]byte(b), a[i:]...))
        //    return concatBA, true
        //}
    }

    return "", false

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