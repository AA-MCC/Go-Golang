package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var fileName = "todolist.json"
var todos []TodoTask

type TodoTask struct {
	ID          string `json:"id"` //package level access as starts with uppercase
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	fmt.Println("Hello Go.. ")

	loadTodos()

	action := flag.String("action", "add", "Action : add, list, update, delete")
	desc := flag.String("desc", "", "Description for add/update")
	status := flag.String("status", "not started", "Todo task status: not started/started/completed")
	id := flag.String("id", "", "Todo task ID for update/delete")

	flag.Parse()

	fmt.Println("id is :", *id)
	fmt.Println("action is :", *action)
	fmt.Println("description is :", *desc)
	fmt.Println("Status is :", *status)

	// handle the commands

	switch *action {
	case "add":
		if *desc != "" {
			addTodo(*desc)
			saveTodos()
			listTodos()
		}
	case "list":
		listTodos()
	case "update":
		fmt.Println(*id, *desc, *status)
		updateTodo(*id, *desc, *status)
		saveTodos()
	case "delete":
		deleteTodo(*id)
		saveTodos()
	default:
		fmt.Println("Invalid action.")
	}

}

func saveTodos() {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		fmt.Println("Error saving todos:", err)
		return
	}
	fmt.Println("Saved Todo")
	_ = os.WriteFile(fileName, data, 0644)
}

// func saveTodos(data map[string][]byte) {
// 	file, err := os.Create(*file)
// 	if err != nil {
// 		log.Fatalf("failed to create file: %v", err)
// 	}
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	if err := encoder.Encode(data); err != nil {
// 		log.Fatalf("failed to encode data: %v", err)
// 	}
// }

func loadTodos() {
	data, err := os.ReadFile(fileName)
	if err == nil {
		_ = json.Unmarshal(data, &todos)

	}
	fmt.Println("Loaded Todos :", todos)
}

func listTodos() {
	fmt.Println("To-Do List:")
	for _, t := range todos {
		fmt.Println(t.ID, t.Status, t.Description)
	}
}

func addTodo(description string) {
	if description == "" {
		fmt.Println("Description is required to add a todo")
		return
	}
	id := uuid.New().String()

	newTodo := TodoTask{
		ID:          id,
		Description: description,
		Status:      "not started",
	}
	todos = append(todos, newTodo)
	fmt.Println("Added todo: ", newTodo.ID, newTodo.Description)
}

func updateTodo(id string, newDescription string, newStatus string) {
	listTodos()
	if newDescription == "" && newStatus == "" {
		fmt.Println("Description or status is needed to update")
		return
	}

	// Find and update the todo by ID
	for i, t := range todos {
		if t.ID == id {
			if newDescription != "" {
				todos[i].Description = newDescription
			}
			if newStatus != "" {
				todos[i].Status = newStatus
			}
			fmt.Println("Updated todo:", todos[i].ID, todos[i].Description, todos[i].Status)
			return
		}
	}
	fmt.Println("Todo with ID", id, "not found")
}

func deleteTodo(id string) {
	fmt.Println("todos in delete function:", todos)
	fmt.Println("id to be deleted is :", id)

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			fmt.Println("Deleted todo:", t.ID, t.Description, t.Status)
			return
		}
	}
	fmt.Println("Todo with ID", id, "not found")
}
