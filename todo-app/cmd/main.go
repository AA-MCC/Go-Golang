package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"todo-app/todostore"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

var fileName = "todolist.json"
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ping %s\n", r.URL.Query().Get("name"))
}

// API Handlers
func createHandler(store *todostore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		todo := store.AddTodo(r.Context(), req.Description)
		if todo == nil {
			http.Error(w, "Description required", http.StatusBadRequest)
			return
		}
		// _ = store.SaveTodos(r.Context())
		store.SaveTodos(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
}

func getHandler(store *todostore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		for _, t := range store.ListTodos(r.Context()) {
			if t.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(t)
				return
			}
		}
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func updateHandler(store *todostore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		ok := store.UpdateTodo(r.Context(), req.ID, req.Description, req.Status)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		_ = store.SaveTodos(r.Context())
		w.WriteHeader(http.StatusOK)
	}
}

func deleteHandler(store *todostore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			req.ID = id
		}
		ok := store.DeleteTodo(r.Context(), req.ID)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		_ = store.SaveTodos(r.Context())
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	ctx := context.Background()
	ctx = withTraceID(ctx)

	fmt.Println("Hello Go.. ")

	store := todostore.NewStore(fileName)
	fmt.Println("Loaded store", store)
	_ = store.LoadTodos(ctx)

	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*** Within /about functionality. ")
		http.ServeFile(w, r, "./static/about.html")
	})

	// mux.Handle("/about/", http.StripPrefix("/about/", http.FileServer(http.Dir(filepath.Join("..", "../static")))))

	action := flag.String("action", "add", "Action : add, list, update, delete, serve")
	desc := flag.String("desc", "", "Description for add/update")
	status := flag.String("status", "", "Todo task status: not started/started/completed")
	id := flag.String("id", "", "Todo task ID for update/delete")
	flag.Parse()

	// Serve dynamic page for the "/list" endpoint
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Within /list functionality. ")
		// Parse the template
		tmpl, err := template.ParseFiles("./templates/list.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			logger.Error("Error parsing template", "traceID", getTraceID(ctx), "error", err)
			return
		}

		logger.Info("Running Execute on template",
			"traceID", getTraceID(ctx),
			"action", *action,
			"description", *desc,
			"status", *status,
			"id", *id)

		fmt.Println("store todos :", store.ListTodos(ctx))

		// Render template with the todos data
		tmpl.Execute(w, store.ListTodos(ctx))

	})

	logger.Info("Command parameters",
		"traceID", getTraceID(ctx),
		"action", *action,
		"description", *desc,
		"status", *status,
		"id", *id)

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	// Start http server
	go func() {
		fmt.Println("Serve 8080")
		logger.Info("Starting server 8080", "traceID", getTraceID(ctx), "action", *action)
		if err := http.ListenAndServe(":8080", mux); err != nil {
			logger.Error("Server failed to start", "error", err)
		}
	}()

	fmt.Println("Served 8080")

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

	mux.HandleFunc("/create", createHandler(store))
	mux.HandleFunc("/get", getHandler(store))
	mux.HandleFunc("/update", updateHandler(store))
	mux.HandleFunc("/delete", deleteHandler(store))

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
