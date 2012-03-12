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
    "time"
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
	line, _, bufErr = uniLineReader.ReadLine();
    }
    sort.Strings(strings)
    for _, s := range strings{
	rst := bytes.Split([]byte(s), []byte("\t"));
	freq, _ := strconv.Atoi(string(rst[1]));
	keys = append(keys, rst[0])
	values = append(values, freq)
    }

    fmt.Printf("input dict length: %v %v\n", len(keys), len(values));
    round := len(keys)
    d := darts.Build(keys[:round], values[:round])
    fmt.Printf("build out length %v\n", len(d))
    t := time.Now()
    for i := 0; i < round; i++ {
	if true != d.ExactMatchSearch(keys[i],0){
	    fmt.Println("wrong", string(keys[i]), i)
	}
    }
    fmt.Println(time.Since(t))
    enc := gob.NewEncoder(ofile);
    enc.Encode(d);
}
