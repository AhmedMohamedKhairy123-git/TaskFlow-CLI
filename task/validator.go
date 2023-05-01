package task

import (
	"fmt"  // Add this import
	"regexp"
	"strings"
	"task-tracker/errors"
	"time"
	"unicode"
)

type ValidationRule func(*Task) error

type Validator struct {
	rules []ValidationRule
}

func NewValidator() *Validator {
	return &Validator{
		rules: []ValidationRule{
			validateTitle,
			validateTags,
			validatePriority,
			validateCreatedAt,
		},
	}
}

func (v *Validator) Validate(task *Task) error {
	var multiErr errors.MultiError
	
	for _, rule := range v.rules {
		if err := rule(task); err != nil {
			multiErr.Add(err)
		}
	}
	
	if multiErr.HasErrors() {
		return &multiErr
	}
	return nil
}

func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

func validateTitle(task *Task) error {
	if task.Title == "" {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateTitle",
			"task title cannot be empty",
		)
	}
	
	if len(task.Title) < 3 {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateTitle",
			"task title must be at least 3 characters",
		).WithContext("title", task.Title)
	}
	
	if len(task.Title) > 100 {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateTitle",
			"task title cannot exceed 100 characters",
		).WithContext("length", len(task.Title))
	}
	
	// Check for valid characters
	for _, r := range task.Title {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return errors.NewAppError(
				errors.ErrValidationFailed,
				"validateTitle",
				"task title contains invalid characters",
			)
		}
	}
	
	return nil
}

func validateTags(task *Task) error {
	if len(task.Tags) > 10 {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateTags",
			"too many tags (max 10)",
		).WithContext("tag_count", len(task.Tags))
	}
	
	tagRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	
	for _, tag := range task.Tags {
		if len(tag) < 2 {
			return errors.NewAppError(
				errors.ErrValidationFailed,
				"validateTags",
				"tag must be at least 2 characters",
			).WithContext("tag", tag)
		}
		
		if len(tag) > 20 {
			return errors.NewAppError(
				errors.ErrValidationFailed,
				"validateTags",
				"tag cannot exceed 20 characters",
			).WithContext("tag", tag)
		}
		
		if !tagRegex.MatchString(tag) {
			return errors.NewAppError(
				errors.ErrValidationFailed,
				"validateTags",
				"tag contains invalid characters (use letters, numbers, underscore, hyphen)",
			).WithContext("tag", tag)
		}
	}
	
	return nil
}

func validatePriority(task *Task) error {
	if task.Priority < Low || task.Priority > Critical {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validatePriority",
			"invalid priority value",
		).WithContext("priority", task.Priority)
	}
	return nil
}

func validateCreatedAt(task *Task) error {
	if task.CreatedAt.IsZero() {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateCreatedAt",
			"created date cannot be zero",
		)
	}
	
	if task.CreatedAt.After(time.Now().Add(time.Hour)) {
		return errors.NewAppError(
			errors.ErrValidationFailed,
			"validateCreatedAt",
			"created date cannot be in the future",
		).WithContext("created_at", task.CreatedAt)
	}
	
	return nil
}

// Custom validators
func MinTitleLength(min int) ValidationRule {
	return func(task *Task) error {
		if len(task.Title) < min {
			return errors.NewAppError(
				errors.ErrValidationFailed,
				"MinTitleLength",
				fmt.Sprintf("title must be at least %d characters", min),
			)
		}
		return nil
	}
}

func NoProfanity(profanityList []string) ValidationRule {
	return func(task *Task) error {
		titleLower := strings.ToLower(task.Title)
		for _, word := range profanityList {
			if strings.Contains(titleLower, word) {
				return errors.NewAppError(
					errors.ErrValidationFailed,
					"NoProfanity",
					"title contains prohibited words",
				)
			}
		}
		return nil
	}
}