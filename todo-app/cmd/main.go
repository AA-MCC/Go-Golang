package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var fileName = "todolist.json"
var todos []TodoTask
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type TodoTask struct {
	ID          string `json:"id"` //package level access as starts with uppercase
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	// Create a context with TraceID at the start
	ctx := context.Background()
	ctx = withTraceID(ctx)

	fmt.Println("Hello Go.. ")

	loadTodos(ctx)

	action := flag.String("action", "add", "Action : add, list, update, delete")
	desc := flag.String("desc", "", "Description for add/update")
	status := flag.String("status", "", "Todo task status: not started/started/completed")
	id := flag.String("id", "", "Todo task ID for update/delete")

	flag.Parse()

	logger.Info("Command parameters",
		"traceID", getTraceID(ctx),
		"action", *action,
		"description", *desc,
		"status", *status,
		"id", *id)

	// handle the commands
	switch *action {
	case "add":
		if *desc != "" {
			addTodo(ctx, *desc)
			saveTodos(ctx)
			listTodos(ctx)
		}
	case "list":
		listTodos(ctx)
	case "update":
		updateTodo(ctx, *id, *desc, *status)
		saveTodos(ctx)
	case "delete":
		deleteTodo(ctx, *id)
		saveTodos(ctx)
	default:
		logger.Error("Invalid action", "traceID", getTraceID(ctx), "action", *action)
	}
}

func saveTodos(ctx context.Context) {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		logger.Error("Error saving todos", "traceID", getTraceID(ctx), "error", err)
		return
	}
	logger.Info("Saved Todo", "traceID", getTraceID(ctx), "file", fileName)
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

func loadTodos(ctx context.Context) {
	data, err := os.ReadFile(fileName)
	if err == nil {
		_ = json.Unmarshal(data, &todos)
	}
	logger.Info("Loaded Todos", "traceID", getTraceID(ctx), "count", len(todos))
}

func listTodos(ctx context.Context) {
	logger.Info("Listing todos", "traceID", getTraceID(ctx), "count", len(todos))
	for _, t := range todos {
		logger.Info("Todo item",
			"traceID", getTraceID(ctx),
			"id", t.ID,
			"status", t.Status,
			"description", t.Description)
	}
}

func addTodo(ctx context.Context, description string) {
	logger.Info("Attempting to add todo", "traceID", getTraceID(ctx), "description", description)

	if description == "" {
		logger.Error("Description required", "traceID", getTraceID(ctx))
		return
	}
	id := uuid.New().String()

	newTodo := TodoTask{
		ID:          id,
		Description: description,
		Status:      "not started",
	}
	todos = append(todos, newTodo)
	logger.Info("Added todo",
		"traceID", getTraceID(ctx),
		"id", newTodo.ID,
		"description", newTodo.Description)
}

func updateTodo(ctx context.Context, id string, newDescription string, newStatus string) {
	logger.Info("Attempting to update todo", "traceID", getTraceID(ctx), "newDescription", newDescription)
	logger.Info("Attempting to update todo", "traceID", getTraceID(ctx), "newStatus", newStatus)
	if newDescription == "" && newStatus == "" {
		logger.Error("Description or status needed", "traceID", getTraceID(ctx))
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
			logger.Info("Updated todo",
				"traceID", getTraceID(ctx),
				"id", todos[i].ID,
				"description", todos[i].Description,
				"status", todos[i].Status)
			return
		}
	}
	logger.Error("Todo not found", "traceID", getTraceID(ctx), "id", id)
}

func deleteTodo(ctx context.Context, id string) {
	logger.Info("Attempting to delete todo", "traceID", getTraceID(ctx), "id", id)

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			logger.Info("Deleted todo",
				"traceID", getTraceID(ctx),
				"id", t.ID,
				"description", t.Description,
				"status", t.Status)
			return
		}
	}
	logger.Error("Todo not found for deletion", "traceID", getTraceID(ctx), "id", id)
}

func withTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, "traceID", uuid.New().String())
}

func getTraceID(ctx context.Context) string {
	if v := ctx.Value("traceID"); v != nil {
		return v.(string)
	}
	return ""
}
