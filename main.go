package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

const (
	appName    = "Task Tracker"
	appVersion = "1.0.0"
)

var tasks []Task
var nextID int = 1
var reader = bufio.NewReader(os.Stdin)

func main() {
	showWelcome()
	
	for {
		showMenu()
		choice := getUserChoice()
		
		if !processChoice(choice) {
			break
		}
	}
	
	fmt.Println("Goodbye!")
}

func showWelcome() {
	fmt.Printf("Welcome to %s v%s\n", appName, appVersion)
	fmt.Println("================================")
}

func showMenu() {
	fmt.Println("\n--- MENU ---")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task Complete")
	fmt.Println("4. Exit")
	fmt.Print("Enter choice: ")
}

func getUserChoice() int {
	var choice int
	fmt.Scanln(&choice)
	return choice
}

func processChoice(choice int) bool {
	switch choice {
	case 1:
		addTask()
		return true
	case 2:
		listTasks()
		return true
	case 3:
		markTaskComplete()
		return true
	case 4:
		return false
	default:
		fmt.Println("Invalid choice. Please try again.")
		return true
	}
}

func addTask() {
	fmt.Print("Enter task title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	
	if title == "" {
		fmt.Println("Task title cannot be empty!")
		return
	}
	
	task := Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}
	
	tasks = append(tasks, task)
	nextID++
	
	fmt.Printf("Task added with ID: %d\n", task.ID)
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n--- TASKS ---")
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Title)
	}
}

func markTaskComplete() {
	if len(tasks) == 0 {
		fmt.Println("No tasks to complete.")
		return
	}
	
	listTasks()
	
	fmt.Print("Enter task ID to mark complete: ")
	var id int
	fmt.Scanln(&id)
	
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			fmt.Printf("Task '%s' marked as complete!\n", task.Title)
			return
		}
	}
	
	fmt.Printf("Task with ID %d not found.\n", id)
}