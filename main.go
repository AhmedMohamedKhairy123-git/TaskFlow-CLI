package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"task-tracker/config"
	"task-tracker/errors"
	"task-tracker/storage"
	"task-tracker/task"
	"task-tracker/web"
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
	startWebServer() // NEW: Start web server
	runApplication()
}

func handlePanic() {
	if r := recover(); r != nil {
		rec := errors.NewRecovery(r)
		fmt.Printf("\n🔥 CRITICAL: Application recovered from panic!\n")
		fmt.Printf("Recovery ID: %s\n", rec.ID)
		
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

// NEW: Start web server with graceful shutdown
func startWebServer() {
	server := web.NewServer(store)
	
	// Graceful shutdown on Ctrl+C
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		
		fmt.Println("\n🛑 Shutting down web server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use ctx in Stop method
if err := server.Stop(ctx); err != nil {
    fmt.Printf("Error stopping server: %v\n", err)
}
	}()
	
	// Start server in background
	go func() {
		if err := server.Start("8080"); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
	
	fmt.Println("🌐 Web API available at http://localhost:8080")
	fmt.Println("   Try: curl http://localhost:8080/tasks")
}

func initializeApp() {
	fmt.Println("🚀 Initializing Task Tracker")
	
	loadConfig()
	initStorage()
	store = task.NewTaskStore()
	validator = task.NewValidator()
	loadTasksFromFile()
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
	
	for _, t := range tasks {
		store.Add(t.Title)
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
	saveChan := make(chan bool, 1)
	
	go func() {
		for range autoSaveTicker.C {
			select {
			case saveChan <- true:
			default:
			}
		}
	}()
	
	go func() {
		for range saveChan {
			if store != nil {
				tasks := store.GetAll()
				if err := jsonStorage.Save(tasks); err != nil {
					fmt.Printf("⚠️ Auto-save failed: %v\n", err)
				}
			}
		}
	}()
	
	fmt.Println("🔄 Auto-saver started")
}

func runApplication() {
	for {
		showMenu()
		
		if !processChoiceSafe() {
			break
		}
	}
	
	saveTasks()
	fmt.Println("👋 Goodbye!")
}

func showMenu() {
	fmt.Println("\n📋 === TASK TRACKER MENU ===")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task Complete")
	fmt.Println("4. Set Priority")
	fmt.Println("5. Add Tag")
	fmt.Println("6. Save Tasks")
	fmt.Println("7. Load Tasks")
	fmt.Println("8. Create Backup")
	fmt.Println("9. Exit")
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
}

// Replace listTasks() function
func listTasks() {
	tasks := store.GetAll()
	
	if len(tasks) == 0 {
		fmt.Println(colorYellow + "📭 No tasks found." + colorReset)
		return
	}
	
	fmt.Printf("\n%s📋 TASKS%s\n", colorBold+colorCyan, colorReset)
	for _, t := range tasks {
		fmt.Println(t.ColorizedDisplay())
	}
}

// Replace showMenu() with colored menu
func showMenu() {
	fmt.Printf("\n%s🎯 TASK TRACKER MENU%s\n", colorBold+colorPurple, colorReset)
	fmt.Printf("%s1.%s Add Task\n", colorGreen, colorReset)
	fmt.Printf("%s2.%s List Tasks\n", colorBlue, colorReset)
	fmt.Printf("%s3.%s Mark Complete\n", colorGreen, colorReset)
	fmt.Printf("%s4.%s Set Priority\n", colorYellow, colorReset)
	fmt.Printf("%s5.%s Add Tag\n", colorCyan, colorReset)
	fmt.Printf("%s6.%s Save\n", colorWhite, colorReset)
	fmt.Printf("%s7.%s Load\n", colorWhite, colorReset)
	fmt.Printf("%s8.%s Backup\n", colorPurple, colorReset)
	fmt.Printf("%s9.%s Exit%s\n", colorRed, colorReset, colorReset)
	fmt.Printf("%sChoice:%s ", colorBold, colorReset)
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
	
	fmt.Printf("✅ Saved %d tasks to %s\n", len(tasks), appConfig.DataFile)
}

func loadTasks() {
	fmt.Println("📂 Loading tasks from file...")
	
	tasks, err := jsonStorage.Load()
	if err != nil {
		handleError(err, "Failed to load tasks")
		return
	}
	
	store = task.NewTaskStore()
	
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
	
	fmt.Printf("✅ Backup created: %s\n", backupFile)
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
// Add new helper functions
func promptString(message string) string {
	fmt.Printf("%s%s:%s ", colorCyan, message, colorReset)
	return readString()
}

func promptInt(message string) int {
	fmt.Printf("%s%s:%s ", colorYellow, message, colorReset)
	return readInt()
}

func promptConfirm(message string) bool {
	fmt.Printf("%s%s (y/n):%s ", colorPurple, message, colorReset)
	response := strings.ToLower(readString())
	return response == "y" || response == "yes"
}

func promptSelect(message string, options []string) int {
	fmt.Printf("%s%s:%s\n", colorCyan, message, colorReset)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt)
	}
	fmt.Printf("%sChoice:%s ", colorYellow, colorReset)
	choice := readInt()
	if choice < 1 || choice > len(options) {
		return 1
	}
	return choice
}

// Replace addTask()
func addTask() {
	title := promptString("Enter task title")
	
	// Sanitize input
	title = sanitizeInput(title)
	
	if title == "" {
		fmt.Printf("%s❌ Title cannot be empty!%s\n", colorRed, colorReset)
		return
	}
	
	// Validate length
	if len(title) < 3 {
		fmt.Printf("%s❌ Title too short (min 3 chars)%s\n", colorRed, colorReset)
		return
	}
	
	task, err := store.Add(title)
	if err != nil {
		handleError(err, "Failed to add task")
		return
	}
	
	fmt.Printf("%s✅ Task added with ID: %d%s\n", colorGreen, task.ID, colorReset)
	
	if promptConfirm("Add priority now?") {
		options := []string{"Low", "Medium", "High", "Critical"}
		choice := promptSelect("Select priority", options)
		store.SetPriority(task.ID, Priority(choice-1))
	}
}

// Replace deleteTask()
func deleteTask() {
	listTasks()
	
	if len(store.GetAll()) == 0 {
		return
	}
	
	id := promptInt("Enter task ID to delete")
	
	if promptConfirm(fmt.Sprintf("Delete task %d?", id)) {
		err := store.Delete(id)
		if err != nil {
			handleError(err, "Failed to delete")
			return
		}
		fmt.Printf("%s✅ Task deleted!%s\n", colorGreen, colorReset)
	}
}

// Replace exit handling
func runApplication() {
	for {
		showMenu()
		
		if !processChoiceSafe() {
			if promptConfirm("Really exit?") {
				break
			}
		}
	}
	
	if promptConfirm("Save before exit?") {
		saveTasks()
	}
	fmt.Printf("%s👋 Goodbye!%s\n", colorCyan, colorReset)
}
// Add command history and completion
var commandHistory []string
var historyIndex = -1

type Command struct {
	Name        string
	Description string
	Handler     func(args []string)
}

var commands = map[string]Command{
	"add":      {"add", "Add a new task", addTask},
	"list":     {"list", "List all tasks", func([]string) { listTasks() }},
	"complete": {"complete", "Mark task complete", func(args []string) { 
		if len(args) > 0 {
			id, _ := strconv.Atoi(args[0])
			store.MarkComplete(id)
		}
	}},
	"priority": {"priority", "Set priority", nil},
	"tag":      {"tag", "Add tag", nil},
	"save":     {"save", "Save tasks", func([]string) { saveTasks() }},
	"exit":     {"exit", "Exit program", func([]string) { 
		os.Exit(0) 
	}},
}

func showHelp() {
	fmt.Printf("\n%s📚 Available Commands:%s\n", colorBold+colorCyan, colorReset)
	for _, cmd := range commands {
		fmt.Printf("  %s%-10s%s %s\n", 
			colorGreen, cmd.Name, colorReset, cmd.Description)
	}
}

func processCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	
	cmdName := strings.ToLower(parts[0])
	args := parts[1:]
	
	// Add to history
	commandHistory = append(commandHistory, input)
	historyIndex = len(commandHistory)
	
	if cmdName == "help" {
		showHelp()
		return
	}
	
	if cmd, exists := commands[cmdName]; exists {
		if cmd.Handler != nil {
			cmd.Handler(args)
		} else {
			fmt.Printf("%s❌ Command '%s' not fully implemented%s\n", 
				colorRed, cmdName, colorReset)
		}
	} else {
		fmt.Printf("%s❌ Unknown command: %s%s\n", colorRed, cmdName, colorReset)
		fmt.Println("Type 'help' for available commands")
	}
}

// Replace readString with command handling
func readString() string {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	// Handle up/down arrows for history (simplified)
	if input == "\033[A" { // Up arrow
		if historyIndex > 0 {
			historyIndex--
			return commandHistory[historyIndex]
		}
		return ""
	}
	
	return input
}

// Add command mode toggling
var commandMode = false

func toggleCommandMode() {
	commandMode = !commandMode
	if commandMode {
		fmt.Printf("%s🔧 Command mode enabled (type 'help')%s\n", colorCyan, colorReset)
	} else {
		fmt.Printf("%s🔧 Menu mode enabled%s\n", colorCyan, colorReset)
	}
}
// Replace loadConfig()
func loadConfig() {
	// Check for profile flag
	profile := ProfileDev
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "dev":
			profile = ProfileDev
		case "test":
			profile = ProfileTest
		case "prod":
			profile = ProfileProd
		// Add to menu options
case "undo":
    if err := undoManager.Undo(); err != nil {
        fmt.Printf("%s❌ %s%s\n", colorRed, err, colorReset)
    } else {
        fmt.Printf("%s↩️ Undo successful%s\n", colorYellow, colorReset)
    }
case "redo":
    if err := undoManager.Redo(); err != nil {
        fmt.Printf("%s❌ %s%s\n", colorRed, err, colorReset)
    } else {
        fmt.Printf("%s↪️ Redo successful%s\n", colorYellow, colorReset)
    }
		}
	}
	
	appConfig = LoadProfile(profile)
	fmt.Printf("%s📋 Loaded %s profile%s\n", colorCyan, profile, colorReset)
	
	// Save config
	appConfig.Save("config.json")
}