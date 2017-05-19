package main

import (
    "testing"
    "log"
)

func TestGenBitmapCombi(t *testing.T) {

    result := []Bitmap16{}
    GenBitmapCombi([]int{0, 1, 2}, 0, 2, 0, &result)
    log.Println(result)

    result = []Bitmap16{}
    GenBitmapCombi([]int{0, 1, 4}, 0, 4, 0, &result)
    log.Println(result)

}
