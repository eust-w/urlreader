package main

import (
	"github.com/eust-w/urlreader/api"
	"github.com/eust-w/urlreader/config"
	"github.com/eust-w/urlreader/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 初始化日志
	logger.InitLogger()
	log := logger.GetLogger()
	defer logger.SyncLogger()

	log.Infow("服务启动中...")

	// 加载配置
	cfg := config.LoadConfig()
	log.Infow("配置加载完成", "port", cfg.Port)

	// 创建Gin引擎
	router := gin.Default()

	// 设置 CORS，允许跨域 DELETE、GET、POST、OPTIONS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可根据需要指定前端域名
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	}))

	// 创建API处理程序
	handler := api.NewHandler(cfg)
	log.Infow("API Handler 初始化完成")

	// 设置路由
	handler.SetupRoutes(router)
	log.Infow("路由设置完成")

	// 启动清理任务
	handler.StartCleanupTask()
	log.Info("定时清理任务启动")

	// 启动服务器
	port := cfg.Port
	log.Infow("服务器正在启动，监听端口...", "port", port)
	// 启动服务
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Errorw("无法启动服务", "error", err)
	} else {
		log.Infow("服务已启动", "port", cfg.Port)
	}
}
