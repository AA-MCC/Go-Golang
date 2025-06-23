package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"todo-app/todostore"

	"crypto/rand"
	"encoding/hex"
)

var fileName = "todolist.json"
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	ctx := context.Background()
	ctx = withTraceID(ctx)

	fmt.Println("Hello Go.. ")

	store := todostore.NewStore(fileName)
	_ = store.LoadTodos(ctx)

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

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		// Now have received signal
		logger.Info("Received interrupt signal, exiting...", "traceID", getTraceID(ctx))
		// Signal that we are done
		done <- true
	}()

	switch *action {
	case "add":
		if *desc != "" {
			store.AddTodo(ctx, *desc)
			_ = store.SaveTodos(ctx)
			printTodos(ctx, store)
		}
	case "list":
		printTodos(ctx, store)
	case "update":
		store.UpdateTodo(ctx, *id, *desc, *status)
		_ = store.SaveTodos(ctx)
	case "delete":
		store.DeleteTodo(ctx, *id)
		_ = store.SaveTodos(ctx)
	default:
		logger.Error("Invalid action", "traceID", getTraceID(ctx), "action", *action)
	}

	fmt.Println("Press Ctrl+C to exit...")
	<-done
}

func printTodos(ctx context.Context, store *todostore.Store) {
	todos := store.ListTodos(ctx)
	logger.Info("Listing todos", "traceID", getTraceID(ctx), "count", len(todos))
	for _, t := range todos {
		logger.Info("Todo item",
			"traceID", getTraceID(ctx),
			"id", t.ID,
			"status", t.Status,
			"description", t.Description)
	}
}

func withTraceID(ctx context.Context) context.Context {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return context.WithValue(ctx, "traceID", hex.EncodeToString(b))
}

func getTraceID(ctx context.Context) string {
	if v := ctx.Value("traceID"); v != nil {
		return v.(string)
	}
	return ""
}
