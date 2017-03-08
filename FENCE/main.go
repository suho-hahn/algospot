package main

import (
    "fmt"
    "log"
    "bufio"
    "os"
    "strconv"
)

type MyScanner struct {
    in *bufio.Scanner
}

var in = MyScanner{bufio.NewScanner(os.Stdin)}

func (m *MyScanner) NextInt() (int, error) {
    if !m.innerNext(ScanIntSplitter) {
        return 0, m.in.Err()
    }
    r, _ := strconv.Atoi(m.in.Text())
    return r, nil
}

func (m *MyScanner) innerNext(f bufio.SplitFunc) bool {
    return m.in.Scan()
}

func ScanIntSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
    advance, token, err = bufio.ScanWords(data, atEOF)
    if err == nil && token != nil {
        _, err = strconv.Atoi(string(token))
    }
    return
}

type Fence struct {
    Index   int // 1 <= I <= 20,000
    Height  int // 1 <= H <= 10,000
}

type Rect struct {
    StartIndex int
    EndIndex int
    Height int
}

type RectResult struct {
    LeftJoin []*Rect
    RightJoin []*Rect
    BothJoin []*Rect
    MaxSize int
}

func NewRect(start, end, height int) *Rect {
    return &Rect{
        StartIndex: start,
        EndIndex: end,
        Height: height,
    }
}

func (r *Rect) GetSize() int{
    return r.Height * (r.EndIndex - r.StartIndex + 1)
}

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

    in.in.Split(ScanIntSplitter)
    // go 1.2.1에는 없나보다...ㅠ
    //in.in.Buffer(make([]byte, 1<<20), 1<<20)

    caseNum, _ := in.NextInt() // <= 50

    for i:=0; i<caseNum; i++ {
        fenceLen, _ := in.NextInt()
        fences := make([]int, fenceLen)
        for j:=0; j<fenceLen; j++ {
            fences[j], _ = in.NextInt()
        }

        fmt.Println(getMaxRectangleSize(fences))

    }

}

func getMaxRectangleSize(fences []int) (int) { //height, start, end, curMaxSize,

    rectResult := getRects(fences, 0, len(fences) - 1)
    result := rectResult.MaxSize

    for _, rect := range rectResult.LeftJoin {
        result = max(result, rect.GetSize())
    }
    for _, rect := range rectResult.RightJoin {
        result = max(result, rect.GetSize())
    }
    for _, rect := range rectResult.BothJoin {
        result = max(result, rect.GetSize())
    }

    return result

}

func getRects(fences []int, start, end int) RectResult {

    if start == end {
        return RectResult{
            LeftJoin: nil,
            RightJoin: nil,
            BothJoin: []*Rect{NewRect(start, end, fences[start])},
            MaxSize: 0,
        }
    }

    midIdx := (start + end) / 2
    leftResult := getRects(fences, start, midIdx)
    rightResult := getRects(fences, midIdx + 1, end)

    //log.Println("start/end:", start, end)
    //log.Println("left :", *leftResult.BothJoin[0], leftResult.MaxSize, leftResult)
    //log.Println("right:", *rightResult.BothJoin[0], rightResult.MaxSize, rightResult)

    joinA := append(leftResult.BothJoin, leftResult.RightJoin...)
    joinB := append(rightResult.BothJoin, rightResult.LeftJoin...)

    result := RectResult{
        BothJoin: make([]*Rect, 0, len(joinA) + len(joinB)),
        LeftJoin: make([]*Rect, 0, len(leftResult.LeftJoin) + len(leftResult.BothJoin)),
        RightJoin: make([]*Rect, 0, len(rightResult.RightJoin) + len(rightResult.BothJoin)),
        MaxSize: max(leftResult.MaxSize, rightResult.MaxSize),
    }

    for idxA, idxB := 0, 0 ; ; {
        rectA, rectB := joinA[idxA], joinB[idxB]
        //log.Println(*rectA, *rectB)
        var appendTarget *Rect
        if rectA.Height < rectB.Height {
            rectA.EndIndex = rectB.EndIndex
            appendTarget = rectA
            idxA ++
        } else if rectA.Height > rectB.Height {
            rectB.StartIndex = rectA.StartIndex
            appendTarget = rectB
            idxB ++
        } else {
            rectA.EndIndex = rectB.EndIndex
            appendTarget = rectA
            idxA ++
            idxB ++
        }

        if appendTarget.StartIndex == start && appendTarget.EndIndex == end {
            //log.Println("append both:", *appendTarget)
            result.BothJoin = append(result.BothJoin, appendTarget)
        } else if appendTarget.StartIndex == start {
            //log.Println("append left:", *appendTarget)
            result.LeftJoin = append(result.LeftJoin, appendTarget)
        } else if appendTarget.EndIndex == end {
            //log.Println("append right:", *appendTarget)
            result.RightJoin = append(result.RightJoin, appendTarget)
        } else {
            //log.Println("append none:", *appendTarget)
            result.MaxSize = max(result.MaxSize, appendTarget.GetSize())
        }

        if idxA == len(joinA) {
            for _, rectB := range joinB[idxB:] {
                if rectB.EndIndex == end {
                    result.RightJoin = append(result.RightJoin, rectB)
                } else {
                    result.MaxSize = max(result.MaxSize, rectB.GetSize())
                }
            }
            break
        }
        if idxB == len(joinB) {
            for _, rectA := range joinA[idxA:] {
                if rectA.StartIndex == start {
                    result.LeftJoin = append(result.LeftJoin, rectA)
                } else {
                    result.MaxSize= max(result.MaxSize, rectA.GetSize())
                }
            }
            break
        }

    }

    result.LeftJoin = append(result.LeftJoin, leftResult.LeftJoin...)
    result.RightJoin = append(result.RightJoin, rightResult.RightJoin...)

    return result

}

func max(i, j int) int {
    if  i >= j {
        return i
    }
    return j
}
