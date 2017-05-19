package main

import (
    "bufio"
    "log"
    "os"
    "strconv"
    "fmt"
    "math"
    "container/heap"
)

type Edge struct {
    To int
    Noise float64
}


type Node struct {
    Dist float64
}

type SortedEdgeListByFrom []Edge

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

func scanCase(scanner *bufio.Scanner) {

    N := scanInt(scanner)
    M := scanInt(scanner)

    edges := make([][]Edge, N)
    nodes := make([]Node, N)

    for i:=0; i<M; i++ {

        a, b, c := scanInt(scanner), scanInt(scanner), scanFloat(scanner)

        edges[a] = append(edges[a], Edge{
            To: b,
            Noise: c,
        })

        edges[b] = append(edges[b], Edge{
            To: a,
            Noise: c,
        })
    }

    for i:=range nodes{
        nodes[i].Dist = math.MaxFloat64
    }
    nodes[0].Dist = 1.0

    Dijkstra(edges, nodes, N)

    resultStr := fmt.Sprintf("%.10f\n", nodes[N-1].Dist)
    //log.Println(iterCount)
    StdoutWriter.WriteString(resultStr)

}

var iterCount = 0

func Dijkstra(edges [][]Edge, nodes []Node, N int) {

    queue := make([]int, 1, N)
    infMap := make(map[int]bool)
    for i:=1; i<N; i ++ {
        infMap[i] = true
    }
    queue[0] = 0


    for len(queue) > 0 && len(infMap) > 0 {

        minDist := math.MaxFloat64
        minNodeIdx := 0
        minQueueIdx := 0

        // find min d[v]
        for queueIdx :=0; queueIdx <len(queue); queueIdx++ {
            nodeIdx := queue[queueIdx]
            if minDist > nodes[nodeIdx].Dist {
                minDist = nodes[nodeIdx].Dist
                minNodeIdx = nodeIdx
                minQueueIdx = queueIdx
            }
        }

        //log.Println(len(nodes), len(queue), minDist, minNodeIdx, minQueueIdx)
        queue = append(queue[:minQueueIdx], queue[minQueueIdx+1:]...)

        for i:=0; i<len(edges[minNodeIdx]); i++ {
            edge := edges[minNodeIdx][i]
            alt := nodes[minNodeIdx].Dist * edge.Noise

            if alt < nodes[edge.To].Dist {
                nodes[edge.To].Dist = alt

                if _, ok := infMap[edge.To]; ok {
                    queue = append(queue, edge.To)
                    delete(infMap, edge.To)
                }

            }

        }

    }



    return

}

func logPrefix(depth int) string {
    result := ""
    for i:=0; i<depth; i++ {
        result += ">"
    }
    return fmt.Sprintf("\t|%s[%d]", result, depth)
}
//func (a SortedEdgeListByFrom) Len() int           { return len(a) }
//func (a SortedEdgeListByFrom) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a SortedEdgeListByFrom) Less(i, j int) bool {
//    return a[i].From < a[j].From
//}
func memsetRepeat(a []float64, v float64) {
    if len(a) == 0 {
        return
    }
    a[0] = v
    for bp := 1; bp < len(a); bp *= 2 {
        copy(a[bp:], a[:bp])
    }
}

func scanString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanInt(scanner *bufio.Scanner) int {
    result, _ := strconv.ParseInt(scanString(scanner), 10, 0)
    return int(result)
}

func scanFloat(scanner *bufio.Scanner) float64 {
    result, _ := strconv.ParseFloat(scanString(scanner), 64)
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
