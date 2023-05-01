package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"  // Add this
	"task-tracker/errors"
	"task-tracker/task"
)

var (
	store         *task.TaskStore
	validator     *task.Validator
	backupManager *task.BackupManager
	reader        = bufio.NewReader(os.Stdin)
)

func main() {
	// Initialize with panic recovery
	defer func() {
		if r := recover(); r != nil {
			rec := errors.NewRecovery(r)
			fmt.Printf("\n🔥 CRITICAL: Application recovered from panic!\n")
			fmt.Printf("Recovery ID: %s\n", rec.ID)
			fmt.Printf("Please check logs and consider restoring from backup.\n")
			
			// Attempt auto-recovery
			if backupManager != nil {
				fmt.Println("Attempting to restore from latest backup...")
				if err := backupManager.RestoreLatest(); err != nil {
					fmt.Printf("Auto-recovery failed: %v\n", err)
				} else {
					fmt.Println("✅ Auto-recovery successful!")
				}
			}
		}
	}()
	
	// Initialize components
	initializeApp()
	
	// Main loop with error handling
	runApplication()
}

func initializeApp() {
	fmt.Println("🚀 Initializing Task Tracker with Advanced Error Handling")
	
	// Initialize store
	store = task.NewTaskStore()
	
	// Initialize validator with custom rules
	validator = task.NewValidator()
	validator.AddRule(task.MinTitleLength(3))
	validator.AddRule(task.NoProfanity([]string{"badword", "spam"}))
	
	// Initialize backup manager
	backupManager = task.NewBackupManager(store, "./backups")
	backupManager.SetErrorHandler(func(err error) {
		fmt.Printf("📦 Backup System: %v\n", err)
	})
	
	// Load sample data with error handling
	loadSampleData()
	
	// Create initial backup
	if _, err := backupManager.CreateBackup(); err != nil {
		handleError(err, "Failed to create initial backup")
	}
	
	fmt.Println("✅ Initialization complete!")
}

func loadSampleData() {
	// Use error handling for sample data creation
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ Warning: Sample data loading recovered from panic: %v\n", r)
		}
	}()
	
	sampleTasks := []string{
		"Learn Go error handling",
		"Implement custom errors",
		"Add panic recovery",
		"Create backup system",
	}
	
	for _, title := range sampleTasks {
		if _, err := store.Add(title); err != nil {
			handleError(err, "Failed to add sample task")
		}
	}
	
	// Set some priorities
	store.SetPriority(1, task.High)
	store.SetPriority(2, task.Critical)
	store.SetPriority(3, task.Medium)
	
	// Add tags
	store.AddTag(1, "learning")
	store.AddTag(2, "critical")
	store.AddTag(3, "backlog")
}

func runApplication() {
	for {
		showMenu()
		
		if !processChoiceSafe() {
			break
		}
	}
	
	// Create final backup before exit
	fmt.Println("\n📦 Creating final backup...")
	if _, err := backupManager.CreateBackup(); err != nil {
		handleError(err, "Failed to create final backup")
	}
	
	fmt.Println("👋 Goodbye!")
}

func showMenu() {
	fmt.Println("\n🔰 === TASK TRACKER (Advanced Error Handling) ===")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task Complete")
	fmt.Println("4. Set Priority")
	fmt.Println("5. Add Tag")
	fmt.Println("6. Validate Task")
	fmt.Println("7. Create Backup")
	fmt.Println("8. List Backups")
	fmt.Println("9. Restore from Backup")
	fmt.Println("10. Test Panic Recovery")
	fmt.Println("11. Exit")
	fmt.Print("Enter choice: ")
}

func processChoiceSafe() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ Recovered from panic in menu handler: %v\n", r)
		}
	}()
	
	choice := readInt()
	
	switch choice {
	case 1:
		addTaskWithValidation()
	case 2:
		listTasksWithErrorHandling()
	case 3:
		markCompleteWithErrorHandling()
	case 4:
		setPriorityWithErrorHandling()
	case 5:
		addTagWithErrorHandling()
	case 6:
		validateSpecificTask()
	case 7:
		createBackup()
	case 8:
		listBackups()
	case 9:
		restoreFromBackup()
	case 10:
		testPanicRecovery()
	case 11:
		return false
	default:
		fmt.Println("Invalid choice!")
	}
	return true
}

func handleError(err error, context string) {
	if err == nil {
		return
	}
	
	fmt.Printf("\n❌ ERROR: %s\n", context)
	
	// Type assertion to check if it's our custom error
	if appErr, ok := err.(*errors.AppError); ok {
		fmt.Printf("   Code: %s\n", appErr.Code)
		fmt.Printf("   Message: %s\n", appErr.Message)
		fmt.Printf("   Operation: %s\n", appErr.Operation)
		fmt.Printf("   Time: %s\n", appErr.Timestamp.Format(time.RFC3339))
		
		if len(appErr.Context) > 0 {
			fmt.Println("   Context:")
			for k, v := range appErr.Context {
				fmt.Printf("     %s: %v\n", k, v)
			}
		}
	} else {
		fmt.Printf("   %v\n", err)
	}
}

func addTaskWithValidation() {
	fmt.Print("Enter task title: ")
	title := readString()
	
	// Create temporary task for validation
	tempTask := &task.Task{Title: title}
	
	// Validate before adding
	if err := validator.Validate(tempTask); err != nil {
		handleError(err, "Task validation failed")
		return
	}
	
	// Add task
	newTask, err := store.Add(title)
	if err != nil {
		handleError(err, "Failed to add task")
		return
	}
	
	fmt.Printf("✅ Task added with ID: %d\n", newTask.ID)
}

func listTasksWithErrorHandling() {
	tasks := store.GetAll()
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n📋 --- ALL TASKS ---")
	for _, t := range tasks {
		// Validate each task while displaying
		if err := validator.Validate(&t); err != nil {
			fmt.Printf("⚠️ Task %d has validation issues\n", t.ID)
		}
		fmt.Println(t.Display())
	}
}

func markCompleteWithErrorHandling() {
	listTasksWithErrorHandling()
	
	fmt.Print("Enter task ID to complete: ")
	id := readInt()
	
	err := store.MarkComplete(id)
	if err != nil {
		handleError(err, "Failed to mark task complete")
		return
	}
	
	fmt.Println("✅ Task marked as complete!")
	
	// Trigger auto-backup after modification
	go func() {
		if _, err := backupManager.CreateBackup(); err != nil {
			handleError(err, "Auto-backup failed")
		}
	}()
}

func setPriorityWithErrorHandling() {
	listTasksWithErrorHandling()
	
	fmt.Print("Enter task ID: ")
	id := readInt()
	
	fmt.Println("Priority levels:")
	fmt.Println("0: Low")
	fmt.Println("1: Medium")
	fmt.Println("2: High")
	fmt.Println("3: Critical")
	fmt.Print("Enter priority (0-3): ")
	p := readInt()
	
	err := store.SetPriority(id, task.Priority(p))
	if err != nil {
		handleError(err, "Failed to set priority")
		return
	}
	
	fmt.Println("✅ Priority updated!")
}

func addTagWithErrorHandling() {
	listTasksWithErrorHandling()
	
	fmt.Print("Enter task ID: ")
	id := readInt()
	
	fmt.Print("Enter tag: ")
	tag := readString()
	
	err := store.AddTag(id, tag)
	if err != nil {
		handleError(err, "Failed to add tag")
		return
	}
	
	fmt.Println("✅ Tag added!")
}

func validateSpecificTask() {
	listTasksWithErrorHandling()
	
	fmt.Print("Enter task ID to validate: ")
	id := readInt()
	
	task, err := store.Get(id)
	if err != nil {
		handleError(err, "Failed to get task")
		return
	}
	
	if err := validator.Validate(&task); err != nil {
		handleError(err, "Task validation failed")
	} else {
		fmt.Println("✅ Task is valid!")
	}
}

func createBackup() {
	fmt.Println("📦 Creating backup...")
	
	metadata, err := backupManager.CreateBackup()
	if err != nil {
		handleError(err, "Backup failed")
		return
	}
	
	fmt.Printf("✅ Backup created successfully!\n")
	fmt.Printf("   File: %s\n", metadata.File)
	fmt.Printf("   Tasks: %d\n", metadata.TaskCount)
	fmt.Printf("   Size: %d bytes\n", metadata.Size)
}

func listBackups() {
	backups, err := backupManager.ListBackups()
	if err != nil {
		handleError(err, "Failed to list backups")
		return
	}
	
	if len(backups) == 0 {
		fmt.Println("No backups found.")
		return
	}
	
	fmt.Println("\n📦 --- BACKUPS ---")
	for i, b := range backups {
		fmt.Printf("%d. %s (%s) - %d bytes\n", 
			i+1, b.File, b.Timestamp.Format("2006-01-02 15:04:05"), b.Size)
	}
}

func restoreFromBackup() {
	listBackups()
	
	fmt.Print("Enter backup number to restore (or 0 for latest): ")
	num := readInt()
	
	var err error
	if num == 0 {
		err = backupManager.RestoreLatest()
	} else {
		backups, _ := backupManager.ListBackups()
		if num < 1 || num > len(backups) {
			fmt.Println("Invalid backup number")
			return
		}
		err = backupManager.RestoreBackup(backups[num-1].File)
	}
	
	if err != nil {
		handleError(err, "Restore failed")
		return
	}
	
	fmt.Println("✅ Restore successful!")
	listTasksWithErrorHandling()
}

func testPanicRecovery() {
	fmt.Println("💥 Testing panic recovery...")
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("✅ Panic successfully recovered: %v\n", r)
		}
	}()
	
	// Deliberately cause a panic
	panic("This is a test panic - the app will recover!")
}

func readString() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt() int {
	input := readString()
	val, _ := strconv.Atoi(input)
	return val
}