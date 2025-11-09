package models

import "time"

// Conversation 对话记录
type Conversation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID string    `gorm:"uniqueIndex;size:100;not null" json:"conversation_id"`
	UserID         string    `gorm:"index;size:50;not null" json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// Message 消息记录
type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID string    `gorm:"index;size:100;not null" json:"conversation_id"`
	Role           string    `gorm:"size:20;not null" json:"role"` // user / assistant
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

// StudyPlan 学习计划
type StudyPlan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PlanID    string    `gorm:"uniqueIndex;size:100;not null" json:"plan_id"`
	UserID    string    `gorm:"index;size:50;not null" json:"user_id"`
	Title     string    `gorm:"size:200" json:"title"`
	Goal      string    `gorm:"type:text" json:"goal"`
	Content   string    `gorm:"type:text" json:"content"` // JSON 格式存储
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (StudyPlan) TableName() string {
	return "study_plans"
}