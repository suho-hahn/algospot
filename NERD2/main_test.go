package main

import (
    "testing"
    "fmt"
    "math/rand"
)

func Test_GenCase1(t *testing.T) {

    p, q := 1, 100000

    for i:=0; i<50000; i++ {
        fmt.Println(p, q)
        p += 2
        q -= 2
    }

}

func TestBTree_AddNerd(t *testing.T) {

    sum = 0
    tree := NewBTree(2)

    pMap := make(map[int]bool)
    qMap := make(map[int]bool)

    for {

        var p, q int

        for  {
            p = rand.Int() % 100001
            if _, ok := pMap[p]; !ok {
                break
            }
        }

        for {
            q = rand.Int() % 100001
            if _, ok := qMap[q]; !ok {
                break
            }
        }



        pMap[p] = true
        qMap[q] = true

        fmt.Println(p, q, len(pMap), len(qMap))

        tree.AddNerd(p, q)
    }


}