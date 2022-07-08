package main

import (
	"ServerListen/Data"
	"bytes"
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
	//t2()
}

var (
	Server string
	Port   string
	Ht     http.Client
)

func ServerListen() {
	mux := http.NewServeMux()
	var c = make(chan string, 1)
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
		format2 := time.UnixMilli(mag.Ts).Format("20060102150405")
		var str1 = fmt.Sprintf("ResultFile/%s", mag.AlarmType)
		var str2 = fmt.Sprintf("%s/%s", str1, format2)
		var fileJPG = fmt.Sprintf("%s/%s%d.jpg", str2, mag.CameraId, mag.Ts)
		var previewJPG = fmt.Sprintf("%s/%dprview.jpg", str2, mag.Ts)
		_, _, videoAddr := varData(mag, false)
		stat, _ := os.Stat(str2)
		if stat == nil {
			os.MkdirAll(str2, 666)
		}
		writeJpg(mag, fileJPG)
		if len(mag.Url) > 1 {
			writeVideo(mag, str2, videoAddr)
		}
		defer func() {
			c <- previewJPG
		}()

	})
	mux.HandleFunc("/preview", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("凭证截留----")
		<-c
		fmt.Println("凭证截留----完毕")
		all, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return
		}
		format2 := time.UnixMilli(time.Now().UnixMilli()).Format("20060102150405")
		var Perview = fmt.Sprintf("ResultFile/Perview")
		var PerviewFile = fmt.Sprintf("%s/%s.jpg", Perview, format2)
		stat, _ := os.Stat(Perview)
		if stat == nil {
			os.MkdirAll(Perview, 666)
			return
		}
		err2 := ioutil.WriteFile(PerviewFile, all, 666)
		if err2 != nil {
			fmt.Println("文件错误", err2.Error())
			return
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

func varData(mag Data.MagiciDemo, prf bool) (v1, v2, v3 string) {
	format2 := time.UnixMilli(mag.Ts).Format("20060102150405")
	var str1 = fmt.Sprintf("ResultFile/%s", mag.AlarmType)
	var str2 = fmt.Sprintf("%s/%s", str1, format2)
	if prf {
		var filestr = fmt.Sprintf("%s/%dprview.jpg", str2, mag.Ts)
		videoAddr := fmt.Sprintf("%s%s", Server, mag.Url)
		return str2, filestr, videoAddr
	} else {
		var filestr = fmt.Sprintf("%s/%d.jpg", str2, mag.Ts)
		videoAddr := fmt.Sprintf("%s%s", Server, mag.Url)
		return str2, filestr, videoAddr
	}

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

type uploadfile struct {
	File     string `json:"file"`
	Filepath string `json:"filepath"`
}

func t2() {
	url := "http://192.168.2.49:38095/group1/upload"
	var f uploadfile
	f.File = "file"
	f.Filepath = "t2.mp4"
	marshal, err := json.Marshal(f)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshal))
	if err != nil {
		return
	}

	//open, err := os.Open("t2.mp4")
	//if err != nil {
	//	return
	//}
	//defer open.Close()
	//body := bufio.NewReader(open)
	//writer := multipart.NewWriter(payloadBuf)
	//request.Header.Add("Content-Type", "multipart/form-data;boundary="+multipart.NewWriter(bytes.NewBufferString("")).Boundary())
	request.Header.Del("Content-Type")
	do, err := Ht.Do(request)
	if err != nil {
		return
	}
	all, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return
	}
	fmt.Println(string(all))
	//post, err := http.Post(url, "multipart/form-data", bytes.NewReader(marshal))
	//if err != nil {
	//	return
	//}
	//all, err := ioutil.ReadAll(post.Body)
	//if err != nil {
	//	return
	//}
	//
	//fmt.Println(string(all))
}
