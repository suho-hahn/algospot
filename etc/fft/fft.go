package main

import (
    //"fmt"
    //"github.com/mjibson/go-dsp/fft"
)
import (
    //"math"
    "fmt"
    "math"
    "log"
)

func main(){

    log.SetFlags(log.LstdFlags | log.Lshortfile)

    a := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
    b := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

    fmt.Println(fftMul(a, b))

}


func dft(p []complex128, w complex128) {

    n := len(p)
    if n == 1 {
        return
    }

    // divide
    even := make([]complex128, n / 2)
    odd := make([]complex128, n / 2)

    for i:=0; i < n/2 ; i++ {
        even[i] = p[i * 2]
        odd[i] = p[i * 2 + 1]
    }

    // conquer
    dft(even, w * w)
    dft(odd, w * w)

    //merge
    w_power := complex(float64(1), float64(0)) // 1 + 0i
    for i:=0; i<n/2; i++ {
        offset := w_power * odd[i]
        p[i] = even[i] + offset
        p[i + n/2] = even[i] - offset
        w_power *= w
    }

}

func roundUp(n int) int{
    power2 := 1;
    for power2 < n {
        power2 *= 2
    }
    return power2
}

func fftMul(a, b []int) []int {

    const PI = math.Pi

    // last *2 is needed because C can have twice the degree from original polys
    n := roundUp(max(len(a), len(b))) * 2

    // Complex representations of a and b
    complexA, complexB := make([]complex128, n), make([]complex128, n)
    for i := range a{
        complexA[i] = complex(float64(a[i]), float64(0))
    }
    for i := range b{
        complexB[i] = complex(float64(b[i]), float64(0))
    }

    // FFT
    w := complex(math.Cos(2 * PI / float64(n)), math.Sin(2 * PI / float64(n)))
    log.Println(w)
    log.Println(complexA)
    log.Println(complexB)
    dft(complexA, w)
    dft(complexB, w)
    log.Println(complexA)
    log.Println(complexB)

    // pointwise multiplication
    for i:=0; i<n; i++ {
        complexA[i] *= complexB[i]
    }

    // Inverse FFT
    dft(complexA, complex(math.Cos(2 * PI / float64(n)), math.Sin(-2 * PI / float64(n))))
    for i:=0; i<n; i++ {
        complexA[i] = complex(real(complexA[i]) / float64(n), imag(complexA[i]) / float64(n))
    }

    log.Println(complexA)

    // Trim zero paddings
    for len(complexA) > 0 &&
        math.Abs(real(complexA[len(complexA) - 1])) < float64(0.0000001) {
        complexA = complexA[:len(complexA) - 1]
    }

    result := make([]int, len(complexA))
    for i:=0; i<len(complexA); i++ {
        result[i] = int(math.Floor(real(complexA[i]) + 0.5))
    }

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