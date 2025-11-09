package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"unimate-ai/config"
	"unimate-ai/models"
	"unimate-ai/services"
)

func ChatHandler(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.UserID == "" {
		req.UserID = "default_user"
	}

	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = fmt.Sprintf("user_%s_conv_%d", req.UserID, time.Now().Unix())
		
		// 创建新对话记录
		conversation := models.Conversation{
			ConversationID: conversationID,
			UserID:         req.UserID,
		}
		config.DB.Create(&conversation)
	}

	systemPrompt := `你是"太傅AI"，一个温暖、耐心的大学生学习助手。你的职责是：
1. 提供学习方法和计划建议
2. 给予情绪支持和鼓励
3. 帮助用户克服学习困难
4. 像朋友一样陪伴用户成长

回答要求：
- 语气温暖、真诚、有同理心
- 提供具体可行的建议
- 适当使用 emoji 表情
- 回答简洁（控制在 200 字以内）`

	// 从数据库加载历史消息（最近20条）
	var dbMessages []models.Message
	config.DB.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Limit(20).
		Find(&dbMessages)

	var history []services.Message
	for _, msg := range dbMessages {
		history = append(history, services.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	log.Printf("💬 [对话] 用户=%s, 对话ID=%s, 历史消息数=%d", req.UserID, conversationID, len(history))

	reply, aiResp, err := services.CallAI(systemPrompt, req.Message, 0.7, history)
	if err != nil {
		log.Printf("❌ AI 调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	cost := float64(aiResp.Usage.TotalTokens) * 0.35 / 1000000

	log.Printf("💰 Token使用: 输入=%d, 输出=%d, 总计=%d, 费用=¥%.6f",
		aiResp.Usage.PromptTokens,
		aiResp.Usage.CompletionTokens,
		aiResp.Usage.TotalTokens,
		cost)

	// 保存消息到数据库
	userMsg := models.Message{
		ConversationID: conversationID,
		Role:           "user",
		Content:        req.Message,
	}
	assistantMsg := models.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        reply,
	}
	config.DB.Create(&userMsg)
	config.DB.Create(&assistantMsg)

	// 更新对话的 updated_at
	config.DB.Model(&models.Conversation{}).
		Where("conversation_id = ?", conversationID).
		Update("updated_at", time.Now())

	var resp models.ChatResponse
	resp.Code = 200
	resp.Data.Reply = reply
	resp.Data.ConversationID = conversationID
	resp.Data.TokenUsage.Input = aiResp.Usage.PromptTokens
	resp.Data.TokenUsage.Output = aiResp.Usage.CompletionTokens
	resp.Data.TokenUsage.Total = aiResp.Usage.TotalTokens
	resp.Data.TokenUsage.Cost = cost

	c.JSON(http.StatusOK, resp)
}

func GetConversationsHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 user_id"})
		return
	}

	var conversations []models.Conversation
	config.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&conversations)

	var result []map[string]interface{}
	for _, conv := range conversations {
		var messageCount int64
		config.DB.Model(&models.Message{}).
			Where("conversation_id = ?", conv.ConversationID).
			Count(&messageCount)

		var lastMsg models.Message
		config.DB.Where("conversation_id = ?", conv.ConversationID).
			Order("created_at DESC").
			First(&lastMsg)

		lastMessage := lastMsg.Content
		if len(lastMessage) > 50 {
			lastMessage = lastMessage[:50] + "..."
		}

		result = append(result, map[string]interface{}{
			"conversation_id": conv.ConversationID,
			"message_count":   messageCount,
			"last_message":    lastMessage,
			"updated_at":      conv.UpdatedAt,
		})
	}

	log.Printf("📋 [获取对话列表] 用户=%s, 对话数=%d", userID, len(result))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"conversations": result}})
}

func DeleteConversationHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	
	// 删除对话和所有消息
	config.DB.Where("conversation_id = ?", conversationID).Delete(&models.Conversation{})
	config.DB.Where("conversation_id = ?", conversationID).Delete(&models.Message{})
	
	log.Printf("🗑️ [删除对话] 对话ID=%s", conversationID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}