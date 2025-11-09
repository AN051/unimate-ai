package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// 硅基流动 API 配置
	SiliconFlowAPIKey = "sk-wihissfoxyzeqmjxqxsagfoedrgyykfqrrzcypwmwucjpwlr"
	SiliconFlowAPIURL = "https://api.siliconflow.cn/v1/chat/completions"  // ✅ 硅基流动地址
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CallAI 调用 AI（主对话功能）
func CallAI(systemPrompt, userMessage string, temperature float64, history []Message) (string, AIResponse, error) {
	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
	}

	for _, msg := range history {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": userMessage,
	})

	// ✅ 使用 Qwen2.5-7B 模型（硅基流动平台）
	reqBody := map[string]interface{}{
		"model":       "Qwen/Qwen2.5-7B-Instruct",
		"messages":    messages,
		"temperature": temperature,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", SiliconFlowAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", AIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+SiliconFlowAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", AIResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", AIResponse{}, fmt.Errorf("JSON解析失败: %v, 响应: %s", err, string(body))
	}

	if len(result.Choices) == 0 {
		return "", AIResponse{}, fmt.Errorf("AI 返回为空, 响应: %s", string(body))
	}

	var aiResp AIResponse
	aiResp.Usage = result.Usage

	return result.Choices[0].Message.Content, aiResp, nil
}

// ChatWithContext 调用 AI（学习计划功能）
func ChatWithContext(messages []map[string]string, modelName string) (string, map[string]interface{}, error) {
	// ✅ 默认使用 Qwen2.5-7B 模型
	if modelName == "" {
		modelName = "Qwen/Qwen2.5-7B-Instruct"
	}

	reqBody := map[string]interface{}{
		"model":    modelName,
		"messages": messages,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", SiliconFlowAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+SiliconFlowAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage map[string]interface{} `json:"usage"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, fmt.Errorf("JSON解析失败: %v, 响应: %s", err, string(body))
	}

	if len(result.Choices) == 0 {
		return "", nil, fmt.Errorf("AI 返回为空, 响应: %s", string(body))
	}

	return result.Choices[0].Message.Content, result.Usage, nil
}