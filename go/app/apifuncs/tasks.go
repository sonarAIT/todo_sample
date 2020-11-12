package apifuncs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"set1.ie.aitech.ac.jp/todo_sample/dbctl"
)

//TasksFunc は/tasksにアクセスされた際に実行される関数
func TasksFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
	w.Header().Set("Access-Control-Allow-Headers", "*")

	r.Header.Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		Tasks, err := dbctl.GetTasks()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("database error(GetTasks)", err)
			return
		}

		jsonBytes, err := json.Marshal(Tasks)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("JSON Marshal error(Tasks)", err)
			return
		}

		jsonString := string(jsonBytes)

		w.WriteHeader(http.StatusOK)

		if Tasks == nil {
			jsonString = "[]"
		}

		fmt.Fprintln(w, jsonString)
	} else if r.Method == http.MethodPost {
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("Can't catch Tasks(io error)", err)
			return
		}

		var recTasks []dbctl.Task
		if err := json.Unmarshal(jsonBytes, &recTasks); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("Can't catch Tasks(JSON Unmarshal error)", err)
			return
		}
		log.Print(recTasks)

		var dbTasks []dbctl.Task
		dbTasks, err = dbctl.GetTasks()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("database error(GetTasks)", err)
			return
		}

		var insertTasks []dbctl.Task

		for _, recTask := range recTasks {
			var findFlag = false
			for _, dbTask := range dbTasks {
				if recTask.ID == dbTask.ID {
					findFlag = true
					break
				}
			}
			if !findFlag {
				insertTasks = append(insertTasks, recTask)
			}
		}

		if insertTasks != nil {
			err := dbctl.InsertTasks(insertTasks)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintln(w, `{"status":"Unavailable"}`)
				fmt.Println("database error(InsertTasks)", err)
				return
			}
		}

		var deleteTaskIDs []int

		for _, dbTask := range dbTasks {
			var findFlag = false
			for _, recTask := range recTasks {
				if recTask.ID == dbTask.ID {
					findFlag = true
					break
				}
			}
			if !findFlag {
				deleteTaskIDs = append(deleteTaskIDs, dbTask.ID)
			}
		}

		if deleteTaskIDs != nil {
			err := dbctl.DeleteTasks(deleteTaskIDs)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintln(w, `{"status":"Unavailable"}`)
				fmt.Println("database error(DeleteTasks)", err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"Available"}`)
	}
}
