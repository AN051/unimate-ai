package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"unimate-ai/models"
	"unimate-ai/services"
)

func GeneratePlanHandler(c *gin.Context) {
	var req models.GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	systemPrompt := `你是一个专业的学习计划制定助手。请根据用户的目标、当前水平和时间限制，制定一个详细、可执行的学习计划。

计划要求：
1. 分阶段设置明确的里程碑
2. 每个阶段包含具体的学习内容和时间安排
3. 考虑用户的实际情况，制定切实可行的计划
4. 提供学习资源建议和复习策略`

	userMessage := "我的学习目标：" + req.Goal
	if req.CurrentLevel != "" {
		userMessage += "\n当前水平：" + req.CurrentLevel
	}
	if req.TimeLimit != "" {
		userMessage += "\n时间限制：" + req.TimeLimit
	}

	reply, aiResp, err := services.CallAI(systemPrompt, userMessage, 0.7, nil)
	if err != nil {
		log.Printf("❌ AI 调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	cost := float64(aiResp.Usage.TotalTokens) * 0.35 / 1000000

	var resp models.GeneratePlanResponse
	resp.Code = 200
	resp.Data.Plan = reply
	resp.Data.Suggestion = "建议每天坚持学习，保持良好的学习习惯。"
	resp.Data.TokenUsage.Total = aiResp.Usage.TotalTokens
	resp.Data.TokenUsage.Cost = cost

	c.JSON(http.StatusOK, resp)
}

func BreakdownGoalHandler(c *gin.Context) {
	var req models.BreakdownGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	systemPrompt := `你是一个目标拆解专家。请将用户的大目标拆解为具体的、可执行的小步骤。

要求：
1. 每个步骤清晰具体
2. 步骤之间有逻辑顺序
3. 提供大致的时间线
4. 步骤数量控制在 5-8 个`

	userMessage := "请帮我拆解这个目标：" + req.Goal

	reply, aiResp, err := services.CallAI(systemPrompt, userMessage, 0.7, nil)
	if err != nil {
		log.Printf("❌ AI 调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	cost := float64(aiResp.Usage.TotalTokens) * 0.35 / 1000000

	var resp models.BreakdownGoalResponse
	resp.Code = 200
	resp.Data.Steps = []string{reply}
	resp.Data.Timeline = "根据个人情况调整"
	resp.Data.TokenUsage.Total = aiResp.Usage.TotalTokens
	resp.Data.TokenUsage.Cost = cost

	c.JSON(http.StatusOK, resp)
}

func EvaluatePlanHandler(c *gin.Context) {
	var req models.EvaluatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	systemPrompt := `你是一个学习计划评估专家。请评估用户的学习计划，并给出改进建议。

评估维度：
1. 目标清晰度
2. 可执行性
3. 时间安排合理性
4. 资源配置
5. 持续性保障

请给出 0-100 的评分，并提供具体的改进建议。`

	userMessage := "请评估这个学习计划：\n" + req.Plan

	reply, aiResp, err := services.CallAI(systemPrompt, userMessage, 0.7, nil)
	if err != nil {
		log.Printf("❌ AI 调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	cost := float64(aiResp.Usage.TotalTokens) * 0.35 / 1000000

	var resp models.EvaluatePlanResponse
	resp.Code = 200
	resp.Data.Evaluation = reply
	resp.Data.Score = 85
	resp.Data.Suggestion = []string{"保持每日学习习惯", "定期复习重点内容"}
	resp.Data.TokenUsage.Total = aiResp.Usage.TotalTokens
	resp.Data.TokenUsage.Cost = cost

	c.JSON(http.StatusOK, resp)
}