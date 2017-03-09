// https://algospot.com/judge/problem/read/FANMEETING

package main

import (
    "fmt"
    "math"
    "bufio"
    "os"
    "strconv"
    "log"
    "bytes"
)

const (
    MAX_LEN = 200000
)

var complexSliceListA = make([][]complex128, 20)
var complexSliceListB = make([][]complex128, 20)

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    for i := range complexSliceListA {
        complexSliceListA[i] = make([]complex128, 1<<uint(i))
        complexSliceListB[i] = make([]complex128, 1<<uint(i))
    }
}

func main() {

    //var caseNum int = 0
    //fmt.Scanln(&caseNum)
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanBytes)

    var caseNum int
    scanBuffer := make([]byte, 0, 1024*512)
    for {
        if ! scanner.Scan() {
            panic(scanner.Err())
        }
        data := scanner.Bytes()
        if i := bytes.IndexByte(data, '\n'); i >= 0 {
            caseNum, _ = strconv.Atoi(string(append(scanBuffer, data[:i]...)))
            scanBuffer = append(scanBuffer[:0], data[i+1:]...)
            break
        } else {
            scanBuffer = append(scanBuffer, data...)
        }

    }

    members, fans := make([]int, 0, MAX_LEN), make([]int, 0, MAX_LEN)
    for i:=0; i<caseNum; i++{

        members = members[:0]
        fans = fans[:0]

        for {
            if ! scanner.Scan() {
                panic(scanner.Err())
            }

            data := scanner.Bytes()
            scanBuffer = append(scanBuffer, data...)
            newLineIndex := -1

            for bufIdx, b := range scanBuffer {
                if b == 'M' {
                    members = append(members, 1)
                } else if b == 'F' {
                    members = append(members, 0)
                } else {
                    newLineIndex = bufIdx
                    break
                }
            }

            if newLineIndex >= 0 {
                scanBuffer = append(scanBuffer[:0], scanBuffer[newLineIndex + 1:]...)
                break
            } else {
                scanBuffer = scanBuffer[:0]
            }

        }

        for {
            if ! scanner.Scan() {
                panic(scanner.Err())
            }

            data := scanner.Bytes()
            scanBuffer = append(scanBuffer, data...)
            newLineIndex := -1

            for bufIdx, b := range scanBuffer {
                if b == 'M' {
                    fans = append(fans, 1)
                } else if b == 'F' {
                    fans = append(fans, 0)
                } else {
                    newLineIndex = bufIdx
                    break
                }
            }

            if newLineIndex >= 0 {
                scanBuffer = append(scanBuffer[:0], scanBuffer[newLineIndex + 1:]...)
                break
            } else {
                scanBuffer = scanBuffer[:0]
            }

        }

        for j:=0; j<len(fans) / 2; j++ {
            fans[j], fans[len(fans) - j - 1] = fans[len(fans) - j - 1], fans[j]
        }

        if len(members) > len(fans) || len(members) == 0 || len(fans) == 0 {
            fmt.Println(0)
        } else{
            hugs(members, fans)
        }

    }

}

func hugs(members, fans []int) {

    //result := karatsuba(members, fans)
    result := fftMul(members, fans)

    hug := 0
    for i:=len(members) - 1; i<len(fans); i++ {

        if i > len(result) - 1{
            hug ++
        } else if result[i] == 0{
            hug ++
        }
    }

    fmt.Println(hug)
}

/*
 * FFT
 */


func fftMul(a, b []int) []int {

    const PI = math.Pi

    for len(a) > 0 && a[len(a) - 1] == 0 {
        a = a[:len(a) - 1]
    }
    for len(b) > 0 && b[len(b) - 1] == 0 {
        b = b[:len(b) - 1]
    }

    // last *2 is needed because C can have twice the degree from original polys
    nLog2 := properLengthLog2(max(len(a), len(b))) + 1
    n := 1<<uint(nLog2)

    //fmt.Println(len(a), len(b), n)

    // FFT
    complexA := complexSliceListA[nLog2]
    complexB := complexSliceListB[nLog2]

    for i := 0; i<n; i++{
        if i >= len(a){
            complexA[i] = complex(0, 0)
        } else{
            complexA[i] = complex(float64(a[i]), 0)
        }

        if i >= len(b) {
            complexB[i] = complex(0, 0)
        } else{
            complexB[i] = complex(float64(b[i]), 0)
        }

    }

    w := complex(math.Cos(2 * PI / float64(n)), math.Sin(2 * PI / float64(n)))
    dft(complexA, w, nLog2)
    dft(complexB, w, nLog2)

    //fmt.Println(1<<uint(n), len(complexA), len(complexB))

    // pointwise multiplication
    for i:=0; i< n; i++ {
        complexA[i] *= complexB[i]
    }

    // Inverse FFT
    dft(complexA, complex(math.Cos(2 * PI / float64(n)), math.Sin(-2 * PI / float64(n))), nLog2)

    // normalize is not required
    //for i:=0; i<n; i++ {
    //    complexA[i] = complex(real(complexA[i]) / float64(n), imag(complexA[i]) / float64(n))
    //}

    // make result
    result := make([]int, len(complexA))
    for i:=0; i<len(complexA); i++ {
        result[i] = int(real(complexA[i]) + 0.5)
    }

    return result

}

var maxHeapAlloc = uint64(0)

func dft(p []complex128, w complex128, lenLog2 int) {

    n := len(p)
    if n == 1 {

        //var memStat runtime.MemStats
        //runtime.ReadMemStats(&memStat)
        //if maxHeapAlloc < memStat.HeapAlloc {
        //    maxHeapAlloc = memStat.HeapAlloc
        //    fmt.Println("heap ", memStat.HeapAlloc)
        //}

        return
    }

    // divide
    even := complexSliceListA[lenLog2 - 1]
    odd := complexSliceListB[lenLog2 - 1]

    for i:=0; i < n/2 ; i++ {
        even[i] = p[i * 2]
        odd[i] = p[i * 2 + 1]
    }

    // conquer
    dft(even, w * w, lenLog2 - 1)
    dft(odd, w * w, lenLog2 - 1)

    //merge
    w_power := complex(float64(1), float64(0)) // 1 + 0i
    for i:=0; i<n/2; i++ {
        offset := w_power * odd[i]

        p[i] = even[i] + offset
        p[i + n/2] = even[i] - offset

        w_power *= w
    }

}

func properLengthLog2(n int) int{
    power2 := 1;
    shiftCount := 0
    for power2 < n {
        power2 *= 2
        shiftCount ++
    }
    return shiftCount

}

/*
 * KARATSUBA
 */

func karatsuba(a, b []int) []int {

    for len(a) > 0 && a[len(a) - 1] == 0 {
        a = a[:len(a) - 1]
    }
    for len(b) > 0 && b[len(b) - 1] == 0 {
        b = b[:len(b) - 1]
    }

    if len(a) == 0 || len(b) == 0 {
        return []int{}
    }

    if len(a) == 1 && len(b) == 1 {
        value := a[0] * b[0]
        if value == 0 {
            return []int{}
        } else{
            return []int{value}
        }
    }

    var half int
    var a0, a1, b0, b1 []int

    if len(a) < len(b) {
        a, b = b, a
    }

    half = len(a) / 2

    bSplitPos := min(half, len(b))
    a0, a1 = a[:half], a[half:]
    b0, b1 = b[:bSplitPos], b[bSplitPos:]

    //log.Println("[", depth, "]", len(a0), len(a1), len(b0), len(b1))

    chan0 := make(chan []int)
    chan1 := make(chan []int)
    chan2 := make(chan []int)

    go func(){
        chan0 <- karatsuba(a0, b0)
    }()
    go func(){
        chan1 <- karatsuba(addBigInt(a0, a1), addBigInt(b0, b1))
    }()
    go func(){
        chan2 <- karatsuba(a1, b1)
    }()

    z0 := <- chan0
    z1 := <- chan1
    z2 := <- chan2

    z1 = subBigInt(z1, z0)
    z1 = subBigInt(z1, z2)

    result := make([]int, max(len(z2) + 2*half, len(z1) + half))
    for i := range z0 {
        result[i] += z0[i]
    }
    for i := range z1 {
        result[i + half] += z1[i]
    }
    for i := range z2 {
        result[i + half + half] += z2[i]
    }

    //log.Println("subresult:", a, b, result)

    for len(result) > 0 && result[len(result) - 1] == 0{
        result = result[:len(result) - 1]
    }
    return result

}

func addBigInt(a, b []int) []int {
    result := make([]int, max(len(a), len(b)))

    for i := range a {
        result[i] += a[i]
    }

    for i := range b {
        result[i] += b[i]
        if result[i] >= math.MaxInt32 {
            result[i] -= math.MaxInt32
            result[i+1] ++
        }
    }
    //log.Println("add:", a, b, result)
    return result

}

func subBigInt(a, b []int) []int {

    //log.Println(a)
    //log.Println(b)

    for i := range b {
        //log.Println(i)
        a[i] -= b[i]
        if a[i] < 0 {
            a[i] += math.MaxInt32
            a[i+1] --
        }
    }
    return a
}

/*
 * UTILS
 */

func max(a, b int) int{
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int{
    if a > b {
        return b
    }
    return a
}