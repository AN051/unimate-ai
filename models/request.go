package models

// ChatRequest AI 对话请求
type ChatRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
}

// ChatResponse AI 对话响应
type ChatResponse struct {
	Code int `json:"code"`
	Data struct {
		Reply          string `json:"reply"`
		ConversationID string `json:"conversation_id"`
		TokenUsage     struct {
			Input  int     `json:"input"`
			Output int     `json:"output"`
			Total  int     `json:"total"`
			Cost   float64 `json:"cost"`
		} `json:"token_usage"`
	} `json:"data"`
}

// GeneratePlanRequest 生成学习计划请求
type GeneratePlanRequest struct {
	Goal         string   `json:"goal" binding:"required"`
	CurrentLevel string   `json:"current_level"`
	TimeLimit    string   `json:"time_limit"`
	Preferences  []string `json:"preferences"`
	UserID       string   `json:"user_id"`
}

// GeneratePlanResponse 生成学习计划响应
type GeneratePlanResponse struct {
	Code int `json:"code"`
	Data struct {
		Plan       string `json:"plan"`
		Suggestion string `json:"suggestion"`
		TokenUsage struct {
			Total int     `json:"total"`
			Cost  float64 `json:"cost"`
		} `json:"token_usage"`
	} `json:"data"`
}

// BreakdownGoalRequest 目标拆解请求
type BreakdownGoalRequest struct {
	Goal   string `json:"goal" binding:"required"`
	UserID string `json:"user_id"`
}

// BreakdownGoalResponse 目标拆解响应
type BreakdownGoalResponse struct {
	Code int `json:"code"`
	Data struct {
		Steps      []string `json:"steps"`
		Timeline   string   `json:"timeline"`
		TokenUsage struct {
			Total int     `json:"total"`
			Cost  float64 `json:"cost"`
		} `json:"token_usage"`
	} `json:"data"`
}

// EvaluatePlanRequest 计划评估请求
type EvaluatePlanRequest struct {
	Plan   string `json:"plan" binding:"required"`
	UserID string `json:"user_id"`
}

// EvaluatePlanResponse 计划评估响应
type EvaluatePlanResponse struct {
	Code int `json:"code"`
	Data struct {
		Evaluation string   `json:"evaluation"`
		Score      int      `json:"score"`
		Suggestion []string `json:"suggestion"`
		TokenUsage struct {
			Total int     `json:"total"`
			Cost  float64 `json:"cost"`
		} `json:"token_usage"`
	} `json:"data"`
}