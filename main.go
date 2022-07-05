package main

import (
	"ServerListen/Data"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	ServerListen()
}

var (
	Server string
	Port   string
)

func ServerListen() {

	mux := http.NewServeMux()
	var mag Data.MagiciDemo
	mux.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("开始回调")
		all, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		err2 := json.Unmarshal(all, &mag)
		if err2 != nil {
			fmt.Println(err2)
		}
		format := time.UnixMilli(mag.Ts).Format("2006-01-02-15:04:05")
		fmt.Printf("CameraId ==%s ,AlarmType==%s,时间===%s,video===%s\n", mag.CameraId, mag.AlarmType, format, mag.Url)
		str2, filestr, videoAddr := varData(mag)
		stat, _ := os.Stat(str2)
		if stat == nil {
			os.MkdirAll(str2, 666)
		}
		go writeJpg(mag, filestr)
		if len(mag.Url) > 1 {
			go writeVideo(mag, str2, videoAddr)
		}

	})
	fmt.Println("==========================回调开始===========================")
	addr := fmt.Sprintf("0.0.0.0:%s", Port)
	fmt.Printf("=======================请使用该回调地址发起请求：%s/call =============================\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		return
	}
}

func varData(mag Data.MagiciDemo) (v1, v2, v3 string) {
	format2 := time.UnixMilli(mag.Ts).Format("20060102150405")
	var str1 = fmt.Sprintf("ResultFile/%s", mag.AlarmType)
	var str2 = fmt.Sprintf("%s/%s", str1, format2)
	var filestr = fmt.Sprintf("%s/%d.jpg", str2, mag.Ts)
	videoAddr := fmt.Sprintf("%s%s", Server, mag.Url)
	return str2, filestr, videoAddr
}

func writeJpg(mag Data.MagiciDemo, filestr string) {
	decode, err := base64.StdEncoding.DecodeString(mag.Scene)
	if err != nil {
		return
	}
	err1 := ioutil.WriteFile(filestr, decode, 666)
	if err1 != nil {
		fmt.Println("jpg图片写入失败=====", err1.Error())
		return
	}
	fmt.Printf("当前本地存入图片地址：%s\n", filestr)
}

func writeVideo(mag Data.MagiciDemo, str2, videoAddr string) {
	videofile := fmt.Sprintf("%s/%d.mp4", str2, mag.Ts)
	vidoeFile, err := os.Create(videofile)
	if err != nil {
		return
	}
	get, err := http.Get(videoAddr)
	if err != nil {
		return
	}
	io.Copy(vidoeFile, get.Body)
	fmt.Printf("当前远端视频存放地址：%s\n", videoAddr)
	fmt.Printf("当前本地视频存放地址：%s\n", videofile)
	if err != nil {
		return
	}
}

func init() {
	load, err := ini.Load("confs.ini")
	if err != nil {
		fmt.Printf("初始化配置文件加载失败%s\n", err.Error())
		fmt.Println("请按照设计文档来进行配置confs")
		return
	}

	Server = load.Section("IP").Key("videoUrl").String()
	Port = load.Section("IP").Key("port").String()
}
