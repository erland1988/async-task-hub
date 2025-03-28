# Async Task Hub项目说明

## 一、项目用途
Async Task Hub致力于满足企业和开发者在异步任务处理方面的需求。在实际业务场景中，如电商平台订单生成后，会触发库存更新、积分计算、消息通知等一系列异步任务，Async Task Hub可有效管理这些任务的执行流程，确保任务高效、稳定运行。

## 二、项目技术架构
Async Task Hub基于Go语言开发，采用MySQL存储数据，Redis管理任务队列，并通过分层架构设计，使各模块职责清晰，具备良好的可维护性与扩展性。

## 三、项目作者
项目作者：erland

## 四、交流方式
交流QQ群：1027932757

## 五、开源协议
本项目采用MIT开源协议，详情请见项目根目录下的`LICENSE`文件。

## 六、项目目录结构
```
./
├── common/               // 通用模块目录
├── data/                 // 数据目录，如日志等
├── global/               // 全局配置等相关目录
├── public/               // 公共资源目录
├── services/             // 服务相关目录
├── src/                  // 项目源代码目录
├── static/               // 静态资源目录
├── vue/                  // Vue相关代码目录
├──.env.production       // 生产环境环境变量配置文件
├── docker-compose.yml    // Docker Compose配置文件
├── Dockerfile            // Docker镜像构建文件
├── go.mod                // Go模块依赖管理文件
├── go.sum                // Go模块依赖版本文件
├── main.go               // 项目入口文件
├── package.json          // Node.js项目依赖管理文件
├── package-lock.json     // Node.js项目依赖版本锁定文件
├── LICENSE               // 开源协议文件
└── README.md             // 项目说明文件
```

## 七、环境变量配置
项目依靠环境变量进行配置，涵盖应用运行环境、访问地址、日志模式、数据库连接及Redis连接等关键信息。

## 八、部署方法
### （一）环境变量配置
开发与部署时，**必须**复制`.env.production`为`.env`，并根据实际需求在`.env`文件中配置环境变量。若不进行此操作，项目将无法正常运行。

### （二）Docker-Compose部署
1. **确保环境准备**：服务器需安装Docker和Docker - Compose。
2. **构建与启动**：
   - 进入项目根目录。
   - 执行命令`docker-compose up -d`。此命令会根据`docker-compose.yml`文件构建镜像并启动容器。`async_task_hub_nginx`服务基于`nginx:1.19.1 - alpine`镜像，负责网络代理，映射端口`8083:9090`，挂载配置文件和日志目录。`async_task_hub_app`服务构建项目自身镜像，设置副本数为3，暴露端口`9090`，挂载日志目录，并传递`.env`文件中的环境变量。
3. **验证运行**：通过浏览器访问`http://服务器IP:8083/task/`，若能正常访问，说明项目已成功部署。

## 九、项目核心功能
1. **任务队列服务**：实现任务的入队、出队、处理以及丢失任务恢复等功能。能根据任务执行时间将任务有序加入队列，并按序弹出执行，同时处理任务执行过程中的各种状态及重试逻辑。
2. **执行器客户端**：负责向执行器发送任务请求，构建包含必要信息的请求头，处理请求发送过程中的超时、错误等情况，并返回执行结果。
3. **任务调度器**：启动任务队列监听器和监控器。监听器通过多线程处理任务队列中的任务，监控器定时恢复丢失任务并动态更新执行器超时时间等配置。

## 十、创建队列接口
系统对外提供创建队列的业务接口，示例如下：
```bash
curl --location --request POST 'http://服务器IP:8083/task/task/api/taskqueue/create' \
--header 'X - App - Key: key_3_1' \
--header 'X - App - Secret: secret' \
--header 'Content - Type: application/json' \
--data - raw '{
    "task_code": "code_7_1",
    "parameters": "{\"a\": \"1\",\"b\": \"2\",\"c\": \"3\"}",
    "relative_delay_time": 60
}'
```
此接口方便外部系统与Async Task Hub交互，灵活创建任务队列。

## 十一、项目入口（main.go）
项目入口`main.go`主要完成配置初始化、日志初始化、数据库与Redis连接初始化，根据环境变量设置Gin框架运行模式，并支持特定命令行参数。同时启动任务调度器、清理服务监控，初始化Gin路由并启动HTTP服务，实现优雅关闭服务。 