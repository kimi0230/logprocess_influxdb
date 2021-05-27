package main

import (
	"fmt"
	"strings"
	"time"
)

type LogProcess struct {
	rc         chan string // read channel
	wc         chan string // write channel
	path       string      // 讀取文件的路徑
	influxDBsn string      // influx data source
}

// 使用引用,如果結構體很大時,可以不需要拷貝效能較好. 還可以修改自身定義的參數
func (l *LogProcess) ReadFromFile() {
	// 讀取模塊
	line := "message"
	l.rc <- line
}

func (l *LogProcess) WriteToInfluxDB() {
	// 寫入模塊
	fmt.Println(<-l.wc)

}

func (l *LogProcess) Process() {
	// 解析模塊
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func main() {
	lp := &LogProcess{
		rc:         make(chan string),
		wc:         make(chan string),
		path:       "/tmp/access.log",
		influxDBsn: "username&password...",
	}
	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxDB()
	time.Sleep(time.Second * 1)
}
