package main

import (
	"asynctaskhub/global"
	_ "asynctaskhub/src/router"
	"asynctaskhub/src/service"
	"asynctaskhub/src/service/queue"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"html/template"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	global.InitializeConfig()
	global.InitializeLogger(global.CONFIG)

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		global.Logger.Fatal("failed to load location", zap.Error(err))
		return
	}
	time.Local = loc

	if err := global.InitializeDB(global.CONFIG, global.Logger); err != nil {
		global.Logger.Fatal("failed to initialize DB", zap.Error(err))
		return
	}

	if err := global.InitializeRedis(global.CONFIG, global.Logger); err != nil {
		global.Logger.Fatal("failed to initialize Redis", zap.Error(err))
		return
	}

	switch global.CONFIG.APP_ENV {
	case "production":
		gin.SetMode(gin.ReleaseMode) // 生产环境
	case "test":
		gin.SetMode(gin.TestMode) // 测试环境
	default:
		gin.SetMode(gin.DebugMode) // 默认是开发环境
	}

	if len(os.Args) > 1 && os.Args[1] == "--test" {
		test()
		return
	}

	//命令行初始化数据库
	if len(os.Args) > 1 && os.Args[1] == "--initdb" {
		service.NewDatabaseService().InitDB()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	taskScheduler := queue.NewTaskScheduler()
	taskScheduler.StartTaskQueueListener(ctx)
	go taskScheduler.StartTaskQueueMonitor(ctx)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02") // 年-月-日格式
		},
		"safeHTML": func(htmlContent string) template.HTML {
			return template.HTML(htmlContent) // 将字符串标记为安全 HTML
		},
	})

	r = global.InitRouter(r)
	go func() {
		if err := r.Run(":8080"); err != nil {
			global.Logger.Fatal("Server error", zap.Error(err))
		}
	}()

	go func() {
		<-signalChan
		global.Logger.Info("Shutting down gracefully...")
		cancel()
	}()

	<-ctx.Done()
	time.Sleep(3 * time.Second)
	global.Logger.Info("Server stopped.")
}

func test() {
	global.Logger.Info("test")
}
