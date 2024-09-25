# Task Tracker - Advanced Go CLI Application

A production-grade command-line task management system built with Go, demonstrating advanced concepts and best practices.

## 🏗️ Architecture & Project Structure

This project follows a modular design to ensure **Separation of Concerns (SoC)** and high maintainability, utilizing Go's interface-driven development.

```text
task-tracker/
├── backups/              # Directory for automated backup storage
├── config/
│   └── config.go         # Configuration management logic
├── errors/
│   └── errors.go         # Custom error types and panic recovery middleware
├── storage/
│   ├── storage.go        # Storage interface (Abstraction layer)
│   └── json_storage.go   # JSON implementation of the storage interface
├── task/
│   ├── task.go           # Core Task domain struct and methods
│   ├── store.go          # Task collection management
│   ├── validator.go      # Business logic validation framework
│   ├── backup.go         # Auto-backup system implementation
│   ├── priority.go       # Priority queue logic for task sorting
│   └── interfaces.go     # Shared interface definitions for decoupling
├── .gitignore            # Git exclusion rules
├── config.json           # Application runtime settings
├── go.mod                # Go module dependencies
├── main.go               # Application entry point
├── readme.md             # Project documentation
└── tasks.json            # Local data persistence layer
```
## 🚀 Features Implemented

### Phase 1-2: Core Go Features
- ✅ Package organization
- ✅ Variables, constants, structs
- ✅ Functions with multiple returns
- ✅ Control flow (if/else, switch, loops)
- ✅ User input handling

### Phase 3: Data Structures
- ✅ Slices for dynamic task lists
- ✅ Maps for O(1) task lookup
- ✅ Struct composition
- ✅ Range iterations

### Phase 4: Methods & Interfaces
- ✅ Pointer receivers vs value receivers
- ✅ Interface satisfaction (implicit)
- ✅ Interface composition
- ✅ Stringer interface implementation
- ✅ Custom types with methods
- ✅ Method chaining

### Phase 5: Advanced Error Handling
- ✅ Custom error types with stack traces
- ✅ Error wrapping and unwrapping
- ✅ Panic recovery with defer/recover
- ✅ Multi-error aggregation
- ✅ Error context enrichment
- ✅ Error middleware pattern
- ✅ Graceful degradation

### Phase 6: File I/O & Persistence
- ✅ JSON marshaling/unmarshaling
- ✅ Atomic file writes (temp + rename)
- ✅ Auto-save with ticker
- ✅ Backup rotation (keep last N)
- ✅ Corrupted file recovery
- ✅ Configuration management
- ✅ File metadata inspection

### Phase 7: Concurrency
- ✅ Goroutines for auto-saver
- ✅ Channels for communication
- ✅ Select with timeouts
- ✅ Producer-consumer pattern
- ✅ Buffered channels
- ✅ WaitGroup coordination
- ✅ Race condition prevention

# Phase 8 & 9: Web API + Advanced Features

## 🚀 Phase 8: Web API (6 Mini-Phases)

### 8.1 Basic HTTP Server
- Simple HTTP server with `/` and `/health` endpoints
- Graceful shutdown support
- Timeout configurations

### 8.2 API Models
- JSON request/response structs
- TaskResponse, CreateTaskRequest, ErrorResponse
- Helper functions for JSON responses

### 8.3 Task Handlers
- `GET /tasks` - List all tasks
- `POST /tasks` - Create new task
- JSON request body parsing

### 8.4 Single Task Operations
- `GET /tasks/{id}` - Get specific task
- `PUT /tasks/{id}` - Update task
- `DELETE /tasks/{id}` - Delete task
- Path parameter parsing

### 8.5 Middleware
- Logging middleware (request/response timing)
- Panic recovery middleware
- CORS middleware for cross-origin requests
- Rate limiting (10 concurrent requests)

### 8.6 Integration
- Web server starts alongside CLI
- Graceful shutdown with context timeout
- Background goroutine for server
- Signal handling (Ctrl+C)

## 🔥 Phase 9: Advanced Features (10 Mini-Phases)

### 9.1 Search & Filter

```go
criteria := SearchCriteria{
    Title:     "meeting",
    Completed: false,
    Priority:  []Priority{High, Critical},
    Tags:      []string{"work"},
}
results := store.Search(criteria)
```
## 9.2 Analytics & Statistics
- Completion rates

- Priority distribution

- Tag frequency analysis

- Average tasks per day

- Oldest/newest task tracking

## 9.3 Export/Import
- JSON export with pretty print

- CSV export for spreadsheet apps

- Text format for readability

- Import from CSV

## 9.4 Reminders & Due Dates
- Set due dates for tasks

- Repeat options (daily, weekly, monthly)

+ Automatic reminder checking

+ Overdue task detection

## 9.5 Task Dependencies
- Blocked/depends-on relationships

- Circular dependency detection

- Get blocked tasks list

- Find dependent tasks

## 9.6 History & Audit Log
- Track all task actions

- Before/after change tracking

- Timestamp for each action

- Export history to JSON

## 9.7 Bulk Operations
- Add multiple tasks at once

- Complete multiple tasks

- Delete multiple tasks

- Parallel processing with goroutines

## 9.8 Task Templates
- Predefined task templates

- Meeting template with subtasks

-  Project template with checklist

- Save/load templates from file

## 9.9 Notes & Attachments
- Add text notes to tasks

- File attachments support

- Multiple notes per task

- Attachment metadata tracking

## 9.10 Sharing & Collaboration
-  Generate share links

- Permission levels (view/edit/admin)

- Expiring shares

- Revoke sharing access

## 🌐 API Endpoints Summary
```text
GET    /health           # Health check
GET    /tasks            # List all tasks
POST   /tasks            # Create new task
GET    /tasks/{id}       # Get specific task
PUT    /tasks/{id}       # Update task
DELETE /tasks/{id}       # Delete task
📦 New Packages Added
```
## 🧠 Design Patterns Used

| Pattern | Implementation |
|---------|---------------|
| Repository | `Storage` interface |
| Factory | `NewTaskStore()`, `NewValidator()` |
| Strategy | Validation rules |
| Observer | Auto-save ticker |
| Chain of Responsibility | Error middleware |
| Decorator | Error handlers |
| Adapter | JSON storage adapter |
| Singleton | Config loader |
| Builder | Error context building |
| Composite | MultiError |
| Producer-Consumer | Auto-save channel |

## 🔧 Advanced Go Concepts

### Interface Satisfaction (Implicit)
```go
type FullTaskManager interface {
    TaskManager
    Completer
    Prioritizer
    Tagger
    StatsReporter
}
// *TaskStore automatically satisfies FullTaskManager
Error Handling with Stack Traces
go
func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", 
        e.Code, e.Timestamp, e.Message)
}
Atomic File Operations
go
// Write to temp file then rename
tempFile := file + ".tmp"
os.WriteFile(tempFile, data, 0644)
os.Rename(tempFile, file) // Atomic on most OS
Graceful Panic Recovery
go
defer func() {
    if r := recover(); r != nil {
        log.Printf("Recovered: %v", r)
        emergencySave()
    }
}()
Producer-Consumer with Channels
go
saveChan := make(chan bool, 1)
go producer(ticker, saveChan)
go consumer(saveChan, storage)
```
## 📊 Performance Optimizations
```
O(1) task lookup using maps

Buffered channels to prevent blocking

Atomic file writes to prevent corruption

Lazy loading of tasks

Goroutine pooling for backups

Batch operations for bulk updates

Connection pooling (conceptual)

Query optimization (conceptual)
```

## 🔒 Security Features
```Input validation

Profanity filtering

Path traversal prevention

Safe file permissions (0644)

Panic recovery (no crashes)

Error sanitization (no sensitive data)

Rate limiting (conceptual)

Authentication (conceptual)
```
## 🧪 Testing Strategy
```go
// Example test pattern
func TestTaskStore(t *testing.T) {
    store := NewTaskStore()
    task, _ := store.Add("Test")
    
    assert.Equal(t, "Test", task.Title)
    assert.False(t, task.Completed)
}
🐳 Docker Support (Conceptual)
dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o task-tracker
CMD ["./task-tracker"]
```
## ☸️ Kubernetes (Conceptual)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-tracker
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-tracker
📈 Monitoring (Conceptual)
go
// Prometheus metrics
var (
    taskCreated = prometheus.NewCounter(...)
    taskDuration = prometheus.NewHistogram(...)
)
```
## 🚦 Production Readiness
```
✅ Graceful shutdown

✅ Signal handling

✅ Configuration management

✅ Logging

✅ Error tracking

✅ Backup/Restore

✅ Data validation

✅ Race condition free

✅ Memory safe

✅ Goroutine leak free
```