package model

import (
	"fmt"
	"io"

	"github.com/RaymondCode/simple-demo/lib"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func VideoOss(file io.Reader, videoName string) error {
	conf := lib.LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		fmt.Println("创建实例Error:", err)
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		fmt.Println("获取存储空间Error:", err)
		return err
	}

	// 上传本地文件。
	err = bucket.PutObject(videoName, file)
	if err != nil {
		fmt.Println("本地文件上传Error:", err)
		return err
	}

	fmt.Println("上传文件成功！")

	return nil
}
