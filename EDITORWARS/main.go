package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "fmt"
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
        scanCase(scanner)
    }

}

type Node struct {
    //Index  int
    Ack    *Node
    Dis    *Node

    Height int // root부터 leaf까지 최대 edge 개수
    Size   int
    Counted bool
}

type NodeList []Node

func scanCase(scanner *bufio.Scanner) {

    //log.Println("===========================================")

    N := scanInt(scanner) // 사람 수 (1≤N≤10000)
    M := scanInt(scanner) // 댓글 수 (1≤M≤100000)

    //log.Println(N, M)
    // 회원 ID: 0에서 N - 1 범위

    users := NodeList(make([]Node, N))
    //comments := make([]Comment, M)
    //for i := range users {
    //    users[i] = Node{
    //        //Index: i,
    //        //Size:1,
    //
    //    }
    //}

    //minSizeTree := 0
    contradiction := -1

    for i:=0; i<M; i++ {

        commentType := scanString(scanner)
        indexA, indexB := scanInt(scanner), scanInt(scanner)

        //log.Println("---------------------------------")
        //log.Println("input:", commentType, indexA, indexB)

        ackA := (&users[indexA]).FindAck() // never null
        ackB := (&users[indexB]).FindAck() // never null

        disA := ackA.FindDis() // nullable
        disB := ackB.FindDis() // nullable

        if commentType[0] == 'A' { // ACK


            if (disB != nil && ackA == disB) || (disA != nil && ackB == disA) {
                contradiction = i
                break
            }

            // ack 끼리 합침
            ackA.Merge(ackB)

            // dis끼리 합침
            if disA != nil && disB != nil  {
                disA.Merge(disB)
            } else if disA == nil {
                ackA.Dis = disB
            } else if ackB == nil {
                ackB.Dis = disA
            }

        } else { // DIS

            if ackA == ackB || (disA == disB && disA != nil && disB != nil) {
                contradiction = i
                break
            }

            // a의 반대파와 b의 찬성파 합침
            if ackA.Dis == nil {
                ackA.Dis = ackB
            } else {
                ackA.Dis.Merge(ackB)
            }

            // b의 반대파와 a의 찬성파 합침
            if ackB.Dis == nil {
                ackB.Dis = ackA
            } else {
                ackB.Dis.Merge(ackA)
            }

            //minSizeTree = max(min(rootA.Size, rootB.Size), minSizeTree)

        }

        //users.Print()

    }

    if contradiction >= 0{
        contradiction ++
        for i:=contradiction; i<M; i++ {
            scanString(scanner)
            scanInt(scanner)
            scanInt(scanner)
        }
        fmt.Println("CONTRADICTION AT", contradiction)

    } else {

        result := 0

        for _, node := range users {

            find := node.FindAck()

            if find.Dis == nil {
                continue
            }

            dis := find.FindDis()

            if find.Counted || dis.Counted {
                continue
            }

            result += min(find.Size, dis.Size) + 1
            find.Counted = true
            dis.Counted = true
        }

        result = N - result

        fmt.Println("MAX PARTY SIZE IS", result)
    }

    // MAX PARTY SIZE IS 3
    // MAX PARTY SIZE IS 100
    // CONTRADICTION AT 3
    //StdoutWriter.WriteString(strconv.Itoa(maxTime))
    //StdoutWriter.WriteByte('\n')

}

func (node *Node) FindAck() *Node {

    if node.Ack == nil {
        return node
    }

    ack := node.Ack
    result := ack.FindAck()
    node.Ack = result
    return result

}

func (node *Node) FindDis() *Node {

    node = node.FindAck()

    if node.Dis == nil {
        return nil
    }

    dis := node.Dis
    result := dis.FindAck()
    node.Dis = result
    return result

}

//func (nodes NodeList) Print() {
//    for i, n := range nodes {
//
//        ackIndex := -1
//        disIndex := -1
//
//        if n.Ack != nil {
//            ackIndex = n.Ack.Index
//        }
//
//        if n.Dis != nil {
//            disIndex = n.Dis.Index
//        }
//
//        log.Println("index:", i, "| ACK:", ackIndex, "| DIS:", disIndex)
//    }
//}

func (a *Node) Merge(b *Node) {

    //log.Println("Merge:", a.Index, b.Index)

    if a.Height < b.Height {
        b.Merge(a)
        return
    }

    if a == b {
        return
    }

    if a.Height == b.Height {
        a.Height ++
    }

    b.Ack = a
    a.Size += b.Size + 1
    return

}

func scanString(scanner *bufio.Scanner) string {
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