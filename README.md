# simple-demo

### 说明
###TESTTTTTT
router注册路由，controller层按模块分为五个文件，各自负责各自的部分。

controller命名规则：模块名加Controller，例如：videoController.go

service,model 命名规则同上。

以“POST douyin/user/login/ 用户登录” 举例：

userController.go调用userService.go,userService.go调用userModel.go

其中userModel直接与数据库进行交互。数据库已经配置为远程服务器。

在开发过程中，若要使用token获取id，使用jwtService的ParseToken即可，具体使用可见jwt测试文件。可以参考我的dev分支写，或者那个优秀项目写。


### 项目运行

go mod tidy

go run main.go

使用任意api测试软件即可测试如下接口(POST)：

http://localhost:8080/douyin/user/login/?username=test&password=test


### 目录说明

conf:配置文件

lib:工具类函数

public：视频本地存储文件

router:路由


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




### 进度汇总

重构了代码架构以及添加了mysql、redis数据库

添加了用户注册、登录、信息的路由,提出了相关issure

