package main

import (
    "bufio"
    "os"
    "strconv"
    "log"
    "math"
)

type Bitmap16  uint16

type Course struct {
    PrevCourses Bitmap16
    //NextCourse Bitmap16
    //Semester Bitmap16
}

type Semester struct {
    Courses Bitmap16
}

type Done struct {
    SemesterCount int
    Courses Bitmap16
}

var bitmap16Mask = []Bitmap16{
    1<<0, 1<<1, 1<<2, 1<<3,
    1<<4, 1<<5, 1<<6, 1<<7,
    1<<8, 1<<9, 1<<10, 1<<11,
    1<<12, 1<<13, 1<<14, 1<<15,
}

var bitIndexListCache = make([][]int, 1<<16)

var factorial = []int{
    0, //0
    1, //1
    2, //2
    6, //3
    24, //4
    120, //5
    720, //6
    5040, //7
    40320, //8
    362880, //9
    3628800, //10
    39916800, //11
    479001600, //12
    6227020800, //13
    87178291200, //14
    1307674368000, //15
}

var StdoutWriter *bufio.Writer

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)

    initBitIndexList(0, Bitmap16(0))

    //for i:=0; i<100; i++ {
    //    log.Println(i, bitIndexListCache[i])
    //}

}

func initBitIndexList(bitIndex int, curValue Bitmap16) {

    if bitIndex == 16 {
        return
    }

    // set 0
    initBitIndexList(bitIndex + 1, curValue)

    // set 1
    newValue := curValue.EnableAt(bitIndex)
    bitIndexListCache[newValue] = make([]int, 0, len(bitIndexListCache[curValue]) + 1)
    bitIndexListCache[newValue] = append(bitIndexListCache[newValue], bitIndexListCache[curValue]...)
    bitIndexListCache[newValue] = append(bitIndexListCache[newValue], bitIndex)

    initBitIndexList(bitIndex + 1, newValue)


}

func main() {

    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    caseNum := scanInt(scanner) // max 50

    for ; caseNum > 0; caseNum -- {
        resultM := scanCase(scanner)

        if resultM < 0 {
            StdoutWriter.WriteString("IMPOSSIBLE\n")
        } else {
            StdoutWriter.WriteString(strconv.Itoa(resultM))
            StdoutWriter.WriteByte('\n')
        }

    }

}

func scanCase(scanner *bufio.Scanner) int {

    //log.Println("----------------")

    N := scanInt(scanner) // 전공 과목 수. 12 이하
    K := scanInt(scanner) // 들어야할 과목 수. N 이하(최대 12)
    M := scanInt(scanner) // 학기 수. 10 이하
    L := scanInt(scanner) // 한 학기 최대 과목 수. 10 이하

    courseList := make([]Course, N)
    for i:=0; i<N; i++ {
        Ri := scanInt(scanner)
        for j:=0; j< Ri; j++ {
            courseId := scanInt(scanner)
            courseList[i].PrevCourses = courseList[i].PrevCourses.EnableAt(courseId)
        }

    }

    semesterList := make([]Semester, M)
    for i:=0; i<M; i++ {
        Ci := scanInt(scanner) // 10 이하
        for j:=0; j<Ci; j++ {
            courseId := scanInt(scanner)
            //courseList[courseId].Semester = courseList[courseId].Semester.EnableAt(i)
            semesterList[i].Courses = semesterList[i].Courses.EnableAt(courseId)
        }
    }


    if K == 0 {
        return 0
    }

    doneMap := make(map[Bitmap16]int) //key: complete courses, value: minimum semester
    doneMap[0] = 0

    minSem := math.MaxInt32

    for i:=0; i<M; i++ {

        newDoneMap := enrollSemester(N, K, M, L, courseList, semesterList, i, doneMap)

        //log.Println("------------")
        //log.Println(doneMap)
        //log.Println(newDoneList)

        for courses, count := range newDoneMap {

            if len(bitIndexListCache[courses]) < K {

                if exMinSem, ok := doneMap[courses]; ok {
                    doneMap[courses] = min(exMinSem, count)
                } else {
                    doneMap[courses] = count
                }

            } else {
                minSem = min(minSem, count)
            }

        }

    }

    if minSem == math.MaxInt32 {
        return -1
    }

    return minSem


}

func enrollSemester(N, K, M, L int, courseList []Course, semesterList []Semester, semesterIndex int, doneMap map[Bitmap16]int) (map[Bitmap16]int) {

    resultDone := make(map[Bitmap16]int)

    for doneCourses, semCount := range doneMap {

        already := doneCourses & semesterList[semesterIndex].Courses
        enrollable := semesterList[semesterIndex].Courses ^ already

        enrollableList := enrollable.GetEnabledIndexes()
        enrollableList2 := make([]int, 0, len(enrollableList))

        for i:=0; i <len(enrollableList); i++ {

            courseId := enrollableList[i]
            course := courseList[courseId]

            if course.PrevCourses.And(doneCourses) == course.PrevCourses { // 현재 조합으로는 수강 불가능한 과목
                enrollableList2 = append(enrollableList2, courseId)
            }

        }

        thisSemCombiList := make([]Bitmap16, 0)
        GenBitmapCombi(enrollableList2, 0, L, Bitmap16(0), &thisSemCombiList)

        for _, thisSemCourses := range thisSemCombiList {

            newDoneCourses := thisSemCourses.Or(doneCourses)
            if exMinSem, ok := resultDone[newDoneCourses]; ok {
                resultDone[newDoneCourses] = min(exMinSem, semCount + 1)
            } else {
                resultDone[newDoneCourses] = semCount + 1
            }

        }

    }

    return resultDone

}

func GenBitmapCombi(bitmapIndexList []int, startIndex int, remainedCount int, curCombi Bitmap16, result *[]Bitmap16) {

    //log.Println(startIndex, remainedCount, curCombi)

    if remainedCount > len(bitmapIndexList) {
        remainedCount = len(bitmapIndexList)
    }

    if remainedCount == 0 || startIndex == len(bitmapIndexList){
        *result = append(*result, curCombi)
        return
    }

    for i:=startIndex; i<=len(bitmapIndexList)-remainedCount; i++ {
        GenBitmapCombi(bitmapIndexList, i + 1, remainedCount - 1, curCombi.EnableAt(bitmapIndexList[i]), result)
    }

}

func (bitmap Bitmap16) IsEnabledAt(index int) bool {
    return bitmap & bitmap16Mask[index] > 0
}

func (bitmap Bitmap16) EnableAt(index int) Bitmap16 {
    return bitmap | bitmap16Mask[index]
}

func (bitmap Bitmap16) DisableAt(index int) Bitmap16 {
    if bitmap.IsEnabledAt(index){
        return bitmap - bitmap16Mask[index]
    }
    return bitmap
}

func (bitmap Bitmap16) GetEnabledIndexes() []int {

    return bitIndexListCache[bitmap]

}

func (b1 Bitmap16) And(b2 Bitmap16) Bitmap16 {
    return b1 & b2
}

func (b1 Bitmap16) Or(b2 Bitmap16) Bitmap16 {
    return b1 | b2
}

func scanString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
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