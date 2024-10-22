package models

import "time"

type TaskInstance struct {
	ID          string     `json:"id"`
	TemplateID  string     `json:"template_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
