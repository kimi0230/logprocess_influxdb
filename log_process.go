package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte   // read channel
	wc    chan *Message // write channel
	read  Reader
	write Writer
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// 系統狀態監控
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 總處理log行數
	Tps          float64 `json:"tps"`          // 系統吞吐量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 長度
	WriteChanLen int     `json:"writeChanLen"` // write channel 長度
	RunTime      string  `json:"runTime"`      // 總運行時間
	ErrNum       int     `json:"errNum"`       // 總錯誤
}

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

var TypeMonitorChan = make(chan int, 200)

type Monitor struct {
	startTime time.Time
	data      SystemInfo
	tpsSli    []int
}

func (m *Monitor) start(lp *LogProcess) {
	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.data.ErrNum += 1
			case TypeHandleLine:
				m.data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-ticker.C
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			if len(m.tpsSli) > 2 {
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)

		// 計算吞吐量
		if len(m.tpsSli) >= 2 {
			m.data.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5
		}

		ret, _ := json.MarshalIndent(m.data, "", "\t")
		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":9193", nil)
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
		TypeMonitorChan <- TypeHandleLine
		rc <- line[:len(line)-1]
	}
}

type WriteToInfluxDB struct {
	influxDBDsn string // influx data source
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	// 寫入模塊
	// http://127.0.0.1:8086@kimiORG@kk@myMeasure@s
	infSli := strings.Split(w.influxDBDsn, "@")

	// You can generate a Token from the "Tokens Tab" in the UI
	org := infSli[1]
	bucket := infSli[2]
	measure := infSli[3]

	client := influxdb2.NewClient("http://127.0.0.1:8086", token)
	// always close client at the end
	defer client.Close()
	client.Options()
	writeAPI := client.WriteAPI(org, bucket)

	// write channel 中讀取監控數據
	for v := range wc {
		// 構造數據並寫入influxdb
		// Tags: Path, Method, Scheme, Status
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status}

		// Fields: UpstreamTime, RequestTime, BytesSent
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		// Write the batch
		p := influxdb2.NewPoint(
			measure,
			tags,
			fields,
			v.TimeLocal)

		// write point asynchronously
		writeAPI.WritePoint(p)

		log.Println("write success!")
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
			TypeMonitorChan <- TypeErrNum
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}
		message := &Message{}

		// [04/Mar/2018:13:49:52 +0000]
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeErrNum
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
			TypeMonitorChan <- TypeErrNum
			log.Println("strings.Split fail", ret[6])
			continue
		}
		// GET
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			TypeMonitorChan <- TypeErrNum
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

const token = "7Vft2nXp1IkgLMu1VaLVEqylPKeJMqO1KLLfwRa1wxOg92DwMqHEjKkTqbqj03k49Inw-cD2rmBQOok-Dij2BQ=="

func main() {
	var path, influxDsn string
	flag.StringVar(&path, "path", "./access.log", "read file path")
	flag.StringVar(&influxDsn, "influxDsn", "http://127.0.0.1:8086@kimiORG@kk@myMeasure@s", "influx data source")
	flag.Parse()

	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: influxDsn,
	}

	lp := &LogProcess{
		rc:    make(chan []byte, 200), // 讀取的模塊會比解析來得快, 所以使用buffer的 channel
		wc:    make(chan *Message, 200),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	for i := 0; i < 2; i++ {
		go lp.Process()
	}
	for i := 0; i < 4; i++ {
		go lp.write.Write(lp.wc)
	}

	// 監控模組
	m := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	m.start(lp)
}
