package main

import (
    "bufio"
    "os"
    "log"
    "strconv"
    "fmt"
    "sort"
)

type Piece struct {

    Type      int
    pattern   byte // if isWildcard == false
    accMinLen int

}

const (
    PIECE_TYPE_ALPHANUMERIC = 1
    PIECE_TYPE_QUESTION = 2
    PIECE_TYPE_ASTERISK = 3
)

type Wildcard struct {
    original string
    pieces []Piece

}


func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)

    caseNum := scanLineAsInt(scanner)

    for ; caseNum > 0 ; caseNum -- {

        wildcardStr := fmt.Sprintf("^%s$", scanLineAsString(scanner))

        wildcard := NewWildcard(wildcardStr)
        //log.Println(wildcard.pieces)

        fileNum := scanLineAsInt(scanner)

        fileList := make([]string, fileNum)
        matchFileList := make([]string, 0, fileNum)
        for i := 0; i<fileNum ; i++ {
            filename := scanLineAsString(scanner)
            fileList[i] = filename

            //log.Println("---------------------")
            //log.Println("pattern:", wildcardStr)
            if wildcard.match(filename) {
                //log.Println(filename, ": match")
                matchFileList = append(matchFileList, filename)
            }
        }

        sort.Strings(matchFileList)
        for _, str := range matchFileList {
            fmt.Println(str)
        }

    }
}


func NewWildcard(str string) *Wildcard {
    result := &Wildcard{}
    result.original = str
    result.pieces = make([]Piece, 0, len(str))

    for i:=0; i<len(str); i++ {
        if str[i] != '*' && str[i] != '?' {

            piece := Piece{
                Type:PIECE_TYPE_ALPHANUMERIC,
                pattern:str[i],
                accMinLen: 1,
            }
            result.pieces = append(result.pieces, piece)

        } else {

            includeAsterisk := false
            j := i
            for ; j < len(str); j++ {
                if str[j] == '?' {
                    piece := Piece{
                        Type:PIECE_TYPE_QUESTION,
                        pattern: '?',
                        accMinLen:1,
                    }

                    result.pieces = append(result.pieces, piece)
                } else if str[j] == '*' {
                    includeAsterisk = true
                } else {
                    break
                }
            }

            if includeAsterisk {
                piece := Piece{
                    Type:PIECE_TYPE_ASTERISK,
                    pattern: '*',
                    accMinLen:0,
                }
                result.pieces = append(result.pieces, piece)
            }
            i = j - 1
        }

    }

    for i := len(result.pieces) - 2; i>=0; i-- {
        result.pieces[i].accMinLen += result.pieces[i + 1].accMinLen
    }

    return result

}

func (w *Wildcard) match(input string) bool {

    return w.matchInternal(fmt.Sprintf("^%s$", input), 0, 0, false)

}

func (w *Wildcard) matchInternal(input string, inputIndex int, pieceIndex int, exPieceAsterisk bool) bool {


    //log.Printf("remained input: [%s], pieces: %v, inputIndex: %d, pieceIndex: %d\n",
    //    input[inputIndex:],
    //    w.pieces[pieceIndex:],
    //    inputIndex, pieceIndex)

    if pieceIndex == len(w.pieces) {
        return inputIndex == len(input)

    }

    piece := w.pieces[pieceIndex]

    if inputIndex + piece.accMinLen > len(input) {
        return false
    }

    // 현재 piece가 wildcard라면...
    if piece.Type == PIECE_TYPE_ASTERISK {
        result := w.matchInternal(input, inputIndex, pieceIndex + 1, true)
        return result
    }

    // wildcard가 아닐 경우...
    // question mark or alphanumeric

    if exPieceAsterisk {
        for i:=inputIndex; i<=len(input) - piece.accMinLen; i++ {

            if piece.match(input[i]) &&
                w.matchInternal(input, i + 1, pieceIndex + 1, false) {
                return true
            }
        }
        return false

    } else {
        isMatch := piece.match(input[inputIndex])
        if isMatch {
            result := w.matchInternal(input, inputIndex + 1, pieceIndex + 1, false)
            return result
        } else {
            return false
        }
    }

}

func (p Piece) match(b byte) bool {

    if p.Type == PIECE_TYPE_QUESTION {
        return true
    } else {
        return p.pattern == b
    }

}


func scanLineAsString(scanner *bufio.Scanner) string {
    if ! scanner.Scan() {
        log.Println(scanner.Err())
        panic(scanner.Err())
    }
    return scanner.Text()
}

func scanLineAsInt(scanner *bufio.Scanner) int {
    result, _ := strconv.Atoi(scanLineAsString(scanner))
    return result
}
