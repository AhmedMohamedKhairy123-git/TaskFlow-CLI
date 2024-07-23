package task

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	TaskID    int       `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Attachment struct {
	ID          string    `json:"id"`
	TaskID      int       `json:"task_id"`
	Filename    string    `json:"filename"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func (s *TaskStore) AddNote(taskID int, content string) (Note, error) {
	task, exists := s.Tasks[taskID]
	if !exists {
		return Note{}, fmt.Errorf("task %d not found", taskID)
	}
	
	note := Note{
		ID:        fmt.Sprintf("note_%d_%d", taskID, time.Now().UnixNano()),
		TaskID:    taskID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	task.Notes = append(task.Notes, note)
	return note, nil
}

func (s *TaskStore) GetNotes(taskID int) []Note {
	task, exists := s.Tasks[taskID]
	if !exists {
		return nil
	}
	return task.Notes
}

func (s *TaskStore) AddAttachment(taskID int, filename string, data io.Reader) (*Attachment, error) {
	task, exists := s.Tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %d not found", taskID)
	}
	
	attachmentsDir := "attachments"
	os.MkdirAll(attachmentsDir, 0755)
	
	attachmentID := fmt.Sprintf("att_%d_%d", taskID, time.Now().UnixNano())
	ext := filepath.Ext(filename)
	attachmentPath := filepath.Join(attachmentsDir, attachmentID+ext)
	
	file, err := os.Create(attachmentPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	size, err := io.Copy(file, data)
	if err != nil {
		return nil, err
	}
	
	attachment := &Attachment{
		ID:         attachmentID,
		TaskID:     taskID,
		Filename:   filename,
		Path:       attachmentPath,
		Size:       size,
		UploadedAt: time.Now(),
	}
	
	task.Attachments = append(task.Attachments, attachment)
	return attachment, nil
}

func (s *TaskStore) GetAttachments(taskID int) []Attachment {
	task, exists := s.Tasks[taskID]
	if !exists {
		return nil
	}
	return task.Attachments
}