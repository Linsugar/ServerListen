package Untils

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	var str2 = "LogFile"
	stat, _ := os.Stat(str2)
	if stat == nil {
		os.MkdirAll(str2, 666)
	}
	time.Sleep(time.Second)
	infoFile, err1 := os.OpenFile("LogFile/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	warnFile, err2 := os.OpenFile("LogFile/warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	errFile, err3 := os.OpenFile("LogFile/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatalln("打开日志文件失败：", err3.Error())
	}
	//Info = log.New(os.Stdout, "[<<<大脑Info日志收集>>>]", log.Ldate|log.Ltime|log.Lshortfile)
	//Warning = log.New(os.Stdout, "[<<<大脑Warning日志收集>>>]", log.Ldate|log.Ltime|log.Lshortfile)
	//Error = log.New(io.MultiWriter(os.Stderr, errFile), "[<<<大脑Error日志收集>>>]", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(os.Stderr, infoFile), "[<<<大脑Info日志收集>>>]", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stderr, warnFile), "[<<<大脑Warning日志收集>>>]:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "[<<<大脑Error日志收集>>>]", log.Ldate|log.Ltime|log.Lshortfile)
}
