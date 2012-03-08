// +build ignore

package main

import (
    "fmt"
    "darts"
    "os"
    "bufio"
    "bytes"
    "sort"
    "strconv"
    "encoding/gob"
    )

func main(){
    strings := make([]string, 0, 4096)
    keys := make([][]byte, 0, 4096)
    values := make([]int, 0, 4096)
    unifile, _ := os.Open("darts.txt");
    defer unifile.Close();
    ofile, _ := os.Create("darts.lib");
    defer ofile.Close();

    uniLineReader := bufio.NewReaderSize(unifile, 400);
    line, _, bufErr := uniLineReader.ReadLine();
    for nil == bufErr {
	strings = append(strings, string(line))
	//fmt.Printf("This is a line: %v:%d, end\n", rst[0], freq);
	line, _, bufErr = uniLineReader.ReadLine();
    }
    sort.Strings(strings)
    for _, s := range strings{
	rst := bytes.Split([]byte(s), []byte("\t"));
	freq, _ := strconv.Atoi(string(rst[1]));
	keys = append(keys, rst[0])
	values = append(values, freq)
    }

//    for key, value := range dict {
//	fmt.Printf("key %s, value %d\n", key, value)
//    }
    fmt.Printf("dict length: %v %v\n", len(keys), len(values));
    d := darts.Build(keys, values)
    enc := gob.NewEncoder(ofile);
    enc.Encode(d);
}
