This project is to create a to-do app which allows users to create a list of to-do tasks.

It is still in development!

Usage

To add a task:-
go run todo-app\cmd\main.go -action add -desc "Buy bread" 

To list all tasks:-
go run todo-app\cmd\main.go -action list

To update a task:-
Use the "log/slog" structured logging package to log errors and when data is saved to disk
- Use the "context" package to add a TraceID to enable traceability of calls through the solution by adding it to all logs

To delete a task:-
 go run todo-app\cmd\main.go -action delete -id  "72cc896b-ef77-4f06-987a-094fd3e47e58"




