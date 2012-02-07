package mmsego

import (
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

type Dict map[string]*DictItem

func LoadDictionary(filename string) Dict{
    dict := make(Dict);
    file, _ := os.Open(filename);
    defer file.Close();

    dec := gob.NewDecoder(file);
    dec.Decode(&dict);
    return dict;
}
