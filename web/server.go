package web

import (
	"fmt"
	"net/http"
	"task-tracker/task"
	"time"
)

type Server struct {
	store     *task.TaskStore
	httpServer *http.Server
}

func NewServer(store *task.TaskStore) *Server {
	return &Server{store: store}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.homeHandler)
	mux.HandleFunc("/health", s.healthHandler)

	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Task Tracker API</h1>")
	fmt.Fprintf(w, "<p>Available endpoints:</p>")
	fmt.Fprintf(w, "<ul>")
	fmt.Fprintf(w, "<li><a href='/health'>/health</a> - Health check</li>")
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