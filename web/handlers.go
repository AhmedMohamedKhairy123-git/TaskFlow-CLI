package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"task-tracker/task"
)

func (s *Server) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := s.store.GetAll()
	response := make([]TaskResponse, len(tasks))
	
	for i, t := range tasks {
		response[i] = toTaskResponse(t)
	}
	
	respondJSON(w, http.StatusOK, response)
}

func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}
	
	newTask, err := s.store.Add(req.Title)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}
	
	// Set priority if provided
	if req.Priority != "" {
		priorityMap := map[string]task.Priority{
			"low":      task.Low,
			"medium":   task.Medium,
			"high":     task.High,
			"critical": task.Critical,
		}
		if p, ok := priorityMap[strings.ToLower(req.Priority)]; ok {
			s.store.SetPriority(newTask.ID, p)
		}
	}
	
	// Add tags if provided
	for _, tag := range req.Tags {
		s.store.AddTag(newTask.ID, tag)
	}
	
	task, _ := s.store.Get(newTask.ID)
	respondJSON(w, http.StatusCreated, toTaskResponse(task))
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	
	t, err := s.store.Get(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	
	respondJSON(w, http.StatusOK, toTaskResponse(t))
}

func (s *Server) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	task, err := s.store.Get(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	
	if req.Title != "" {
		task.Title = req.Title
	}
	
	s.store.Update(id, task)
	respondJSON(w, http.StatusOK, toTaskResponse(task))
}

func (s *Server) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	
	if err := s.store.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int, error) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	return strconv.Atoi(path)
}