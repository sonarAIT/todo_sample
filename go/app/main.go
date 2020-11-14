package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"set1.ie.aitech.ac.jp/todo_sample/apifuncs"
)

// http.HandleFuncでリクエストに対応するURLと実行する関数を定義する．
func main() {
	// http://localhost:8081/tasks にアクセスされたら，apifuncs.TasksFuncを実行．
	http.HandleFunc("/tasks", apifuncs.TasksFunc)
	// http://localhost:8081/labels にアクセスされたら，apifuncs.LabelsFuncを実行．
	http.HandleFunc("/labels", apifuncs.LabelsFunc)
	// 80番ポートを使用．
	// dockerによってポートが変更されるので，実際に使用するのは8081番ポートである．
	http.ListenAndServe(":80", nil)
}
