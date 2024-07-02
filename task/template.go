package task

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type TaskTemplate struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Tags        []string   `json:"tags"`
	Subtasks    []string   `json:"subtasks"`
}

type TemplateStore struct {
	templates map[string]TaskTemplate
}

func NewTemplateStore() *TemplateStore {
	return &TemplateStore{
		templates: make(map[string]TaskTemplate),
	}
}

func (ts *TemplateStore) AddTemplate(template TaskTemplate) {
	ts.templates[template.Name] = template
}

func (ts *TemplateStore) GetTemplate(name string) (TaskTemplate, bool) {
	t, ok := ts.templates[name]
	return t, ok
}

func (ts *TemplateStore) ListTemplates() []TaskTemplate {
	var list []TaskTemplate
	for _, t := range ts.templates {
		list = append(list, t)
	}
	return list
}

func (s *TaskStore) CreateFromTemplate(template TaskTemplate, titleSuffix string) ([]Task, error) {
	var created []Task
	
	mainTitle := template.Name
	if titleSuffix != "" {
		mainTitle = mainTitle + " - " + titleSuffix
	}
	
	mainTask, err := s.Add(mainTitle)
	if err != nil {
		return nil, err
	}
	
	mainTask.Priority = template.Priority
	mainTask.Tags = append(mainTask.Tags, template.Tags...)
	
	created = append(created, *mainTask)
	
	for _, subtaskTitle := range template.Subtasks {
		subtask, err := s.Add(subtaskTitle)
		if err != nil {
			continue
		}
		subtask.Priority = template.Priority
		s.AddDependency(subtask.ID, mainTask.ID, DependsOn)
		created = append(created, *subtask)
	}
	
	return created, nil
}

func (ts *TemplateStore) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(ts.templates, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (ts *TemplateStore) LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &ts.templates)
}

func DefaultTemplates() *TemplateStore {
	ts := NewTemplateStore()
	
	ts.AddTemplate(TaskTemplate{
		Name:        "Meeting",
		Description: "Schedule and prepare for meeting",
		Priority:    Medium,
		Tags:        []string{"work", "meeting"},
		Subtasks:    []string{"Prepare agenda", "Send invites", "Take notes"},
	})
	
	ts.AddTemplate(TaskTemplate{
		Name:        "Project",
		Description: "Start a new project",
		Priority:    High,
		Tags:        []string{"work", "project"},
		Subtasks:    []string{"Define scope", "Create timeline", "Assign roles"},
	})
	
	return ts
}