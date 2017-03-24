package main

import (
    "bufio"
    "log"
    "os"
    "strconv"
    "math"
    "fmt"
)

type CitySign struct {
    From int
    To   int
    Gap  int
    Count int
}

type CitySignListByFrom []CitySign
type CitySignListByTo []CitySign

var StdoutWriter *bufio.Writer

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {

    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    T := scanAsInt(scanner) // max 50

    for ; T > 0; T -- {
        result := startTest(scanner)
        resultString := fmt.Sprintf("%s\n", strconv.Itoa(result))
        //log.Println(resultString)
        StdoutWriter.WriteString(resultString)
    }

}

func startTest(scanner *bufio.Scanner) int {
    N := scanAsInt(scanner) // 도시의 수 <= 5000
    K := scanAsInt(scanner) // K번째 표지판의 위치 ( 1 ~ 2^31)

    total := 0

    citySignList := make([]CitySign, N)

    minFrom := math.MaxInt32
    maxTo := 0

    for i:=0; i<N; i++ {
        L := scanAsInt(scanner) // 시작점으로부터 각 도시까지의 거리. max 8,030,000
        M := scanAsInt(scanner) // 도시 시작하기 M 미터 전부터
        G := scanAsInt(scanner) // G 미터 간격으로 표지판 설치되어 있음

        citySign := CitySign{
            From: L - M,
            To: L,
            Gap: G,
            Count: M/G + 1,
        }

        citySignList[i] = citySign

        minFrom = min(minFrom, citySign.From)
        maxTo = max(maxTo, citySign.To)

        total += citySign.Count

    }

    low := minFrom - 1
    high := maxTo

    for  {

        mid := (low + 1 + high)/2

        accCount := 0

        for i:=0; i<len(citySignList); i++ {
            accCount += citySignList[i].Passed(mid)
        }

        if accCount < K {

            low = mid

            if high - low == 1 {
                return high
            }

        } else if accCount >= K {

            high = mid

            if high - low == 1 {
                return high
            }

        } 

    }


}

func (citySign *CitySign) Passed(point int) (int) {

    if point > citySign.To {
        return citySign.Count
    } else if point < citySign.From {
        return 0
    } else {
        count := (point - citySign.From) / citySign.Gap + 1
        return count
    }

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
