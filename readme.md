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
