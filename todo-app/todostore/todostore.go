package todostore

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type TodoTask struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Store struct {
	fileName string
	todos    []TodoTask
}

func NewStore(fileName string) *Store {
	return &Store{fileName: fileName}
}

func (s *Store) LoadTodos(ctx context.Context) error {
	data, err := os.ReadFile(s.fileName)
	if err == nil {
		_ = json.Unmarshal(data, &s.todos)
	}
	return nil
}

func (s *Store) SaveTodos(ctx context.Context) error {
	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileName, data, 0644)
}

func (s *Store) ListTodos(ctx context.Context) []TodoTask {
	return s.todos
}

func (s *Store) AddTodo(ctx context.Context, description string) *TodoTask {
	if description == "" {
		return nil
	}
	id := uuid.New().String()
	newTodo := TodoTask{
		ID:          id,
		Description: description,
		Status:      "not started",
	}
	s.todos = append(s.todos, newTodo)
	return &newTodo
}

func (s *Store) UpdateTodo(ctx context.Context, id string, newDescription string, newStatus string) bool {
	for i, t := range s.todos {
		if t.ID == id {
			if newDescription != "" {
				s.todos[i].Description = newDescription
			}
			if newStatus != "" {
				s.todos[i].Status = newStatus
			}
			return true
		}
	}
	return false
}

func (s *Store) DeleteTodo(ctx context.Context, id string) bool {
	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return true
		}
	}
	return false
}
