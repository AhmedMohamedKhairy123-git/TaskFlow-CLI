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
	
	// Home page
	mux.HandleFunc("/", s.homeHandler)
	
	// Health check
	mux.HandleFunc("/health", s.healthHandler)
	
	// Tasks collection
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.getTasksHandler(w, r)
		case http.MethodPost:
			s.createTaskHandler(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
	
	// Single task
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
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
	})
	
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	fmt.Printf("🚀 Web server starting on port %s\n", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Task Tracker API</h1>")
	fmt.Fprintf(w, "<p>Available endpoints:</p>")
	fmt.Fprintf(w, "<ul>")
	fmt.Fprintf(w, "<li>GET /health - Health check</li>")
	fmt.Fprintf(w, "<li>GET /tasks - List all tasks</li>")
	fmt.Fprintf(w, "<li>POST /tasks - Create new task</li>")
	fmt.Fprintf(w, "<li>GET /tasks/{id} - Get specific task</li>")
	fmt.Fprintf(w, "<li>PUT /tasks/{id} - Update task</li>")
	fmt.Fprintf(w, "<li>DELETE /tasks/{id} - Delete task</li>")
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