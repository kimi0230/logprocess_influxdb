package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc     chan string // read channel
	wc     chan string // write channel
	read   Reader
	writer Writer
}

type ReadFromFile struct {
	path string // 讀取文件的路徑
}

func (r *ReadFromFile) Read(rc chan string) {
	// 讀取模塊
	line := "message"
	rc <- line
}

type WriteToInfluxDB struct {
	influxDBsn string // influx data source
}

func (l *WriteToInfluxDB) Write(wc chan string) {
	// 寫入模塊
	fmt.Println(<-wc)
}

func (l *LogProcess) Process() {
	// 解析模塊
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func main() {
	r := &ReadFromFile{
		path: "/tmp/access.log",
	}
	w := &WriteToInfluxDB{
		influxDBsn: "username&password...",
	}

	lp := &LogProcess{
		rc:     make(chan string),
		wc:     make(chan string),
		read:   r,
		writer: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)
	time.Sleep(time.Second * 1)
}
