package models

import (
	"fmt"
	"time"
)

// Task represents a task in our system.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // "pending", "in-progress", "completed", "cancelled"
	Priority    string    `json:"priority"` // "low", "medium", "high", "critical"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

// TaskFilter represents filtering options for tasks.
type TaskFilter struct {
	Status     string   `json:"status,omitempty"`
	Priority   string   `json:"priority,omitempty"`
	AssignedTo string   `json:"assigned_to,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}

// TaskSearchQuery represents a search query for tasks.
type TaskSearchQuery struct {
	Query    string     `json:"query"`
	Fields   []string   `json:"fields"` // Fields to search in: "title", "description"
	Filters  TaskFilter `json:"filters"`
	SortBy   string     `json:"sort_by"` // "created_at", "updated_at", "priority"
	SortDesc bool       `json:"sort_desc"`
}

// TaskStats provides statistics about tasks.
type TaskStats struct {
	TotalTasks      int            `json:"total_tasks"`
	TasksByStatus   map[string]int `json:"tasks_by_status"`
	TasksByPriority map[string]int `json:"tasks_by_priority"`
	TasksByUser     map[string]int `json:"tasks_by_user"`
	LastUpdated     time.Time      `json:"last_updated"`
}

// Validation methods for Task.

// Validate checks if the task has valid data.
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if len(t.Title) > 200 {
		return fmt.Errorf("task title must be less than 200 characters")
	}
	if len(t.Description) > 1000 {
		return fmt.Errorf("task description must be less than 1000 characters")
	}
	if !IsValidStatus(t.Status) {
		return fmt.Errorf("invalid task status: %s", t.Status)
	}
	if !IsValidPriority(t.Priority) {
		return fmt.Errorf("invalid task priority: %s", t.Priority)
	}
	return nil
}

// IsValidStatus checks if the status is valid.
func IsValidStatus(status string) bool {
	validStatuses := []string{"pending", "in-progress", "completed", "cancelled"}
	for _, v := range validStatuses {
		if v == status {
			return true
		}
	}
	return false
}

// IsValidPriority checks if the priority is valid.
func IsValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high", "critical"}
	for _, v := range validPriorities {
		if v == priority {
			return true
		}
	}
	return false
}

// GetValidStatuses returns all valid task statuses.
func GetValidStatuses() []string {
	return []string{"pending", "in-progress", "completed", "cancelled"}
}

// GetValidPriorities returns all valid task priorities.
func GetValidPriorities() []string {
	return []string{"low", "medium", "high", "critical"}
}
