package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"task-tracker/config"
	"task-tracker/errors"
	"task-tracker/storage"
	"task-tracker/task"
)

var (
	store         *task.TaskStore
	validator     *task.Validator
	backupManager *task.BackupManager
	jsonStorage   *storage.JSONStorage
	appConfig     *config.Config
	reader        = bufio.NewReader(os.Stdin)
	autoSaveTicker *time.Ticker
)

func main() {
	defer handlePanic()
	initializeApp()
	startAutoSave()
	runApplication()
}

func handlePanic() {
	if r := recover(); r != nil {
		rec := errors.NewRecovery(r)
		fmt.Printf("\n🔥 CRITICAL: Application recovered from panic!\n")
		fmt.Printf("Recovery ID: %s\n", rec.ID)
		
		// Emergency save
		if store != nil && jsonStorage != nil {
			fmt.Println("Attempting emergency save...")
			tasks := store.GetAll()
			if err := jsonStorage.Save(tasks); err != nil {
				fmt.Printf("Emergency save failed: %v\n", err)
			} else {
				fmt.Println("✅ Emergency save successful!")
			}
		}
	}
}

func initializeApp() {
	fmt.Println("🚀 Initializing Task Tracker with File I/O")
	
	// Load configuration
	loadConfig()
	
	// Initialize storage
	initStorage()
	
	// Initialize components
	store = task.NewTaskStore()
	validator = task.NewValidator()
	
	// Load tasks from file
	loadTasksFromFile()
	
	// Initialize backup manager with storage
	backupManager = task.NewBackupManager(store, appConfig.BackupDir)
	
	fmt.Println("✅ Initialization complete!")
}

func loadConfig() {
	var err error
	appConfig, err = config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("⚠️ Warning: Using default config (%v)\n", err)
		appConfig = config.DefaultConfig()
	}
	
	// Save config if it doesn't exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		appConfig.Save("config.json")
	}
}

func initStorage() {
	jsonStorage = storage.NewJSONStorage(
		appConfig.DataFile,
		appConfig.BackupDir,
	)
}

func loadTasksFromFile() {
	tasks, err := jsonStorage.Load()
	if err != nil {
		handleError(err, "Failed to load tasks from file")
		return
	}
	
	// Populate store
	for _, t := range tasks {
		store.Add(t.Title)
		// Restore additional properties
		if task, exists := store.Tasks[t.ID]; exists {
			task.Completed = t.Completed
			task.Priority = t.Priority
			task.Tags = t.Tags
			task.CreatedAt = t.CreatedAt
		}
	}
	
	fmt.Printf("📂 Loaded %d tasks from %s\n", len(tasks), appConfig.DataFile)
}

func startAutoSave() {
	if !appConfig.AutoSave {
		return
	}
	
	autoSaveTicker = time.NewTicker(time.Duration(appConfig.SaveInterval) * time.Second)
	
	go func() {
		for range autoSaveTicker.C {
			if store != nil {
				tasks := store.GetAll()
				if err := jsonStorage.Save(tasks); err != nil {
					fmt.Printf("⚠️ Auto-save failed: %v\n", err)
				} else {
					fmt.Printf("💾 Auto-saved %d tasks\n", len(tasks))
				}
			}
		}
	}()
}

func runApplication() {
	for {
		showMenu()
		
		if !processChoiceSafe() {
			break
		}
	}
	
	// Save before exit
	saveTasks()
	fmt.Println("👋 Goodbye!")
}

func showMenu() {
	fmt.Println("\n💾 === TASK TRACKER (File I/O) ===")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task Complete")
	fmt.Println("4. Set Priority")
	fmt.Println("5. Add Tag")
	fmt.Println("6. Save Tasks")
	fmt.Println("7. Load Tasks")
	fmt.Println("8. Create Backup")
	fmt.Println("9. List Backups")
	fmt.Println("10. Restore from Backup")
	fmt.Println("11. View File Info")
	fmt.Println("12. Exit")
	fmt.Print("Enter choice: ")
}

func processChoiceSafe() bool {
	defer recoverFromPanic()
	
	choice := readInt()
	
	switch choice {
	case 1:
		addTask()
	case 2:
		listTasks()
	case 3:
		markComplete()
	case 4:
		setPriority()
	case 5:
		addTag()
	case 6:
		saveTasks()
	case 7:
		loadTasks()
	case 8:
		createBackup()
	case 9:
		listBackups()
	case 10:
		restoreFromBackup()
	case 11:
		viewFileInfo()
	case 12:
		return false
	default:
		fmt.Println("Invalid choice!")
	}
	return true
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Printf("⚠️ Recovered from panic: %v\n", r)
	}
}

func addTask() {
	fmt.Print("Enter task title: ")
	title := readString()
	
	if title == "" {
		fmt.Println("Title cannot be empty!")
		return
	}
	
	task, err := store.Add(title)
	if err != nil {
		handleError(err, "Failed to add task")
		return
	}
	
	fmt.Printf("✅ Task added with ID: %d\n", task.ID)
	
	if appConfig.AutoSave {
		saveTasks()
	}
}

func listTasks() {
	tasks := store.GetAll()
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n📋 --- ALL TASKS ---")
	for _, t := range tasks {
		fmt.Println(t.Display())
	}
}

func markComplete() {
	listTasks()
	
	fmt.Print("Enter task ID to complete: ")
	id := readInt()
	
	err := store.MarkComplete(id)
	if err != nil {
		handleError(err, "Failed to mark task complete")
		return
	}
	
	fmt.Println("✅ Task marked as complete!")
	
	if appConfig.AutoSave {
		saveTasks()
	}
}

func setPriority() {
	listTasks()
	
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

func addTag() {
	listTasks()
	
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

func saveTasks() {
	fmt.Println("💾 Saving tasks...")
	
	tasks := store.GetAll()
	if err := jsonStorage.Save(tasks); err != nil {
		handleError(err, "Failed to save tasks")
		return
	}
	
	// Get file info
	info, err := os.Stat(appConfig.DataFile)
	if err != nil {
		fmt.Printf("✅ Saved %d tasks to %s\n", len(tasks), appConfig.DataFile)
	} else {
		fmt.Printf("✅ Saved %d tasks to %s (%d bytes)\n", 
			len(tasks), appConfig.DataFile, info.Size())
	}
}

func loadTasks() {
	fmt.Println("📂 Loading tasks from file...")
	
	tasks, err := jsonStorage.Load()
	if err != nil {
		handleError(err, "Failed to load tasks")
		return
	}
	
	// Clear current store
	store = task.NewTaskStore()
	
	// Populate store
	for _, t := range tasks {
		store.Add(t.Title)
		if task, exists := store.Tasks[t.ID]; exists {
			task.Completed = t.Completed
			task.Priority = t.Priority
			task.Tags = t.Tags
			task.CreatedAt = t.CreatedAt
		}
	}
	
	fmt.Printf("✅ Loaded %d tasks from %s\n", len(tasks), appConfig.DataFile)
}

func createBackup() {
	fmt.Println("📦 Creating backup...")
	
	backupFile, err := jsonStorage.Backup()
	if err != nil {
		handleError(err, "Failed to create backup")
		return
	}
	
	info, err := os.Stat(backupFile)
	if err != nil {
		fmt.Printf("✅ Backup created: %s\n", backupFile)
	} else {
		fmt.Printf("✅ Backup created: %s (%d bytes)\n", backupFile, info.Size())
	}
}

func listBackups() {
	backups, err := jsonStorage.ListBackups()
	if err != nil {
		handleError(err, "Failed to list backups")
		return
	}
	
	if len(backups) == 0 {
		fmt.Println("No backups found.")
		return
	}
	
	fmt.Println("\n📦 --- BACKUPS ---")
	for i, backup := range backups {
		info, err := os.Stat(backup)
		size := "unknown"
		if err == nil {
			size = fmt.Sprintf("%d bytes", info.Size())
		}
		fmt.Printf("%d. %s (%s)\n", i+1, filepath.Base(backup), size)
	}
}

func restoreFromBackup() {
	backups, err := jsonStorage.ListBackups()
	if err != nil || len(backups) == 0 {
		fmt.Println("No backups available.")
		return
	}
	
	listBackups()
	
	fmt.Print("Enter backup number to restore: ")
	num := readInt()
	
	if num < 1 || num > len(backups) {
		fmt.Println("Invalid backup number")
		return
	}
	
	fmt.Printf("Restoring from %s...\n", filepath.Base(backups[num-1]))
	
	if err := jsonStorage.Restore(backups[num-1]); err != nil {
		handleError(err, "Restore failed")
		return
	}
	
	// Reload tasks
	loadTasks()
	fmt.Println("✅ Restore complete!")
}

func viewFileInfo() {
	info, err := os.Stat(appConfig.DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist yet.\n", appConfig.DataFile)
		} else {
			handleError(err, "Failed to get file info")
		}
		return
	}
	
	fmt.Printf("\n📁 File: %s\n", appConfig.DataFile)
	fmt.Printf("Size: %d bytes\n", info.Size())
	fmt.Printf("Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("Permissions: %s\n", info.Mode())
	
	// Count tasks in file
	tasks, err := jsonStorage.Load()
	if err != nil {
		fmt.Printf("Tasks in file: error reading (%v)\n", err)
	} else {
		fmt.Printf("Tasks in file: %d\n", len(tasks))
	}
}

func handleError(err error, context string) {
	if err == nil {
		return
	}
	
	fmt.Printf("\n❌ ERROR: %s\n", context)
	
	if appErr, ok := err.(*errors.AppError); ok {
		fmt.Printf("   Code: %s\n", appErr.Code)
		fmt.Printf("   Message: %s\n", appErr.Message)
	} else {
		fmt.Printf("   %v\n", err)
	}
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