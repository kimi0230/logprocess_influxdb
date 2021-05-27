package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc     chan []byte // read channel
	wc     chan string // write channel
	read   Reader
	writer Writer
}

type ReadFromFile struct {
	path string // 讀取文件的路徑
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 讀取模塊
	// 打開文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error : %s", err.Error()))
	}
	// 從文件末尾開始執行讀取
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			// 讀到末尾
			time.Sleep(500 * time.Microsecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error : %s", err.Error()))
		}
		rc <- line[:len(line)-1]
	}
}

type WriteToInfluxDB struct {
	influxDBsn string // influx data source
}

func (l *WriteToInfluxDB) Write(wc chan string) {
	// 寫入模塊
	for v := range wc {
		fmt.Println(v)
	}
}

func (l *LogProcess) Process() {
	// 解析模塊
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDBsn: "username&password...",
	}

	lp := &LogProcess{
		rc:     make(chan []byte),
		wc:     make(chan string),
		read:   r,
		writer: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)
	time.Sleep(time.Second * 30)
}
