package handlers

import (
	"net/http"

	"merge-queue/internal/config"
	"merge-queue/pkg/utils"
)

// StaticHandler handles static content and web interface.
type StaticHandler struct {
	config *config.Config
	logger *utils.Logger
}

// NewStaticHandler creates a new StaticHandler instance.
func NewStaticHandler(cfg *config.Config, logger *utils.Logger) *StaticHandler {
	return &StaticHandler{
		config: cfg,
		logger: logger,
	}
}

// ServeHome handles GET / requests with a simple web interface.
func (sh *StaticHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	sh.logger.Debug("Serving home page")

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + sh.config.App.Name + `</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }

        .header {
            text-align: center;
            color: white;
            margin-bottom: 3rem;
        }

        .header h1 {
            font-size: 3rem;
            margin-bottom: 0.5rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .header p {
            font-size: 1.2rem;
            opacity: 0.9;
        }

        .card {
            background: white;
            border-radius: 12px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            backdrop-filter: blur(10px);
        }

        .endpoints {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        .endpoint {
            background: #f8f9fa;
            padding: 1.5rem;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }

        .endpoint h3 {
            color: #667eea;
            margin-bottom: 0.5rem;
            font-size: 1.1rem;
        }

        .endpoint p {
            color: #666;
            font-size: 0.9rem;
        }

        .method {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 4px;
            font-size: 0.8rem;
            font-weight: bold;
            margin-right: 0.5rem;
        }

        .method.get { background: #d4edda; color: #155724; }
        .method.post { background: #cce5ff; color: #004085; }
        .method.put { background: #fff3cd; color: #856404; }
        .method.delete { background: #f8d7da; color: #721c24; }

        .quick-test {
            background: #e8f5e8;
            padding: 1.5rem;
            border-radius: 8px;
            border-left: 4px solid #28a745;
        }

        .quick-test h3 {
            color: #28a745;
            margin-bottom: 1rem;
        }

        .code {
            background: #2d3748;
            color: #e2e8f0;
            padding: 1rem;
            border-radius: 6px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 0.9rem;
            overflow-x: auto;
            margin: 0.5rem 0;
        }

        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1rem;
        }

        .feature {
            text-align: center;
            padding: 1rem;
        }

        .feature-icon {
            font-size: 2rem;
            margin-bottom: 0.5rem;
        }

        .stats {
            display: flex;
            justify-content: space-around;
            text-align: center;
            margin: 2rem 0;
        }

        .stat {
            color: white;
        }

        .stat-number {
            font-size: 2rem;
            font-weight: bold;
            display: block;
        }

        .stat-label {
            opacity: 0.8;
            font-size: 0.9rem;
        }

        @media (max-width: 768px) {
            .container { padding: 1rem; }
            .header h1 { font-size: 2rem; }
            .endpoints { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 ` + sh.config.App.Name + `</h1>
            <p>Version ` + sh.config.App.Version + ` • Built for Hackathon Excellence</p>

            <div class="stats">
                <div class="stat">
                    <span class="stat-number">6</span>
                    <span class="stat-label">API Endpoints</span>
                </div>
                <div class="stat">
                    <span class="stat-number">4</span>
                    <span class="stat-label">Sample Tasks</span>
                </div>
                <div class="stat">
                    <span class="stat-number">100%</span>
                    <span class="stat-label">Ready to Hack</span>
                </div>
            </div>
        </div>

        <div class="card">
            <h2>🌟 Features</h2>
            <div class="features">
                <div class="feature">
                    <div class="feature-icon">⚡</div>
                    <h4>Lightning Fast</h4>
                    <p>Built with Go for maximum performance</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🔒</div>
                    <h4>Thread Safe</h4>
                    <p>Concurrent operations with mutex protection</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🎯</div>
                    <h4>RESTful API</h4>
                    <p>Clean, intuitive endpoints</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🛠️</div>
                    <h4>Configurable</h4>
                    <p>JSON configuration with environment overrides</p>
                </div>
            </div>
        </div>

        <div class="card">
            <h2>📋 API Endpoints</h2>
            <div class="endpoints">
                <div class="endpoint">
                    <h3><span class="method get">GET</span>/api/v1/health</h3>
                    <p>Health check endpoint for monitoring</p>
                </div>
                <div class="endpoint">
                    <h3><span class="method get">GET</span>/api/v1/tasks</h3>
                    <p>Get all tasks with optional filtering (?status=pending)</p>
                </div>
                <div class="endpoint">
                    <h3><span class="method post">POST</span>/api/v1/tasks</h3>
                    <p>Create a new task with title, description, etc.</p>
                </div>
                <div class="endpoint">
                    <h3><span class="method get">GET</span>/api/v1/tasks/{id}</h3>
                    <p>Get a specific task by ID</p>
                </div>
                <div class="endpoint">
                    <h3><span class="method put">PUT</span>/api/v1/tasks/{id}</h3>
                    <p>Update an existing task</p>
                </div>
                <div class="endpoint">
                    <h3><span class="method delete">DELETE</span>/api/v1/tasks/{id}</h3>
                    <p>Delete a task by ID</p>
                </div>
            </div>
        </div>

        <div class="quick-test">
            <h3>🧪 Quick Test Commands</h3>
            <p>Try these commands in your terminal:</p>

            <div class="code">curl http://localhost` + sh.config.Server.Port + `/api/v1/health</div>
            <div class="code">curl http://localhost` + sh.config.Server.Port + `/api/v1/tasks</div>
            <div class="code">curl -X POST http://localhost` + sh.config.Server.Port + `/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Task","description":"Created from curl"}'</div>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
