package models

import "time"

type GeneratePlanRequest struct {
	UserID     string `json:"user_id"`
	Flag       string `json:"flag" binding:"required"`
	Difficulty string `json:"difficulty"`
	Deadline   string `json:"deadline"`
	DailyHours int    `json:"daily_hours" binding:"required,min=1,max=24"`
}

type Stage struct {
	StageID     int    `json:"stage_id"`
	StageName   string `json:"stage_name"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	IsCompleted bool   `json:"is_completed"`
}

type Evaluation struct {
	Score     int      `json:"score"`
	Details   string   `json:"details"`
	KeyPoints []string `json:"key_points"`
}

type StudyPlan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlanID     string    `json:"plan_id" gorm:"uniqueIndex;size:100"`
	UserID     string    `json:"user_id" gorm:"index;size:50"`
	Flag       string    `json:"flag" gorm:"size:200"`
	Difficulty string    `json:"difficulty" gorm:"size:50"`
	Deadline   string    `json:"deadline" gorm:"size:20"`
	DailyHours int       `json:"daily_hours"`
	Stages     string    `json:"stages" gorm:"type:text"`
	Evaluation string    `json:"evaluation" gorm:"type:text"`
	Progress   int       `json:"progress"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (StudyPlan) TableName() string {
	return "study_plans"
}

type GeneratePlanResponse struct {
	Code int `json:"code"`
	Data struct {
		PlanID     string     `json:"plan_id"`
		Flag       string     `json:"flag"`
		Difficulty string     `json:"difficulty"`
		Deadline   string     `json:"deadline"`
		DailyHours int        `json:"daily_hours"`
		Stages     []Stage    `json:"stages"`
		Evaluation Evaluation `json:"evaluation"`
		TokenUsage struct {
			Total int     `json:"total"`
			Cost  float64 `json:"cost"`
		} `json:"token_usage"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
}

type UpdateStageRequest struct {
	IsCompleted *bool `json:"is_completed" binding:"required"`
}

type UpdateStageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PlanID      string `json:"plan_id"`
		StageID     int    `json:"stage_id"`
		IsCompleted bool   `json:"is_completed"`
		Progress    string `json:"progress"`
	} `json:"data"`
}