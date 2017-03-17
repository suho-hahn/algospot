package main

import (
    "testing"
    "log"
)


func TestGetCombCount(t *testing.T) {

    join := JoinStr{
        strList:[]string{"1", "2", "3", "4"},
    }

    log.Println(join.GetSubKeys(0xf))

}