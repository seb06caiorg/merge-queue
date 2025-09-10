package models

import "time"

// APIResponse represents a standard API response format.
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationMeta represents pagination metadata.
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime,omitempty"`
}

// CreateTaskRequest represents a request to create a task.
type CreateTaskRequest struct {
	Title       string   `json:"title" validate:"required,max=200"`
	Description string   `json:"description" validate:"max=1000"`
	Status      string   `json:"status" validate:"omitempty,oneof=pending in-progress completed cancelled"`
	Priority    string   `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	AssignedTo  string   `json:"assigned_to" validate:"omitempty,max=50"`
	Tags        []string `json:"tags" validate:"omitempty,dive,max=50"`
}

// UpdateTaskRequest represents a request to update a task.
type UpdateTaskRequest struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=pending in-progress completed cancelled"`
	Priority    *string  `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
	AssignedTo  *string  `json:"assigned_to,omitempty" validate:"omitempty,max=50"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}
