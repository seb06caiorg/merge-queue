package services

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"merge-queue/internal/models"
	"merge-queue/pkg/utils"
)

// TaskService handles business logic for task operations.
type TaskService struct {
	tasks     map[int]*models.Task
	nextID    int
	mutex     sync.RWMutex
	validator *utils.ValidationUtils
	timeUtils *utils.TimeUtils
	maxTasks  int
}

// NewTaskService creates a new TaskService instance.
func NewTaskService(maxTasks int) *TaskService {
	service := &TaskService{
		tasks:     make(map[int]*models.Task),
		nextID:    1,
		validator: utils.NewValidationUtils(),
		timeUtils: utils.NewTimeUtils(),
		maxTasks:  maxTasks,
	}

	// Add sample data for demonstration.
	service.addSampleTasks()

	return service
}

// CreateTask creates a new task.
func (ts *TaskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Validate request.
	if err := ts.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check task limit.
	if len(ts.tasks) >= ts.maxTasks {
		return nil, fmt.Errorf("maximum number of tasks (%d) reached", ts.maxTasks)
	}

	// Set defaults.
	status := req.Status
	if status == "" {
		status = "pending"
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	// Create task.
	task := &models.Task{
		ID:          ts.nextID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AssignedTo:  strings.TrimSpace(req.AssignedTo),
		Tags:        req.Tags,
	}

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task, nil
}

// GetTask retrieves a task by ID.
func (ts *TaskService) GetTask(id int) (*models.Task, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	return task, nil
}

// GetAllTasks returns all tasks with optional filtering.
func (ts *TaskService) GetAllTasks(filter *models.TaskFilter) ([]*models.Task, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	var tasks []*models.Task

	for _, task := range ts.tasks {
		if ts.matchesFilter(task, filter) {
			tasks = append(tasks, task)
		}
	}

	// Apply sorting.
	ts.sortTasks(tasks)

	// Apply pagination.
	if filter != nil && (filter.Limit > 0 || filter.Offset > 0) {
		tasks = ts.applyPagination(tasks, filter.Limit, filter.Offset)
	}

	return tasks, nil
}

// UpdateTask updates an existing task.
func (ts *TaskService) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	// Validate update request.
	if err := ts.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Apply updates.
	if req.Title != nil {
		task.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		task.Description = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.AssignedTo != nil {
		task.AssignedTo = strings.TrimSpace(*req.AssignedTo)
	}
	if req.Tags != nil {
		task.Tags = req.Tags
	}

	task.UpdatedAt = time.Now()

	return task, nil
}

// DeleteTask removes a task by ID.
func (ts *TaskService) DeleteTask(id int) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// SearchTasks searches for tasks based on query.
func (ts *TaskService) SearchTasks(query *models.TaskSearchQuery) ([]*models.Task, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	var results []*models.Task
	searchTerm := strings.ToLower(strings.TrimSpace(query.Query))

	for _, task := range ts.tasks {
		// Check if task matches filter criteria.
		if !ts.matchesFilter(task, &query.Filters) {
			continue
		}

		// Check if task matches search query.
		if ts.matchesSearchQuery(task, searchTerm, query.Fields) {
			results = append(results, task)
		}
	}

	// Apply sorting.
	ts.sortTasksBy(results, query.SortBy, query.SortDesc)

	return results, nil
}

// GetTaskStats returns statistics about tasks.
func (ts *TaskService) GetTaskStats() *models.TaskStats {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	stats := &models.TaskStats{
		TotalTasks:      len(ts.tasks),
		TasksByStatus:   make(map[string]int),
		TasksByPriority: make(map[string]int),
		TasksByUser:     make(map[string]int),
		LastUpdated:     time.Now(),
	}

	for _, task := range ts.tasks {
		stats.TasksByStatus[task.Status]++
		stats.TasksByPriority[task.Priority]++
		if task.AssignedTo != "" {
			stats.TasksByUser[task.AssignedTo]++
		}
	}

	return stats
}

// Helper methods.

func (ts *TaskService) validateCreateRequest(req *models.CreateTaskRequest) error {
	if err := ts.validator.ValidateRequired("title", req.Title); err != nil {
		return err
	}

	if err := ts.validator.ValidateLength("title", req.Title, 1, 200); err != nil {
		return err
	}

	if req.Description != "" {
		if err := ts.validator.ValidateLength("description", req.Description, 0, 1000); err != nil {
			return err
		}
	}

	if req.Status != "" && !models.IsValidStatus(req.Status) {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	if req.Priority != "" && !models.IsValidPriority(req.Priority) {
		return fmt.Errorf("invalid priority: %s", req.Priority)
	}

	if err := ts.validator.ValidateTagList(req.Tags, 10, 50); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) validateUpdateRequest(req *models.UpdateTaskRequest) error {
	if req.Title != nil {
		if err := ts.validator.ValidateRequired("title", *req.Title); err != nil {
			return err
		}
		if err := ts.validator.ValidateLength("title", *req.Title, 1, 200); err != nil {
			return err
		}
	}

	if req.Description != nil {
		if err := ts.validator.ValidateLength("description", *req.Description, 0, 1000); err != nil {
			return err
		}
	}

	if req.Status != nil && !models.IsValidStatus(*req.Status) {
		return fmt.Errorf("invalid status: %s", *req.Status)
	}

	if req.Priority != nil && !models.IsValidPriority(*req.Priority) {
		return fmt.Errorf("invalid priority: %s", *req.Priority)
	}

	if err := ts.validator.ValidateTagList(req.Tags, 10, 50); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) matchesFilter(task *models.Task, filter *models.TaskFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Status != "" && task.Status != filter.Status {
		return false
	}

	if filter.Priority != "" && task.Priority != filter.Priority {
		return false
	}

	if filter.AssignedTo != "" && task.AssignedTo != filter.AssignedTo {
		return false
	}

	if len(filter.Tags) > 0 {
		hasTag := false
		for _, filterTag := range filter.Tags {
			for _, taskTag := range task.Tags {
				if taskTag == filterTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	return true
}

func (ts *TaskService) matchesSearchQuery(task *models.Task, searchTerm string, fields []string) bool {
	if searchTerm == "" {
		return true
	}

	// If no fields specified, search in title and description.
	if len(fields) == 0 {
		fields = []string{"title", "description"}
	}

	for _, field := range fields {
		var content string
		switch field {
		case "title":
			content = strings.ToLower(task.Title)
		case "description":
			content = strings.ToLower(task.Description)
		default:
			continue
		}

		if strings.Contains(content, searchTerm) {
			return true
		}
	}

	return false
}

func (ts *TaskService) sortTasks(tasks []*models.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
}

func (ts *TaskService) sortTasksBy(tasks []*models.Task, sortBy string, desc bool) {
	switch sortBy {
	case "created_at":
		sort.Slice(tasks, func(i, j int) bool {
			if desc {
				return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	case "updated_at":
		sort.Slice(tasks, func(i, j int) bool {
			if desc {
				return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
			}
			return tasks[i].UpdatedAt.Before(tasks[j].UpdatedAt)
		})
	case "priority":
		priorityOrder := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
		sort.Slice(tasks, func(i, j int) bool {
			pi, pj := priorityOrder[tasks[i].Priority], priorityOrder[tasks[j].Priority]
			if desc {
				return pi > pj
			}
			return pi < pj
		})
	default:
		ts.sortTasks(tasks) // Default sort by creation time.
	}
}

func (ts *TaskService) applyPagination(tasks []*models.Task, limit, offset int) []*models.Task {
	if offset >= len(tasks) {
		return []*models.Task{}
	}

	end := len(tasks)
	if limit > 0 && offset+limit < len(tasks) {
		end = offset + limit
	}

	return tasks[offset:end]
}

func (ts *TaskService) addSampleTasks() {
	sampleTasks := []*models.CreateTaskRequest{
		{
			Title:       "Setup project structure",
			Description: "Create basic Go project layout with proper package organization",
			Status:      "completed",
			Priority:    "high",
			AssignedTo:  "alice",
			Tags:        []string{"setup", "infrastructure"},
		},
		{
			Title:       "Implement API endpoints",
			Description: "Create REST API endpoints for task management with proper error handling",
			Status:      "in-progress",
			Priority:    "high",
			AssignedTo:  "bob",
			Tags:        []string{"api", "backend"},
		},
		{
			Title:       "Add authentication",
			Description: "Implement JWT-based authentication and authorization middleware",
			Status:      "pending",
			Priority:    "medium",
			AssignedTo:  "charlie",
			Tags:        []string{"auth", "security"},
		},
		{
			Title:       "Write documentation",
			Description: "Create comprehensive API documentation and user guides",
			Status:      "pending",
			Priority:    "low",
			Tags:        []string{"docs", "documentation"},
		},
	}

	for _, req := range sampleTasks {
		ts.CreateTask(req)
	}
}
