package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// ErrorCode represents specific error types
type ErrorCode int

const (
	ErrNotFound ErrorCode = iota + 1000
	ErrInvalidInput
	ErrDuplicateTask
	ErrStorageFailure
	ErrPermissionDenied
	ErrValidationFailed
	ErrConcurrentModification
	ErrBackupFailed
	ErrRecoveryFailed
)

func (e ErrorCode) String() string {
	switch e {
	case ErrNotFound:
		return "NOT_FOUND"
	case ErrInvalidInput:
		return "INVALID_INPUT"
	case ErrDuplicateTask:
		return "DUPLICATE_TASK"
	case ErrStorageFailure:
		return "STORAGE_FAILURE"
	case ErrPermissionDenied:
		return "PERMISSION_DENIED"
	case ErrValidationFailed:
		return "VALIDATION_FAILED"
	case ErrConcurrentModification:
		return "CONCURRENT_MODIFICATION"
	case ErrBackupFailed:
		return "BACKUP_FAILED"
	case ErrRecoveryFailed:
		return "RECOVERY_FAILED"
	default:
		return "UNKNOWN_ERROR"
	}
}

// AppError is a custom error type with rich information
type AppError struct {
	Code      ErrorCode    `json:"code"`
	Message   string       `json:"message"`
	Operation string       `json:"operation"`
	Timestamp time.Time    `json:"timestamp"`
	Stack     string       `json:"-"`
	Prev      error        `json:"-"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

func NewAppError(code ErrorCode, operation, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Operation: operation,
		Timestamp: time.Now(),
		Stack:     getStack(),
		Context:   make(map[string]interface{}),
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s (operation: %s)", 
		e.Code, e.Timestamp.Format(time.RFC3339), e.Message, e.Operation)
}

func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

func (e *AppError) WithCause(err error) *AppError {
	e.Prev = err
	return e
}

func (e *AppError) Unwrap() error {
	return e.Prev
}

func (e *AppError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// Recovery represents a panic recovery handler
type Recovery struct {
	ID        string
	Timestamp time.Time
	PanicValue interface{}
	Stack     string
	Handled   bool
}

func NewRecovery(r interface{}) *Recovery {
	return &Recovery{
		ID:         fmt.Sprintf("recovery-%d", time.Now().UnixNano()),
		Timestamp:  time.Now(),
		PanicValue: r,
		Stack:      getStack(),
		Handled:    false,
	}
}

func getStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// ErrorHandler function type
type ErrorHandler func(error) error

// ErrorMiddleware creates a chain of error handlers
func ErrorMiddleware(err error, handlers ...ErrorHandler) error {
	current := err
	for _, handler := range handlers {
		if current == nil {
			return nil
		}
		current = handler(current)
	}
	return current
}

// Logging handler
func LoggingHandler(next ErrorHandler) ErrorHandler {
	return func(err error) error {
		fmt.Printf("🔴 ERROR: %v\n", err)
		if appErr, ok := err.(*AppError); ok {
			fmt.Printf("   Stack: %s\n", appErr.Stack)
		}
		return next(err)
	}
}

// Recovery handler
func RecoveryHandler(next ErrorHandler) ErrorHandler {
	return func(err error) error {
		defer func() {
			if r := recover(); r != nil {
				rec := NewRecovery(r)
				fmt.Printf("🔥 PANIC RECOVERED: %+v\n", rec)
			}
		}()
		return next(err)
	}
}

// MultiError handles multiple errors
type MultiError struct {
	Errors []error
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred", len(m.Errors))
}

func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}