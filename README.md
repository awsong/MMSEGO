MMSEGO
=====
This is a GO implementation of [MMSEG](http://technology.chtsai.org/mmseg/) which a Chinese word splitting algorithm.

TO DO list
----------
* Documentation/comments
* Benchmark

Usage
---------
#Input Dictionary Format
```sh
Key\tFreq
```
Each key occupies one line. The file should be utf-8 encoded, please refer to [go-darts](https://github.com/awsong/go-darts)

#Code example
```go
package main

import (
    "fmt"
    "time"
    "os"
    "mmsego"
    "bufio"
    "log"
    )

func main() {
    var s = new(mmsego.Segmenter)
    s.Init("darts.lib")
    if err != nil {
	log.Fatal(err)
    }

    t := time.Now()
    offset := 0

    unifile, _ := os.Open("/tmp/a.txt")
    uniLineReader := bufio.NewReaderSize(unifile, 4000)
    line, bufErr := uniLineReader.ReadString('\n')
    for nil == bufErr {
	//takeWord := func(off int, length int){ fmt.Printf("%s ", string(line[off-offset:off-offset+length])) }
	takeWord := func(off, length int){ }
	s.Mmseg(line[:], offset, takeWord, nil, false)
	offset += len(line)
	line, bufErr = uniLineReader.ReadString('\n')
    }
    takeWord := func(off int, length int){ fmt.Printf("%s ", string(line[off-offset:off-offset+length])) }
    s.Mmseg(line, offset, takeWord, nil, true)

    fmt.Printf("Duration: %v\n", time.Since(t))
}
```
LICENSE
-----------
Apache License 2.0
