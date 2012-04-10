// +build ignore

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
    var s = new(mmsego.Segmenter)
    s.Init("/public/development/go/src/mmsego/darts.lib")
    f, err := os.Create("/tmp/gprof")
    if err != nil {
	log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    t := time.Now()
    offset := 0

    unifile, _ := os.Open("/tmp/a.txt")
    uniLineReader := bufio.NewReaderSize(unifile, 4000)
    line, bufErr := uniLineReader.ReadString('\n')
    for nil == bufErr {
	takeWord := func(off int, length int){ fmt.Printf("%s ", string(line[off-offset:off-offset+length])) }
	s.Mmseg(line[:], offset, takeWord, nil, false)
	offset += len(line)
	line, bufErr = uniLineReader.ReadString('\n')
    }
    takeWord := func(off int, length int){ fmt.Printf("%s ", string(line[off-offset:off-offset+length])) }
    s.Mmseg(line, offset, takeWord, nil, true)

    fmt.Printf("Duration: %v\n", time.Since(t))
}
