package todostore

import (
	"context"
	"testing"
)

func TestAddTodo(t *testing.T) {
	store := NewStore("")
	desc := "A todo task"
	todo := store.AddTodo(context.Background(), desc)
	if todo == nil || todo.Description != desc {
		t.Fatalf("expected todo with description %q, got %+v", desc, todo)
	}
	if todo.Status != "not started" {
		t.Errorf("expected status 'not started', got %q", todo.Status)
	}
}

func TestUpdateTodo(t *testing.T) {
	store := NewStore("")
	todo := store.AddTodo(context.Background(), "Initial")
	ok := store.UpdateTodo(context.Background(), todo.ID, "Updated", "completed")
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if store.todos[0].Description != "Updated" || store.todos[0].Status != "completed" {
		t.Errorf("update did not apply correctly: %+v", store.todos[0])
	}
}

func TestDeleteTodo(t *testing.T) {
	store := NewStore("")
	todo := store.AddTodo(context.Background(), "To delete")
	ok := store.DeleteTodo(context.Background(), todo.ID)
	if !ok {
		t.Fatal("expected delete to succeed")
	}
	if len(store.todos) != 0 {
		t.Errorf("expected todos to be empty after delete, got %+v", store.todos)
	}
}

func TestListTodos(t *testing.T) {
	store := NewStore("")
	store.AddTodo(context.Background(), "Task 1")
	store.AddTodo(context.Background(), "Task 2")
	todos := store.ListTodos(context.Background())
	if len(todos) != 2 {
		t.Errorf("expected 2 todos, got %d", len(todos))
	}
}
