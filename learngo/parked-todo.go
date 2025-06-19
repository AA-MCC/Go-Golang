package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var pl = fmt.Println
var file = "todolist.json"

type TodoTask struct {
	ID          int    `json:id`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	pl("Hello Go.. \n")

	//load todos

	addItem := flag.String("add", "", "Add a new todo task")
	listItems := flag.String("list", "", "List all todo tasks")
	updateItem := flag.String("update", "", "Update todo task")
	deleteItem := flag.String("delete", "", "Delete todo task")
	status := flag.String("status", "not started", "Todo task status: not started/started/completed")

	flag.Parse()

	// handle the commands - investigate options

	switch {
	case *addItem:
		addTodo(*additem)

	}

	pl("addItem :", *addItem)

	todos := []string{}

	if *addItem != "" {
		todos = append(todos, *addItem)
		pl("Added to-do task:", *addItem)
	}

	data, err := json.Marshal(todo)
	saveTodos()

	// reader adn a writer

	// try json encoder or encoding and decoding text encoding
	// https://pkg.go.dev/golang.org/x/text/encoding

	// os.File

	// os.WriteFile(
	// check read write permissions

}

func saveTodos(data map[string][]byte) {
	file, err := os.Create(*file)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		log.Fatalf("failed to encode data: %v", err)
	}
}

func loadTodos() {
	file, err := os.ReadFile(todos_list_File)

}
