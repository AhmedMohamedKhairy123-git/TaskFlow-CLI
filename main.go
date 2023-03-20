package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-tracker/task"
)

var store *task.TaskStore
var reader = bufio.NewReader(os.Stdin)

func main() {
	store = task.NewTaskStore()
	
	// Add some initial tasks
	store.AddTask("Learn Go basics")
	store.AddTask("Complete Phase 3")
	store.AddTask("Practice slices and maps")
	
	for {
		showMenu()
		if !processChoice() {
			break
		}
	}
	
	fmt.Println("Goodbye!")
}

func showMenu() {
	fmt.Println("\n=== TASK TRACKER ===")
	fmt.Println("1. Add Task")
	fmt.Println("2. List All Tasks")
	fmt.Println("3. Mark Task Complete")
	fmt.Println("4. Delete Task")
	fmt.Println("5. Find Tasks by Title")
	fmt.Println("6. Show Statistics")
	fmt.Println("7. Exit")
	fmt.Print("Enter choice: ")
}

func processChoice() bool {
	choice := readInt()
	
	switch choice {
	case 1:
		addTask()
	case 2:
		listTasks()
	case 3:
		markComplete()
	case 4:
		deleteTask()
	case 5:
		findTasks()
	case 6:
		store.DisplayStats()
	case 7:
		return false
	default:
		fmt.Println("Invalid choice!")
	}
	return true
}

func addTask() {
	fmt.Print("Enter task title: ")
	title := readString()
	
	if title == "" {
		fmt.Println("Title cannot be empty!")
		return
	}
	
	task := store.AddTask(title)
	fmt.Printf("Task added with ID: %d\n", task.ID)
}

func listTasks() {
	tasks := store.GetAllTasks()
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n--- ALL TASKS ---")
	for _, t := range tasks {
		status := " "
		if t.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
	}
}

func markComplete() {
	listTasks()
	
	if len(store.GetAllTasks()) == 0 {
		return
	}
	
	fmt.Print("Enter task ID to complete: ")
	id := readInt()
	
	if store.MarkComplete(id) {
		fmt.Println("Task marked as complete!")
	} else {
		fmt.Printf("Task with ID %d not found.\n", id)
	}
}

func deleteTask() {
	listTasks()
	
	if len(store.GetAllTasks()) == 0 {
		return
	}
	
	fmt.Print("Enter task ID to delete: ")
	id := readInt()
	
	if store.DeleteTask(id) {
		fmt.Println("Task deleted!")
	} else {
		fmt.Printf("Task with ID %d not found.\n", id)
	}
}

func findTasks() {
	fmt.Print("Enter title to search: ")
	title := readString()
	
	results := store.FindByTitle(title)
	
	if len(results) == 0 {
		fmt.Println("No matching tasks found.")
		return
	}
	
	fmt.Println("\n--- SEARCH RESULTS ---")
	for _, t := range results {
		status := " "
		if t.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
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