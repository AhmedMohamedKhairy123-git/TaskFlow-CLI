package task

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

func (s *TaskStore) ExportToFile(format ExportFormat, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	switch format {
	case FormatJSON:
		return s.exportJSON(file)
	case FormatCSV:
		return s.exportCSV(file)
	case FormatText:
		return s.exportText(file)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (s *TaskStore) exportJSON(file *os.File) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.GetAll())
}

func (s *TaskStore) exportCSV(file *os.File) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	writer.Write([]string{"ID", "Title", "Completed", "Priority", "Tags", "CreatedAt"})
	
	for _, task := range s.GetAll() {
		writer.Write([]string{
			fmt.Sprint(task.ID),
			task.Title,
			fmt.Sprint(task.Completed),
			task.Priority.String(),
			strings.Join(task.Tags, "|"),
			task.CreatedAt.Format(time.RFC3339),
		})
	}
	
	return nil
}

func (s *TaskStore) exportText(file *os.File) error {
	for _, task := range s.GetAll() {
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}
		line := fmt.Sprintf("%s %d: %s (Priority: %s, Tags: %v)\n",
			status, task.ID, task.Title, task.Priority, task.Tags)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func (s *TaskStore) ImportFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 2 {
			s.Add(record[1])
		}
	}
	
	return nil
}