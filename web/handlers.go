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