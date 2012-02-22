// +build ignore

package main

import (
"fmt"
//"strings"
)

func main(){
    a := "ä¸­å½äººæ°é¢"
    for i, v := range a{
	fmt.Printf("%v, %v\n", i, v);
    }
    var mm = make([]int , 200)
    mm[100] = -1
    for i, v := range mm[10:]{
	if v == -1{
	    fmt.Println((i+10))
	}
    }
}
