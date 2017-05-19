package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
)

type Soldier struct {
    Index int
    //Prev *Soldier
    Next *Soldier
}

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
    K := scanInt(scanner)

    first := &Soldier {
        Index: 1,
        Next: nil,
    }
    current := first
    for i:=2; i<=N; i++ {
        current.Next = &Soldier{
            Index: i,
            Next: nil,
        }

        current = current.Next
    }
    current.Next = first

    //current = first
    //for i:=0; i<N; i++ {
    //    log.Println(current.Prev.Index, current.Index, current.Next.Index)
    //    current = current.Next
    //}


    alive := N
    for {
        remove := current.Next
        alive --

        //log.Println(remove.Index)
        current.Next = remove.Next

        if alive <= 2 {
            break
        }

        for i:=0; i<K-1; i++ {
            current = current.Next
        }

    }

    a, b := current.Index, current.Next.Index

    if a <= b {
        StdoutWriter.WriteString(strconv.Itoa(a))
        StdoutWriter.WriteByte(' ')
        StdoutWriter.WriteString(strconv.Itoa(b))
    } else {
        StdoutWriter.WriteString(strconv.Itoa(b))
        StdoutWriter.WriteByte(' ')
        StdoutWriter.WriteString(strconv.Itoa(a))
    }

    StdoutWriter.WriteByte('\n')

}

func scanInt(scanner *bufio.Scanner) int {
    scanner.Scan()
    result, _ := strconv.Atoi(scanner.Text())
    return result
}
