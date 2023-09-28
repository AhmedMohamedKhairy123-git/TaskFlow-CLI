package web

import (
	"encoding/json"
	"net/http"
	"task-tracker/task"
	"time"
)

type TaskResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Priority  string    `json:"priority"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTaskRequest struct {
	Title    string   `json:"title"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func toTaskResponse(t task.Task) TaskResponse {
	return TaskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.Completed,
		Priority:  t.Priority.String(),
		Tags:      t.Tags,
		CreatedAt: t.CreatedAt,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Error:   http.StatusText(status),
		Code:    status,
		Message: message,
	}
	respondJSON(w, status, resp)
}