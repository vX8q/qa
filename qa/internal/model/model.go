package model

import "time"

type Question struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Answers   []Answer  `gorm:"constraint:OnDelete:CASCADE" json:"answers"`
}

type Answer struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID int       `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:text;not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
