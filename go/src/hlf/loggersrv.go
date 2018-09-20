package hlf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (

	//SysConfFile log system parameter file
	SysConfFile = "logsys.conf"

	//LogConfFile loggers parameter file
	LogConfFile = "logs.conf"
)

type logItem struct {
	target string
	text   string
}

var _ch = make(chan logItem)

var _srvRunning = false

var _logRoot = "./log/"

//Init turn on server
func init() {
	loadLogConf(SysConfFile, LogConfFile)
	dirCache := make(map[string]time.Time)
	t := time.Now()
	_logRoot += t.Format(time.RFC3339) + "/"

	go run(dirCache)
}

func run(dirCache map[string]time.Time) {
	_srvRunning = true

	for {
		select {
		case li := <-_ch:
			if li.target[:8] == "console:" {
				printc(li.text)
			} else if li.target[:5] == "file:" {
				path := li.target[5:]
				printf(&dirCache, path, li.text)
			}
		}
	}
}

func printc(text string) {
	fmt.Printf(text)
}

func printf(dc *map[string]time.Time, path string, text string) {
	_, found := (*dc)[path]
	if !found {
		c := strings.Split(path, "/")
		if len(c) > 1 {
			d := path[:len(path)-len(c[len(c)-1])-1]
			os.MkdirAll(d, 0700)
		}
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err == nil {
		_, err = f.WriteString(text)
	} else if os.IsNotExist(err) {
		err = ioutil.WriteFile(path, []byte(text), 0700)
	} else if err != nil {
		printc(text + "[FCF]")
	}

	(*dc)[path] = time.Now()
}
