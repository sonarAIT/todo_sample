package main

import (
	"net/http"

	//"set1.ie.aitech.ac.jp/todo_sample/apifuncs"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// http.HandleFunc("/hoge", apifuncs.HogeFunc)

	http.ListenAndServe(":80", nil)
}
