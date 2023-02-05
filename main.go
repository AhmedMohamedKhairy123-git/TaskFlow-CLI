package main

import (
	"fmt"
)

// Task represents a single task in our tracker
type Task struct {
	ID        int
	Title     string
	Completed bool
}

const (
	appName    = "Task Tracker"
	appVersion = "1.0.0"
)

func main() {
	// Basic variables demonstration
	var userName string
	var userChoice int
	
	// Constants and basic output
	fmt.Printf("Welcome to %s v%s\n", appName, appVersion)
	fmt.Println("================================")
	
	// Get user name
	fmt.Print("Enter your name: ")
	fmt.Scanln(&userName)
	
	// Create a simple task to demonstrate struct
	task1 := Task{
		ID:        1,
		Title:     "Learn Go basics",
		Completed: false,
	}
	
	task2 := Task{2, "Complete Phase 1", false} // shorthand notation
	
	// Display welcome and tasks
	fmt.Printf("\nHello, %s! Here are your initial tasks:\n\n", userName)
	
	// Manual task display (we'll improve this later)
	fmt.Printf("Task #%d: %s - Completed: %v\n", task1.ID, task1.Title, task1.Completed)
	fmt.Printf("Task #%d: %s - Completed: %v\n", task2.ID, task2.Title, task2.Completed)
	
	// Simple menu placeholder
	fmt.Println("\n--- Menu Options (coming in Phase 2) ---")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Exit")
	fmt.Print("\nEnter your choice (demo only, press any number): ")
	fmt.Scanln(&userChoice)
	
	fmt.Printf("\nYou selected option %d. More features coming in Phase 2!\n", userChoice)
	fmt.Println("Phase 1 complete! Ready for commit.")
}