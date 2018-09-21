package hlf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	_ConfFile = "hlf.conf"
)

type logItem struct {
	target string
	text   string
}

var _logSrvCh = make(chan logItem)

var _logRoot string

func init() {
	//load settings
	err := loadConfFile(_ConfFile, &_conf)

	if err != nil {
		fmt.Println("[Warning][HLF] Failed to load settings: " + _ConfFile + " " + err.Error())
	}

	//create root directory for the log session
	_logRoot = generateLogRoot()

	//turn server on
	go run()
}

func generateLogRoot() string {
	t := time.Now()
	logRoot := _conf.LogRoot + t.Format(time.RFC3339) + "/"
	logRoot = strings.Replace(logRoot, "-", "", -1)
	logRoot = strings.Replace(logRoot, ":", "", -1)
	logRoot = strings.Replace(logRoot, "+", "", -1)
	return logRoot
}

func run() {
	for {
		select {
		case li := <-_logSrvCh:
			if li.target[:8] == "console:" {
				send2c(li.text)
			} else if li.target[:5] == "file:" {
				path := li.target[5:]
				send2f(path, li.text)
			}
		}
	}
}

func send2c(text string) {
	fmt.Printf(text)
}

func send2f(path string, text string) {
	c := strings.Split(path, "/")
	if len(c) > 1 {
		d := path[:len(path)-len(c[len(c)-1])-1]
		os.MkdirAll(d, 0700)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err == nil {
		_, err = f.WriteString(text)
	} else if os.IsNotExist(err) {
		err = ioutil.WriteFile(path, []byte(text), 0700)
	} else if err != nil {
		send2c(text + "[FCF]")
	}
}
