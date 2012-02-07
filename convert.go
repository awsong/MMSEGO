package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "encoding/gob"
    )

type Thesaurus struct{
    Offset int;
    Body string;
}
type DictItem struct{
    Freq int;
    ThesaurusCount int;
    TT [8]Thesaurus;
}

func main(){
    dict := make(map[string]*DictItem);
    unifile, _ := os.Open("unigram.txt");
    defer unifile.Close();
    thesaurusfile, _ := os.Open("thesaurus.txt");
    defer thesaurusfile.Close();
    ofile, _ := os.Create("uni.lib");
    defer ofile.Close();

    uniLineReader, _ := bufio.NewReaderSize(unifile, 400);
    line, _, bufErr := uniLineReader.ReadLine();
    for nil == bufErr {
	rst := strings.Split(string(line)[:], "\t");
	key := rst[0];
	freq, _ := strconv.Atoi(rst[1]);
	//fmt.Printf("This is a line: %v:%d, end\n", rst[0], freq);
	dict[key] = &DictItem{freq, 0, [8]Thesaurus{}};
	line, _, bufErr = uniLineReader.ReadLine();
    }

    thesaurusLineReader, _ := bufio.NewReaderSize(thesaurusfile, 400);
    line, _, bufErr = thesaurusLineReader.ReadLine();
    for nil == bufErr {
	key := string(line);
	line, _, bufErr = thesaurusLineReader.ReadLine();
	ttts := strings.Split(string(line)[:], "\t");
//	fmt.Printf("key=%v\n", key);
	dict[key].ThesaurusCount = len(ttts);
	for i:=0; i<len(ttts); i++{
	    dict[key].TT[i] = Thesaurus{strings.Index(key, ttts[i]), ttts[i]}
//	    fmt.Printf("\t%v end\n", ttts[i]);
	}
	line, _, bufErr = thesaurusLineReader.ReadLine();
    }
//    for key, value := range dict {
//	fmt.Printf("key %s, value %d\n", key, value)
//    }
    fmt.Printf("dict length: %v\n", len(dict));
    enc := gob.NewEncoder(ofile);
    enc.Encode(dict);
}
