# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./simple-demo
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试


### apk模拟器
mac m1版:
https://developer.android.com/studio
https://blog.csdn.net/qq285744011/article/details/126200016

### 开发相关文档
https://bytedance.feishu.cn/docx/WZDddh2Lqoyfu6x93u1c8km9nug
https://bytedance.feishu.cn/docs/doccnKrCsU5Iac6eftnFBdsXTof#
https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7
https://www.apifox.cn/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707530