package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc     chan []byte   // read channel
	wc     chan *Message // write channel
	read   Reader
	writer Writer
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
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
		// 寫入到 read channel
		rc <- line[:len(line)-1]
	}
}

type WriteToInfluxDB struct {
	influxDBsn string // influx data source
}

func (l *WriteToInfluxDB) Write(wc chan *Message) {
	// 寫入模塊
	for v := range wc {
		fmt.Println(v)
	}
}

func (l *LogProcess) Process() {
	// 解析模塊

	/**
	172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	*/

	// 正規提取所需的監控數據(path, status, method 等)
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Taipei")
	// 從 read channel 中讀取每行日誌數據
	for v := range l.rc {
		// 第 0 項是數據本身
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}
		message := &Message{}

		// [04/Mar/2018:13:49:52 +0000]
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		// 2133
		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail", ret[6])
			continue
		}
		// GET
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err)
			continue
		}
		message.Path = u.Path

		// http
		message.Scheme = ret[5]
		// 200
		message.Status = ret[7]

		// 1.005
		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		// 1.854
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		// 寫入 write channel
		l.wc <- message
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
		wc:     make(chan *Message),
		read:   r,
		writer: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)
	time.Sleep(time.Second * 30)
}
