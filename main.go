package main

import (
"log"

"github.com/gin-gonic/gin"
"unimate-ai/config"
"unimate-ai/handlers"
"unimate-ai/models"
)

func main() {
if err := config.InitDB(); err != nil {
log.Fatalf("数据库初始化失败: %v", err)
}

if err := config.AutoMigrate(&models.StudyPlan{}); err != nil {
log.Fatalf("数据表创建失败: %v", err)
}

r := gin.Default()

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

r.GET("/health", func(c *gin.Context) {
c.JSON(200, gin.H{"status": "ok", "message": "AI 学习计划服务运行中"})
})

api := r.Group("/api/ai")
{
api.POST("/generate-plan", handlers.GeneratePlanHandler)
api.PUT("/plans/:plan_id/stages/:stage_id", handlers.UpdateStageHandler)
}

log.Println(" AI 学习计划服务启动成功！")
log.Println(" 接口地址: http://localhost:8080")

r.Run(":8080")
}