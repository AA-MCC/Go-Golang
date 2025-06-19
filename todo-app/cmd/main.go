package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var pl = fmt.Println
var fileName = "todolist.json"
var todos []TodoTask

type TodoTask struct {
	ID          string `json:"id"` //package level access as starts with uppercase
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	pl("Hello Go.. \n")

	loadTodos()
	// listTodos()

	action := flag.String("action", "add", "Action : add, list, update, delete")

	desc := flag.String("desc", "", "Description for add/update")

	// listItems := flag.String("list", "", "List all todo tasks")
	// updateItem := flag.String("update", "", "Update todo task")
	// deleteItem := flag.String("delete", "", "Delete todo task")
	// status := flag.String("status", "not started", "Todo task status: not started/started/completed")

	flag.Parse()
	pl("description is :", *desc)
	pl("action is :", *action)

	// handle the commands - investigate options

	switch *action {
	case "add":
		if *desc != "" {
			addTodo(*desc)
			saveTodos()
			listTodos()
		}

	case "list":
		listTodos()
	default:
		pl("Invalid action.")

	}

}

// Saves the current todos slice to file
func saveTodos() {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		pl("Error saving todos:", err)
		return
	}
	pl("Saved Todo")
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
	// pl("load to do data: ", data)
	pl("Loaded Todos :", todos)
}

func listTodos() {
	pl("To-Do List:")
	for _, t := range todos {
		fmt.Printf("#%d [%s]: %s\n", t.ID, t.Status, t.Description)
	}
}

func addTodo(description string) {
	if description == "" {
		pl("Description is required to add a todo")
		return
	}
	id := uuid.New().String()
	// id := getNextID()
	newTodo := TodoTask{
		ID:          id,
		Description: description,
		Status:      "not started",
	}
	todos = append(todos, newTodo)
	pl("Added todo: ", newTodo.ID, newTodo.Description)
}
