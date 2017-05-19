package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "container/heap"
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

    caseNum := scanInt(scanner) // max 20

    for ; caseNum > 0; caseNum -- {


        N := scanInt(scanner)
        a := scanInt(scanner)
        b := scanInt(scanner)

        result := solve(N, a, b)

        StdoutWriter.WriteString(strconv.Itoa(result))
        StdoutWriter.WriteByte('\n')

    }

}

type MinHeap []int
type MaxHeap []int

func solve(N, a, b int) int {

    A := 1983
    sum := 0

    minHeapArr := make([]int, 0, N)
    maxHeapArr := make([]int, 0, N)

    smaller := (*MaxHeap)(&maxHeapArr)
    larger := (*MinHeap)(&minHeapArr)

    heap.Init(smaller)
    heap.Init(larger)

    heap.Push(smaller, A)
    sum += A

    for i :=1; i <N; i++ {

        // new input value
        A = (A * a + b) % 20090711

        if A > (*smaller)[0] {
            heap.Push(larger, A)
        } else {
            heap.Push(smaller, A)
        }

        for smaller.Len() > larger.Len() + 1 {
            moveValue := heap.Pop(smaller)
            //log.Println(moveValue)
            heap.Push(larger, moveValue)
        }

        for smaller.Len() < larger.Len() {
            moveValue := heap.Pop(larger)
            //log.Println(moveValue)
            heap.Push(smaller, moveValue)
        }

        //log.Println(smaller, larger)
        sum += (*smaller)[0]
        sum %= 20090711




    }

    return sum

}

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func scanInt(scanner *bufio.Scanner) int {
    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
    return result
}
