package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var pl = fmt.Println
var fileName = "todolist.json"
var todos []TodoTask

type TodoTask struct {
	ID          int    `json:id` //package level access as starts with uppercase
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	pl("Hello Go.. \n")

	loadTodos()
	listTodos()

	addItem := flag.String("add", "", "Add a new todo task")

	action := flag.String("action", "add", "Action : add, list, update, delete")

	// listItems := flag.String("list", "", "List all todo tasks")
	// updateItem := flag.String("update", "", "Update todo task")
	// deleteItem := flag.String("delete", "", "Delete todo task")
	// status := flag.String("status", "not started", "Todo task status: not started/started/completed")

	flag.Parse()

	//load todos

	// handle the commands - investigate options

	strtodos := []string{}

	switch *action {
	case "add":
		if *addItem != "" {
		}
		strtodos = append(strtodos, *addItem)
		pl("Added to-do task:", *addItem)
		// addTodo()
	case "list":
		listTodos()
	default:
		pl("Invalid action, only use add, list, update, or delete.")

	}

	pl("addItem :", *addItem)

	if *addItem != "" {
		strtodos = append(strtodos, *addItem)
		pl("Added to-do task:", *addItem)
	}

	// data, err := json.Marshal(todo)
	// saveTodos()

	// reader adn a writer

	// try json encoder or encoding and decoding text encoding
	// https://pkg.go.dev/golang.org/x/text/encoding

	// os.File

	// os.WriteFile(
	// check read write permissions

}

// func saveTodos() {
// 	data, err := json.MarshalIndent(todos, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error saving todos:", err)
// 		return
// 	}
// 	_ = os.WriteFile(fileName, data, 0644)
// }

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
}

func listTodos() {
	fmt.Println("To-Do List:")
	for _, t := range todos {
		fmt.Printf("#%d [%s]: %s\n", t.ID, t.Status, t.Description)
	}
}

// func addTodo() {
// 	if *desc == "" {
// 		fmt.Println("Description is required to add a todo")
// 		return
// 	}
// 	id := getNextID()
// 	todos = append(todos, Todo{ID: id, Description: *desc, Status: NotStarted})
// 	fmt.Printf("Added todo: #%d - %s\n", id, *desc)
// }
