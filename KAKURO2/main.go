package main

import (
    "bufio"
    "os"
    "log"
    "strconv"
    "fmt"
)

const (
    SlotTypeBlack = 0
    SlotTypeWhite = 1
    DirectionHorizontal = 0
    DirectionVertical = 1

    HintValidFalse = 0
    HintValidTrue = 1
    HintValidUnknown = -1
)



type Board struct {
    Size  int
    Slots [][]Slot
    Hints []Hint
}

type Slot struct {

    Type  int8
    Value int

    Row        int
    Col        int

    HHint *Hint
    VHint *Hint

    CandidateBitmap Bitmap16

}

type Hint struct {

    Row        int
    Col        int
    Direction  int
    Sum        int

    Index        int
    Slots        []*Slot

    CombiList    []Bitmap16
    Used         Bitmap16

}

//type HintSlice []Hint
//type SlotSlice []*Slot
type Bitmap16  uint16

var StdoutWriter *bufio.Writer
var combiCache [][][]Bitmap16 // sum(1~45) / slot #(1~9) / combinations(bitmap)
var bitmap16Mask = []Bitmap16{
    1<<0, 1<<1, 1<<2, 1<<3,
    1<<4, 1<<5, 1<<6, 1<<7,
    1<<8, 1<<9, 1<<10, 1<<11,
    1<<12, 1<<13, 1<<14, 1<<15,
}
var factorial = []int{
    0,
    1,
    2,
    6,
    24,
    120,
    720,
    5040,
    40320,
    362880,
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)


    combiCache = make([][][]Bitmap16, 46) // 1 ~ 45
    for i:=0; i<=45; i++ {
        combiCache[i] = make([][]Bitmap16, 10)
        combiCache[i][0] = []Bitmap16{} // (0,0), (1,0), ..., (44,0), (45,0)
    }

    for i:=1; i<=9; i++ {
        combiCache[0][i] = []Bitmap16{}
        combiCache[i][1] = []Bitmap16{ 1<<uint(i) }
    }

    factorial[1] = 1
    for i:=2; i<len(factorial); i++ {
        factorial[i] = factorial[i-1] * i
    }

    //for i:=range combiCache{
    //    log.Println(i, combiCache[i])
    //}

}

func main() {
    StdoutWriter = bufio.NewWriter(os.Stdout)
    defer StdoutWriter.Flush()

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)

    caseNum := scanAsInt(scanner) // max 50

    for ; caseNum > 0; caseNum -- {
        scanCase(scanner)
    }

}

func scanCase(scanner *bufio.Scanner) {
    boardSize := scanAsInt(scanner) // max 20

    board := Board{}
    board.Size = boardSize
    board.Slots = make([][]Slot, boardSize)

    // read board info
    for i:=0; i<boardSize; i++ {
        board.Slots[i] = make([]Slot, boardSize)
        for j:=0; j<boardSize; j++ {
            board.Slots[i][j].Type = int8(scanAsInt(scanner))
            board.Slots[i][j].Row = i
            board.Slots[i][j].Col = j
        }

        //log.Println(board.Slots[i])

    }

    hintSize := scanAsInt(scanner)
    board.Hints = make([]Hint, hintSize)

    // read hint info
    for i:=0; i<hintSize; i++ {

        hint := &board.Hints[i]

        hint.Row = scanAsInt(scanner) - 1
        hint.Col = scanAsInt(scanner) - 1
        hint.Direction = scanAsInt(scanner)
        hint.Sum = scanAsInt(scanner)
        hint.Index = i
        hint.Slots = make([]*Slot, 0, 9)

        // set related slots for each hint
        if hint.Direction == DirectionVertical {
            for j:=hint.Row + 1; j<boardSize && board.Slots[j][hint.Col].Type == SlotTypeWhite ; j++ {
                hint.Slots = append(hint.Slots, &board.Slots[j][hint.Col])
                board.Slots[j][hint.Col].VHint = &board.Hints[i]
            }
        } else { //horizontal
            for j:=hint.Col + 1; j<boardSize && board.Slots[hint.Row][j].Type == SlotTypeWhite ; j++ {
                //log.Println(j)
                hint.Slots = append(hint.Slots, &board.Slots[hint.Row][j])
                board.Slots[hint.Row][j].HHint = &board.Hints[i]
            }
        }

        hint.CombiList = getCombiList(hint.Sum, len(hint.Slots))

        //log.Println(hint.Row, hint.Col, hint.CombiList)
    }

    board.Solve()

    //for i:=0; i<hintSize; i++ {
    //    hint := &board.Hints[i]
    //    log.Println(hint.Row, hint.Col, hint.Direction, hint.CombiList)
    //}

}

func (board *Board) Solve() {

    board.CalcSlotCandidate()

    remainedSlots := make([]*Slot, 0, board.Size * board.Size)
    for i:=0; i<board.Size; i++ {
        for j:=0; j<board.Size; j++ {
            slot := &board.Slots[i][j]
            if slot.Type == SlotTypeBlack {
                continue
            } else if slot.Value > 0 {
                continue
            }
            remainedSlots = append(remainedSlots, slot)
            //log.Println(slot)
        }
    }

    //log.Println("--------------")
    board.FillSlots(remainedSlots)

    //board.LogSlots()
    //time.Sleep(1 * time.Second)
    board.PrintSlots()

}


func (board *Board)CalcSlotCandidate() {

    slots := make([]*Slot, 0, board.Size * board.Size * 2)

    for i:=0; i<board.Size; i++ {
        for j:=0; j<board.Size; j++ {
            if board.Slots[i][j].Type == SlotTypeBlack {
                continue
            }
            slots = append(slots, &board.Slots[i][j])
        }
    }

    for len(slots) > 0 {


        //slot := slots[0]
        //slots = append(slots[:0], slots[1:]...)
        slot := slots[len(slots) - 1]
        slots = slots[:len(slots) - 1]


        //log.Println("----------------")
        //log.Println(slot)

        if slot.Value > 0 { // value is fixed
            continue
        }

        vHint := slot.VHint
        hHint := slot.HHint

        vUsed := vHint.Used
        hUsed := hHint.Used

        newCandidateBitmap := Bitmap16(0)

        vCombis := Bitmap16(0)
        for _, vCombi := range vHint.CombiList {
            vCombis = vCombis.Or(vCombi)
        }
        vCombis = (vCombis ^ vUsed) & vCombis

        hCombis := Bitmap16(0)
        for _, hCombi := range hHint.CombiList {
            hCombis = hCombis.Or(hCombi)
        }
        hCombis = (hCombis ^ hUsed) & hCombis

        newCandidateBitmap = hCombis.And(vCombis)

        if len(vHint.CombiList) == 1 {
            newCandidateBitmap = ReduceCandidate(vHint, slot, newCandidateBitmap)
        }

        if len(hHint.CombiList) == 1 {
            newCandidateBitmap = ReduceCandidate(hHint, slot, newCandidateBitmap)
        }

        if slot.CandidateBitmap == newCandidateBitmap {
            continue
        }

        slot.CandidateBitmap = newCandidateBitmap
        candidates := slot.CandidateBitmap.GetEnabledIndexes()
        //valueFixed := false
        if len(candidates) == 1 {
            slot.Value = candidates[0]
            //valueFixed = true
            vHint.Used = vHint.Used.Or(bitmap16Mask[slot.Value])
            hHint.Used = hHint.Used.Or(bitmap16Mask[slot.Value])
        }

        //OptimizeHint(vHint)
        //OptimizeHint(hHint)

        slots = append(slots, vHint.Slots...)
        slots = append(slots, hHint.Slots...)

        //if OptimizeHint(vHint, slot) || valueFixed {
        //    slots = append(slots, vHint.Slots...)
        //}

        //if OptimizeHint(hHint, slot) || valueFixed {
        //    slots = append(slots, hHint.Slots...)
        //}

    }

}

func (board *Board) FillSlots(slots []*Slot) bool {
    if len(slots) == 0 {
        return true
    }

    slot := slots[0]

    //log.Println(slot)

    vHint := slot.VHint
    hHint := slot.HHint

    candidates := slot.CandidateBitmap.GetEnabledIndexes()

    for _, cand := range candidates {

        if vHint.Used.IsEnabledAt(cand) {
            continue
        }
        if hHint.Used.IsEnabledAt(cand) {
            continue
        }

        slot.Value = cand
        vHint.Used = vHint.Used.EnableAt(cand)
        hHint.Used = hHint.Used.EnableAt(cand)

        if  vHint.IsValid() != HintValidFalse &&
            hHint.IsValid() != HintValidFalse &&
            board.FillSlots(slots[1:]) {
            return true
        }
        slot.Value = 0
        vHint.Used = vHint.Used.DisableAt(cand)
        hHint.Used = hHint.Used.DisableAt(cand)


    }

    return false

}

func (hint *Hint) IsValid() int8 {

    sum := 0

    for _, slot := range hint.Slots {
        if slot.Value == 0 {
            return HintValidUnknown
        }

        sum += slot.Value

    }

    if sum == hint.Sum {
        return HintValidTrue
    } else{
        return HintValidFalse
    }

}

func OptimizeHint(hint *Hint) bool {

    allCandMap := Bitmap16(0)
    for _, slot := range hint.Slots {
        if slot.CandidateBitmap == 0 {
            return false
        }
        allCandMap = allCandMap.Or(slot.CandidateBitmap)
    }

    hintCombiChanged := false
    for i:=0; i<len(hint.CombiList); i++ {

        combi := hint.CombiList[i]

        if combi.And(allCandMap) == 0 {

            //log.Println("unused value")
            hint.CombiList = append(hint.CombiList[:i], hint.CombiList[i+1:]...)
            i --
            hintCombiChanged = true
        }

    }

    //log.Println(hint.CombiList)
    return hintCombiChanged

}

func OptimizeHint2(hint *Hint) bool {

    candidateBitmap := Bitmap16(0)
    for _, slot := range hint.Slots {
        if slot.CandidateBitmap == 0 {
            return false
        }
        candidateBitmap = candidateBitmap.Or(slot.CandidateBitmap)
    }

    hintCombiChanged := false
    for i:=0; i<len(hint.CombiList); i++ {

        combi := hint.CombiList[i]

        if combi.And(candidateBitmap) == 0 {

            //log.Println("unused value")
            hint.CombiList = append(hint.CombiList[:i], hint.CombiList[i+1:]...)
            i --
            hintCombiChanged = true
        }

    }

    //log.Println(hint.CombiList)
    return hintCombiChanged

}

func ReduceCandidate(hint *Hint, slot *Slot, newBitmap Bitmap16) Bitmap16 {

    otherBitmap := Bitmap16(0)

    for _, otherSlot := range hint.Slots {
        if otherSlot == slot {
            continue
        } else if otherSlot.CandidateBitmap == 0 {
            return newBitmap
        }

        otherBitmap = otherBitmap | otherSlot.CandidateBitmap
    }

    xor := otherBitmap ^ newBitmap
    result := xor & newBitmap
    if result == 0 {
        return newBitmap
    } else {
        return result
    }

}


func (bitmap Bitmap16) IsEnabledAt(index int) bool {
    return bitmap & bitmap16Mask[index] > 0
}

func (bitmap Bitmap16) EnableAt(index int) Bitmap16 {
    return bitmap | bitmap16Mask[index]
}

func (bitmap Bitmap16) DisableAt(index int) Bitmap16 {
    return bitmap - bitmap16Mask[index]
}

func (bitmap Bitmap16) GetEnabledIndexes() []int {

    result := make([]int, 0, 9)
    for i:=1; i<=9; i++ {
        mask := bitmap16Mask[i]
        if mask > bitmap {
            break
        }

        if bitmap.And(mask) == mask {
            result = append(result, i)
        }
    }

    return result

}

func (b1 Bitmap16) And(b2 Bitmap16) Bitmap16 {
    return b1 & b2
}

func (b1 Bitmap16) Or(b2 Bitmap16) Bitmap16 {
    return b1 | b2
}

func (board *Board) LogSlots() {

    for _, row := range board.Slots {
        line := ""
        for _, slot := range row {
            line = fmt.Sprintf("%s%d:%d(%v)\t\t", line, slot.Type, slot.Value, slot.CandidateBitmap.GetEnabledIndexes())
            //if slot.Type == SlotTypeWhite && len(slot.CandidateBitmap.GetEnabledIndexes()) == 0 {
            //    panic("")
            //}
        }
        log.Println(line)
    }

}


func (board *Board) PrintSlots() {

    for _, row := range board.Slots {
        for i, slot := range row {
            StdoutWriter.WriteString(strconv.Itoa(slot.Value))
            if i == board.Size - 1 {
                StdoutWriter.WriteByte('\n')
            } else{
                StdoutWriter.WriteByte(' ')
            }
        }
    }

}

func getCombiList(sum, slotCount int) []Bitmap16{


    //log.Println(sum, slotCount)

    if combiCache[sum][slotCount] != nil {
        //log.Println("cache hit")
        return combiCache[sum][slotCount]
    }

    maxValue := 2 * sum / slotCount - 1
    minValue := 2 * sum / slotCount - 9

    if minValue < 1 {
        minValue = 1
    }

    if maxValue > 9 {
        maxValue = 9
    }


    minValue = max(minValue, 1)
    maxValue = min(maxValue, 9)

    //log.Println(">>", sum, slotCount, minValue, maxValue)

    resultMap := make(map[Bitmap16]bool)

    for i := minValue; i <= maxValue; i ++ {

        subList := getCombiList(sum - i, slotCount - 1)
        for _, subBitmap := range subList {
            if subBitmap.IsEnabledAt(i) {
                continue
            }
            resultKey := subBitmap.EnableAt(i)

            //log.Println(resultMap[resultKey])

            resultMap[resultKey] = true
        }
    }

    resultList := make([]Bitmap16, 0, len(resultMap))
    for key := range resultMap {
        resultList = append(resultList, key)
    }
    //sort.Sort([]uint16(resultList))

    combiCache[sum][slotCount] = resultList

    return resultList

}

func scanAsString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanAsInt(scanner *bufio.Scanner) int {
    result, _ := strconv.Atoi(scanAsString(scanner))
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