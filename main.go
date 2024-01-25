package main

import (
	"ServerListen/Data"
	"ServerListen/Untils"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fishtailstudio/imgo"
	"github.com/go-ini/ini"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	ServerListen()
	//t2()
}

var (
	Server string
	Port   string
)

func ServerListen() {
	mux := http.NewServeMux()
	var c = make(chan string, 1)
	var mag Data.MagiciDemo
	fs := http.FileServer(http.Dir("./videos"))
	mux.Handle("/", http.StripPrefix("/videos", fs))
	mux.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		log.Println("开始回调")
		all, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		err2 := json.Unmarshal(all, &mag)
		if err2 != nil {
			Untils.Error.Println(err2.Error())
		}
		format := time.UnixMilli(mag.Ts).Format("2006-01-02-15:04:05")
		Untils.Info.Printf("CameraId ==%s ,AlarmType==%s,时间===%s,video===%s\n", mag.CameraId, mag.AlarmType, format, mag.Url)
		format2 := time.UnixMilli(mag.Ts).Format("20060102150405")
		var str1 = fmt.Sprintf("ResultFile/%s", mag.AlarmType)
		var str2 = fmt.Sprintf("%s/%s", str1, format2)
		stat, _ := os.Stat(str2)
		if stat == nil {
			os.MkdirAll(str2, 666)
		}
		go SaveLocalJPG(mag, str2)

		if len(mag.Url) > 1 {
			go SaveLocalVideo(mag, str2)
		}

	})
	mux.HandleFunc("/preview", func(writer http.ResponseWriter, request *http.Request) {
		Untils.Info.Println("凭证截留----")
		<-c
		Untils.Info.Println("凭证截留----完毕")
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
			Untils.Error.Println("文件错误", err2.Error())
			return
		}
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		file, header, err := r.FormFile("file") // 获取上传的文件
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 创建一个用于保存文件的本地路径
		dstPath := fmt.Sprintf("./videos/%s", header.Filename)
		UpFile := "/videos/" + header.Filename
		dstFile, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dstFile.Close()
		// 将上传的文件内容复制到目标文件
		_, err = io.Copy(dstFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := SaveUploadFile(UpFile)
		marshal, err := json.Marshal(data)
		if err != nil {
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			Untils.Error.Println(err.Error())
			return
		}

	})
	Untils.Info.Println("==========================回调开始===========================")
	addr := fmt.Sprintf("0.0.0.0:%s", Port)
	Untils.Info.Printf("=======================请使用该回调地址发起请求：%s/call =============================\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		return
	}
}

// SaveLocalJPG 将图片进行本地保存
func SaveLocalJPG(mag Data.MagiciDemo, str2 string) {
	var fileJPG = fmt.Sprintf("%s/%s%d.jpg", str2, mag.CameraId, mag.Ts)
	decode, err := base64.StdEncoding.DecodeString(mag.Scene)
	if err != nil {
		return
	}
	err1 := os.WriteFile(fileJPG, decode, 666)
	if err1 != nil {
		Untils.Error.Println("jpg图片写入失败=====", err1.Error())
		return
	}
	Untils.Info.Printf("当前本地存入图片地址：%s\n", fileJPG)
	defer func() {
		CanvasImage(fileJPG, str2, mag)
	}()
}

func init() {
	load, err := ini.Load("confs.ini")
	err2 := os.MkdirAll("./videos", os.ModePerm) // 创建存储视频文件的文件夹
	if err2 != nil {
		Untils.Error.Println(err2.Error())
	}
	if err != nil {
		fmt.Printf("初始化配置文件加载失败%s\n", err.Error())
		Untils.Error.Println("请按照设计文档来进行配置confs")
		return
	}

	Server = load.Section("IP").Key("ServerUrl").String()
	Port = load.Section("IP").Key("port").String()
}

// SaveUploadFile 视频存储到本地后-返回给接口的响应值
func SaveUploadFile(path string) map[string]any {
	m1 := md5.New()
	NowTime := time.Now().Unix()
	randStr := strconv.Itoa(int(NowTime))
	m1.Write([]byte(randStr))
	Code := hex.EncodeToString(m1.Sum(nil))
	data := map[string]any{
		"domain":  Server + "/",
		"md5":     Code,
		"mtime":   NowTime,
		"path":    path,
		"retcode": 0,
		"retmsg":  "",
		"scene":   "default",
		"scenes":  "default",
		"size":    132171,
		"src":     path,
		"url":     Server + "/" + path,
	}
	return data
}

// SaveLocalVideo 存储视频到本地
func SaveLocalVideo(mag Data.MagiciDemo, str2 string) {
	//拼接远程地址
	ServerVideoAddr := fmt.Sprintf("%s%s", Server, mag.Url)
	Untils.Info.Println("远程地址：", ServerVideoAddr)
	//创建一个同预期的路径文件
	VideoFile := fmt.Sprintf("%s/%d.mp4", str2, mag.Ts)
	Untils.Info.Println("本地地址：", VideoFile)
	get, err := http.Get(ServerVideoAddr)
	if err != nil {
		Untils.Error.Println(err.Error())
		return
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		Untils.Error.Println(err.Error())
		return
	}
	file, err := os.OpenFile(VideoFile, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		Untils.Error.Println(err.Error())
		return
	}
	_, err2 := file.Write(all)
	if err2 != nil {
		Untils.Error.Println(err2.Error())
		return
	}
	defer file.Close()
}

// CanvasImage 根据坐标点进行绘图-
func CanvasImage(LocalImage, str2 string, Data Data.MagiciDemo) {
	im := imgo.Load(LocalImage)
	SaveOnePath := fmt.Sprintf("%s/Canvas.png", str2)
	RandPng := fmt.Sprintf("%d.png", Data.Ts)
	for i := 0; i < len(Data.Boxes); i++ {
		imgo.Canvas(Data.Boxes[i].Width, Data.Boxes[i].Height, color.RGBA{R: 220, G: 20, B: 60, A: 80}).
			Save(RandPng)
		im.Insert(RandPng, Data.Boxes[i].X, Data.Boxes[i].Y)
	}
	im.Save(SaveOnePath)
	defer func() {
		//绘制完毕后--删除坐标图
		err := os.Remove(RandPng)
		if err != nil {
			Untils.Error.Println("删除失败：", err.Error())
			return
		}
		Untils.Info.Println("删除成功")
	}()
}
