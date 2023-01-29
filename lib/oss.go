package lib

import (
	"fmt"
	"io"
	
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 说明：fileHash是存储在oss内部的文件名，提取文件时需要文件名以及其后缀
// 这里上传文件是读取本地文件后再上传

//上传文件至阿里云
func UploadOss(filename, fileHash string) {
	//获取文件后缀
	fileSuffix := path.Ext(filename)
	conf := LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		fmt.Println("创建实例Error:", err)
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		fmt.Println("获取存储空间Error:", err)
		return
	}

	// 上传本地文件。
	err = bucket.PutObjectFromFile("files/"+fileHash+fileSuffix, "/Users/toma/Projects/GoLearning/douyin-minimalist-service/public/"+filename)
	if err != nil {
		fmt.Println("本地文件上传Error:", err)
		return
	}

	fmt.Println("上传文件成功！")
}

//从oss下载文件
func DownloadOss(fileHash, fileType string) []byte {
	conf := LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 下载文件到流。
	body, err := bucket.GetObject("files/" + fileHash + fileType)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return data
}

//从oss删除文件
func DeleteOss(fileHash, fileType string) {
	conf := LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = bucket.DeleteObject("files/" + fileHash + fileType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}