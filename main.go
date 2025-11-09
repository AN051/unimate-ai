package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"unimate-ai/config"
	"unimate-ai/handlers"
	"unimate-ai/models"
)

func main() {
	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}

	// 自动创建表
	if err := config.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.StudyPlan{},
	); err != nil {
		log.Fatalf("❌ 数据表创建失败: %v", err)
	}

	r := gin.Default()

	// CORS 跨域配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "AI 服务运行中"})
	})

	// API 路由组
	api := r.Group("/api/ai")
	{
		// AI 对话相关
		api.POST("/chat", handlers.ChatHandler)
		api.GET("/conversations", handlers.GetConversationsHandler)
		api.DELETE("/conversations/:conversation_id", handlers.DeleteConversationHandler)

		// 学习计划相关
		api.POST("/generate-plan", handlers.GeneratePlanHandler)
		api.POST("/breakdown-goal", handlers.BreakdownGoalHandler)
		api.POST("/evaluate-plan", handlers.EvaluatePlanHandler)
	}

	log.Println("🚀 AI 后端服务启动成功！")
	log.Println("📝 接口地址: http://localhost:8080")
	log.Println("💾 数据存储: SQLite 本地文件（unimate_ai.db）")
	log.Println("📦 部署提示: 迁移到服务器时，修改 config/database.go 切换到 MySQL")

	r.Run(":8080")
}