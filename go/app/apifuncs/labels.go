package apifuncs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"set1.ie.aitech.ac.jp/todo_sample/dbctl"
)

//LabelsFunc は/labelsにアクセスされた際に実行される関数
func LabelsFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
	w.Header().Set("Access-Control-Allow-Headers", "*")

	r.Header.Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		Labels, err := dbctl.GetLabels()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("database error(GetLabels)", err)
			return
		}

		jsonBytes, err := json.Marshal(Labels)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("JSON Marshal error(Labels)", err)
			return
		}

		jsonString := string(jsonBytes)

		w.WriteHeader(http.StatusOK)

		if Labels == nil {
			jsonString = "[]"
		}

		fmt.Fprintln(w, jsonString)
	}

}
