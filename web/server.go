package web

import (
	"fmt"
	"net/http"
	"task-tracker/task"
	"time"
)

type Server struct {
	store      *task.TaskStore
	httpServer *http.Server
}

func NewServer(store *task.TaskStore) *Server {
	return &Server{store: store}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// 1. Home Route
	mux.HandleFunc("/", s.loggingMiddleware(
		s.recoveryMiddleware(
			s.corsMiddleware(
				s.rateLimitMiddleware(s.homeHandler)))))

	// 2. Health Check Route
	mux.HandleFunc("/health", s.loggingMiddleware(
		s.recoveryMiddleware(s.healthHandler)))

	// 3. Global Tasks Route (GET all / POST new)
	mux.HandleFunc("/tasks", s.loggingMiddleware(
		s.recoveryMiddleware(
			s.corsMiddleware(
				s.rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						s.getTasksHandler(w, r)
					case http.MethodPost:
						s.createTaskHandler(w, r)
					default:
						respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
					}
				})))))

	// 4. Specific Task Route (Integrated Version with Middleware)
	// This handles /tasks/{id}
	mux.HandleFunc("/tasks/", s.loggingMiddleware(
		s.recoveryMiddleware(
			s.corsMiddleware(
				s.rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
					// Method Switching Logic
					switch r.Method {
					case http.MethodGet:
						s.getTaskHandler(w, r)
					case http.MethodPut:
						s.updateTaskHandler(w, r)
					case http.MethodDelete:
						s.deleteTaskHandler(w, r)
					default:
						respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
					}
				})))))

	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	return s.httpServer.ListenAndServe()
}

// --- Handlers ---

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Task Tracker API</h1>")
	fmt.Fprintf(w, "<p>Available endpoints:</p>")
	fmt.Fprintf(w, "<ul>")
	fmt.Fprintf(w, "<li><a href='/health'>/health</a> - Health check</li>")
	fmt.Fprintf(w, "<li><a href='/tasks'>/tasks</a> - Get all tasks (GET) or Create task (POST)</li>")
	fmt.Fprintf(w, "<li>/tasks/{id} - Get, Update, Delete specific task</li>")
	fmt.Fprintf(w, "</ul>")
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "time": "%s"}`, time.Now().Format(time.RFC3339))
}

func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}