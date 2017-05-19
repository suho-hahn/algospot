package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
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
    Value int
    Left *Node
    Right *Node
}

func scanCase(scanner *bufio.Scanner) {

    //log.Println("----------------")

    nodeCount := scanInt(scanner)

    preorder := make([]int, nodeCount)
    inorder := make([]int, nodeCount)

    for i:=0; i<nodeCount; i++ {
        preorder[i] = scanInt(scanner)
    }

    for i:=0; i<nodeCount; i++ {
        inorder[i] = scanInt(scanner)
    }

    tree := Fill(preorder, inorder)

    //log.Println("====================================")
    //result := make([]int, 0, nodeCount)
    //tree.TraversePre(&result)
    //log.Println(result)
    //log.Println("--------------------")
    //result = make([]int, 0, nodeCount)
    //tree.TraverseIn(&result)
    //log.Println(result)

    result := make([]int, 0, nodeCount)
    tree.TraversePost(&result)

    for _, value := range result {
        StdoutWriter.WriteString(strconv.Itoa(value))
        StdoutWriter.WriteByte(' ')
    }
    StdoutWriter.WriteByte('\n')

}

func Fill(preorder, inorder []int) *Node {

    if len(preorder) == 0 {
        return nil
    }

    curNodeValue := preorder[0]

    inorderLeftLen := 0
    for inorder[inorderLeftLen] != curNodeValue {
        inorderLeftLen ++
    }

    preorderLeft := preorder[1:inorderLeftLen + 1]
    preorderRight := preorder[inorderLeftLen + 1 :]

    inorderLeft := inorder[:inorderLeftLen]
    inorderRight := inorder[inorderLeftLen + 1 :]

    node := &Node {
        Value: curNodeValue,
        Left: Fill(preorderLeft, inorderLeft),
        Right: Fill(preorderRight, inorderRight),
    }

    return node

}

func (n *Node) TraversePre(result *[]int) {

    *result = append(*result, n.Value)

    if n.Left != nil {
        n.Left.TraversePre(result)
    }

    if n.Right != nil {
        n.Right.TraversePre(result)
    }

}

func (n *Node) TraverseIn(result *[]int) {

    if n.Left != nil {
        n.Left.TraverseIn(result)
    }

    *result = append(*result, n.Value)

    if n.Right != nil {
        n.Right.TraverseIn(result)
    }

}

func (n *Node) TraversePost(result *[]int) {

    if n.Left != nil {
        n.Left.TraversePost(result)
    }

    if n.Right != nil {
        n.Right.TraversePost(result)
    }

    *result = append(*result, n.Value)

}

func scanInt(scanner *bufio.Scanner) int {
    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
    return result
}