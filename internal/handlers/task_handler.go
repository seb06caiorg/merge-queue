package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"merge-queue/internal/models"
	"merge-queue/internal/services"
	"merge-queue/pkg/utils"
)

// TaskHandler handles HTTP requests for task operations.
type TaskHandler struct {
	taskService *services.TaskService
	response    *utils.ResponseHelper
	validator   *utils.ValidationUtils
	logger      *utils.Logger
}

// NewTaskHandler creates a new TaskHandler instance.
func NewTaskHandler(taskService *services.TaskService, logger *utils.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		response:    utils.NewResponseHelper(),
		validator:   utils.NewValidationUtils(),
		logger:      logger,
	}
}

// GetTasks handles GET /tasks requests.
func (th *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	th.logger.Debug("Getting tasks with filters")

	// Parse query parameters for filtering.
	filter := &models.TaskFilter{
		Status:     r.URL.Query().Get("status"),
		Priority:   r.URL.Query().Get("priority"),
		AssignedTo: r.URL.Query().Get("assigned_to"),
	}

	// Parse pagination parameters.
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Parse tags filter.
	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		filter.Tags = []string{tagsStr} // Simple implementation - could support multiple tags.
	}

	tasks, err := th.taskService.GetAllTasks(filter)
	if err != nil {
		th.logger.Error("Failed to get tasks: %v", err)
		th.response.SendError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}

	th.response.SendSuccess(w, response)
}

// GetTask handles GET /tasks/{id} requests.
func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		th.response.SendError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	th.logger.Debug("Getting task with ID: %d", id)

	task, err := th.taskService.GetTask(id)
	if err != nil {
		th.logger.Warn("Task not found: %d", id)
		th.response.SendError(w, http.StatusNotFound, "Task not found")
		return
	}

	th.response.SendSuccess(w, task)
}

// CreateTask handles POST /tasks requests.
func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	th.logger.Debug("Creating new task")

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation.
	if th.validator.IsEmpty(req.Title) {
		th.response.SendError(w, http.StatusBadRequest, "Task title is required")
		return
	}

	task, err := th.taskService.CreateTask(&req)
	if err != nil {
		th.logger.Error("Failed to create task: %v", err)
		th.response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	th.logger.Info("Created task with ID: %d", task.ID)
	th.response.SendCreated(w, task)
}

// UpdateTask handles PUT /tasks/{id} requests.
func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		th.response.SendError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	th.logger.Debug("Updating task with ID: %d", id)

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	task, err := th.taskService.UpdateTask(id, &req)
	if err != nil {
		th.logger.Error("Failed to update task %d: %v", id, err)
		th.response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	th.logger.Info("Updated task with ID: %d", task.ID)
	th.response.SendSuccess(w, task)
}

// DeleteTask handles DELETE /tasks/{id} requests.
func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		th.response.SendError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	th.logger.Debug("Deleting task with ID: %d", id)

	if err := th.taskService.DeleteTask(id); err != nil {
		th.logger.Error("Failed to delete task %d: %v", id, err)
		th.response.SendError(w, http.StatusNotFound, "Task not found")
		return
	}

	th.logger.Info("Deleted task with ID: %d", id)
	th.response.SendNoContent(w)
}

// SearchTasks handles POST /tasks/search requests.
func (th *TaskHandler) SearchTasks(w http.ResponseWriter, r *http.Request) {
	th.logger.Debug("Searching tasks")

	var query models.TaskSearchQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		th.response.SendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	tasks, err := th.taskService.SearchTasks(&query)
	if err != nil {
		th.logger.Error("Failed to search tasks: %v", err)
		th.response.SendError(w, http.StatusInternalServerError, "Failed to search tasks")
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
		"query": query.Query,
	}

	th.response.SendSuccess(w, response)
}

// GetTaskStats handles GET /tasks/stats requests.
func (th *TaskHandler) GetTaskStats(w http.ResponseWriter, r *http.Request) {
	th.logger.Debug("Getting task statistics")

	stats := th.taskService.GetTaskStats()
	th.response.SendSuccess(w, stats)
}
