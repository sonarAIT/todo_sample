package dbctl

import (
	"log"
	"runtime"
)

// Task はタスク1つ分の構造体
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SubmitTime  string `json:"submitTime"`
	Label       int    `json:"label"`
}

// GetTasks タスク一覧を取得する関数
func GetTasks() ([]Task, error) {
	rows, err := db.Query("select id, name, description, submitTime, label from tasks")
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.Printf(errFormat, err, f.Name(), file, line)
		return nil, err
	}
	defer rows.Close()

	var Tasks []Task

	for rows.Next() {
		var id int
		var name string
		var description string
		var submitTime string
		var label int
		rows.Scan(&id, &name, &description, &submitTime, &label)
		Tasks = append(Tasks, Task{ID: id, Name: name, Description: description, SubmitTime: submitTime, Label: label})
	}

	return Tasks, nil
}

// InsertTasks は挿入するべきTaskを全て挿入する関数
func InsertTasks(Tasks []Task) error {
	for _, task := range Tasks {
		_, err := db.Exec("insert into tasks values(?, ?, ?, ?, ?)", task.ID, task.Name, task.Description, task.SubmitTime, task.Label)
		if err != nil {
			pc, file, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.Printf(errFormat, err, f.Name(), file, line)
			return err
		}
	}

	return nil
}

// DeleteTasks は削除するべきTaskを全て削除する関数
func DeleteTasks(TaskIDs []int) error {
	for _, taskID := range TaskIDs {
		_, err := db.Exec("delete from tasks where id = ?", taskID)
		if err != nil {
			pc, file, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.Printf(errFormat, err, f.Name(), file, line)
			return err
		}
	}

	return nil
}
