package ProxyLog

import (
	"io"
	"log"
	"os"
)

var (
	Info *log.Logger
	Err  *log.Logger
)

func init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("LogInitErr ", err)
	}
	Info = log.New(io.MultiWriter(file, os.Stderr), "Info: ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	Err = log.New(io.MultiWriter(file, os.Stderr), "Err: ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
}

// TODO 日志分级
