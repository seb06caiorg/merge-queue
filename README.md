# 🚀 Task Manager API - Hackathon Project

A lightweight but feature-rich task management API built in Go, perfect for demonstrating both human and AI collaboration during hackathons.

## 🌟 Features

- **RESTful API** with full CRUD operations for tasks
- **Thread-safe** task management with goroutine safety
- **JSON configuration** support for easy customization
- **Middleware** for CORS and request logging
- **Validation** helpers and error handling
- **Sample data** pre-loaded for immediate testing
- **Built-in web interface** for API exploration

## 🚦 Quick Start

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the server:**
   ```bash
   go run .
   ```

3. **Test the API:**
   ```bash
   curl http://localhost:8080/api/v1/tasks
   ```

4. **View the web interface:**
   Open http://localhost:8080 in your browser

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/tasks` | Get all tasks (supports `?status=pending` filter) |
| POST | `/api/v1/tasks` | Create a new task |
| GET | `/api/v1/tasks/{id}` | Get specific task |
| PUT | `/api/v1/tasks/{id}` | Update task |
| DELETE | `/api/v1/tasks/{id}` | Delete task |

## 💡 Perfect for Hackathon Collaboration

### Areas for Human Enhancement:
- Add new task fields (tags, due dates, attachments)
- Implement user authentication and permissions
- Add task search and filtering capabilities
- Create batch operations for multiple tasks
- Add task comments and activity tracking

### Areas for AI Enhancement:
- Implement smart task prioritization algorithms
- Add automated task assignment based on workload
- Create task analytics and reporting features
- Implement notification systems
- Add data persistence (database integration)

## 🛠 Project Structure

```
├── main.go          # Main application and HTTP handlers
├── models.go        # Data structures and validation
├── utils.go         # Utility functions and helpers
├── config.json      # Application configuration
├── go.mod           # Go module definition
└── README.md        # This file
```

## 🔧 Configuration

Edit `config.json` to customize:
- Server port and host
- Feature toggles (CORS, logging)
- Default values for tasks
- Application metadata

## 📊 Sample Data

The API comes pre-loaded with sample tasks to demonstrate functionality:
- Project setup task (completed)
- API implementation (in-progress)
- Authentication feature (pending)
- Documentation task (pending)

## 🎯 Hackathon Ideas

1. **Frontend Integration**: Build a React/Vue frontend
2. **Database Layer**: Add PostgreSQL/MongoDB support
3. **Real-time Updates**: Implement WebSocket notifications
4. **Team Features**: Add user management and team collaboration
5. **Analytics Dashboard**: Create task metrics and visualizations
6. **Mobile API**: Extend for mobile app integration
7. **Automation**: Add task scheduling and automation rules

---

*Built for merge-queue testing and hackathon collaboration! 🎉*

test22

this is a fork.