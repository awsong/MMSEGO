package main

import (
    "fmt"
    "time"
    "os"
    "mmsego"
    "bufio"
    "log"
    "runtime/pprof"
    )

func main() {
    f, err := os.Create("/tmp/gprof")
    if err != nil {
	log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    var buf [10000000]int;
    inChan := make(chan string, 50);
    outChan := make(chan [2]int, 50);
    go mmsego.Mmseg(inChan, outChan)

    t := time.Now()
    go func(){
	unifile, _ := os.Open("/tmp/a.txt");
	uniLineReader := bufio.NewReaderSize(unifile, 4000);
	line, _, bufErr := uniLineReader.ReadLine();
	for nil == bufErr {
	    //fmt.Println(string(line[:]))
	    inChan <- string(line)[:]
	    line, _, bufErr = uniLineReader.ReadLine();
	}
	close(inChan)
    }()
    i := 0;
    for m := range outChan {
	buf[i] = m[0]
	buf[i+1] = m[1]
	i += 2
    }
    fmt.Printf("Duration: %v\n", time.Since(t));
    fmt.Println(i)
    /*
    for j := 0; j<i; j++ {
	fmt.Printf("%v ", buf[j])
    }*/
}
