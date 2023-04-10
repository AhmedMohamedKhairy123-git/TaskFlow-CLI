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
	
	// Add sample tasks with different priorities
	store.Add("Learn Go basics")
	store.Add("Complete Phase 4")
	store.Add("Practice interfaces")
	store.Add("Write documentation")
	
	store.SetPriority(1, task.High)
	store.SetPriority(2, task.Critical)
	store.SetPriority(3, task.Medium)
	
	store.AddTag(1, "learning")
	store.AddTag(2, "project")
	store.AddTag(3, "practice")
	
	store.MarkComplete(1)
	
	for {
		showMenu()
		if !processChoice() {
			break
		}
	}
	
	fmt.Println("\nFinal Statistics:")
	store.DisplayStats()
	fmt.Println("Goodbye!")
}

func showMenu() {
	fmt.Println("\n🎯 === TASK TRACKER (Methods & Interfaces) ===")
	fmt.Println("1. Add Task")
	fmt.Println("2. List All Tasks")
	fmt.Println("3. List Tasks (Simple)")
	fmt.Println("4. Mark Task Complete")
	fmt.Println("5. Set Task Priority")
	fmt.Println("6. Add Tag to Task")
	fmt.Println("7. Find by Priority")
	fmt.Println("8. Find by Tag")
	fmt.Println("9. Show Statistics")
	fmt.Println("10. Delete Task")
	fmt.Println("11. Exit")
	fmt.Print("Enter choice: ")
}

func processChoice() bool {
	choice := readInt()
	
	switch choice {
	case 1:
		addTask()
	case 2:
		listTasksDetailed()
	case 3:
		listTasksSimple()
	case 4:
		markComplete()
	case 5:
		setPriority()
	case 6:
		addTag()
	case 7:
		findByPriority()
	case 8:
		findByTag()
	case 9:
		store.DisplayStats()
	case 10:
		deleteTask()
	case 11:
		return false
	default:
		fmt.Println("Invalid choice!")
	}
	return true
}

func addTask() {
	fmt.Print("Enter task title: ")
	title := readString()
	
	task, err := store.Add(title)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Task added with ID: %d\n", task.ID)
}

func listTasksDetailed() {
	tasks := store.GetAll()
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n📋 --- ALL TASKS (Detailed) ---")
	for _, t := range tasks {
		fmt.Println(t.Display())
	}
}

func listTasksSimple() {
	tasks := store.GetAll()
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	
	fmt.Println("\n📋 --- ALL TASKS ---")
	for _, t := range tasks {
		fmt.Println(t.DisplaySimple())
	}
}

func markComplete() {
	listTasksSimple()
	
	fmt.Print("Enter task ID to complete: ")
	id := readInt()
	
	err := store.MarkComplete(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Task marked as complete!")
}

func setPriority() {
	listTasksSimple()
	
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
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Priority updated!")
}

func addTag() {
	listTasksSimple()
	
	fmt.Print("Enter task ID: ")
	id := readInt()
	
	fmt.Print("Enter tag: ")
	tag := readString()
	
	err := store.AddTag(id, tag)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Tag added!")
}

func findByPriority() {
	fmt.Println("Priority levels:")
	fmt.Println("0: Low")
	fmt.Println("1: Medium")
	fmt.Println("2: High")
	fmt.Println("3: Critical")
	fmt.Print("Enter priority to search (0-3): ")
	p := readInt()
	
	tasks := store.GetByPriority(task.Priority(p))
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found with this priority.")
		return
	}
	
	fmt.Printf("\n📋 --- %s PRIORITY TASKS ---\n", task.Priority(p))
	for _, t := range tasks {
		fmt.Println(t.Display())
	}
}

func findByTag() {
	fmt.Print("Enter tag to search: ")
	tag := readString()
	
	tasks := store.GetByTag(tag)
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found with this tag.")
		return
	}
	
	fmt.Printf("\n📋 --- TASKS WITH TAG '%s' ---\n", tag)
	for _, t := range tasks {
		fmt.Println(t.Display())
	}
}

func deleteTask() {
	listTasksSimple()
	
	fmt.Print("Enter task ID to delete: ")
	id := readInt()
	
	err := store.Delete(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Task deleted!")
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