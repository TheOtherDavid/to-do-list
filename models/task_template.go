package models

import "time"

type TaskTemplate struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	RecurrenceRule string    `json:"recurrence_rule"` // e.g., "DAILY", "WEEKLY", "FORTNIGHTLY", "MONTHLY"
	CreatedAt      time.Time `json:"created_at"`
}
