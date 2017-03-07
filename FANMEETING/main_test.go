package main

import (
    "testing"
    _ "net/http/pprof"
    "net/http"
)

func Test_fft_mul(t *testing.T) {

    go func(){
        http.ListenAndServe(":6060", nil)
    }()

    const cap = 200000

    a, b := make([]int, cap), make([]int, cap)

    for i:=0; i<cap; i++{
        a[i], b[i] = 1, 1
    }


    for i:=0; ; i++{
        hugs(a, b)

    }



}
