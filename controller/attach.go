package controller

import (
	"fmt"
	"im/utils"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// todo 定期更新
const (
	AccessKeyId     = "5p2RZKnrUanMuQw9"
	AccessKeySecret = "bsNmjU8Au08axedV40TRPCS5XIFAkK"
	EndPoint        = "oss-cn-shenzhen.aliyuncs.com"
	Bucket          = "winliondev"
)

// 存储位置 ./upload,需要确保已经创建好
func init() {
	os.MkdirAll("./upload", os.ModePerm)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	UploadLocal(w, r)
	//UploadOss(w,r)
}

// url格式 /upload/xxxx.png  需要确保网络能访问/upload/
func UploadLocal(w http.ResponseWriter, req *http.Request) {
	// 获得上传的源文件s
	srcfile, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	// 创建一个新文件d
	suffix := ".png"
	// 如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	//formdata.append("filetype",".png")
	filetype := req.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstfile, err := os.Create("./upload/" + filename)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 将源文件内容copy到新文件
	_, err = io.Copy(dstfile, srcfile)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 将新文件路径转换成url地址
	url := "/upload/" + filename
	// 响应到前端
	utils.RespOK(w, url, "")
}

//需要安装, 权限设置为公共读状态
func UploadOss(w http.ResponseWriter, req *http.Request) {
	// 获得上传的文件
	srcfile, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 获得文件后缀.png/.mp3
	suffix := ".png"
	//如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	//formdata.append("filetype",".png")
	filetype := req.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	// 初始化ossclient
	client, err := oss.New(EndPoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 获得bucket
	bucket, err := client.Bucket(Bucket)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 设置文件名称
	//time.Now().Unix()
	filename := fmt.Sprintf("upload/%d%04d%s",
		time.Now().Unix(), rand.Int31(),
		suffix)
	// 通过bucket上传
	err = bucket.PutObject(filename, srcfile)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	// 获得url地址
	url := "http://" + Bucket + "." + EndPoint + "/" + filename
	// 响应到前端
	utils.RespOK(w, url, "")
}
