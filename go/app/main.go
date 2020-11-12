package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"set1.ie.aitech.ac.jp/todo_sample/apifuncs"
)

func main() {
	http.HandleFunc("/tasks", apifuncs.TasksFunc)
	http.HandleFunc("/labels", apifuncs.LabelsFunc)

	http.ListenAndServe(":80", nil)
}
