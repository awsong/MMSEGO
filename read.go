// +build ignore

package main

import (
    "fmt"
    "os"
    "encoding/gob"
    )

type Thesaurus struct{
    Offset int;
    Body string;
}
type DictItem struct{
    Freq int;
    ThesaurusCount int;
    TT [8]Thesaurus
}

func main() {
    dict := make(map[string]DictItem);
    file, _ := os.Open("uni.lib");
    defer file.Close();

    dec := gob.NewDecoder(file);
    dec.Decode(&dict);
    fmt.Printf("å: %v\n", dict["çå±çäºº"].TT[0]);
//    for key, value := range dict {
//	fmt.Printf("key %s, value %d\n", key, value)
//    }
}
