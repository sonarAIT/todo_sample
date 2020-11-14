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
	// SQLを実行し，タスクの数だけid, name, description, submitTime, labelを取得．
	rows, err := db.Query("select id, name, description, submitTime, label from tasks")
	// エラー処理も行う．
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		// フォーマットに合わせてエラー，関数の名前，ファイル名，何行目かを表示している．
		// errFormatについては，dbctl/util.goを参照．
		log.Printf(errFormat, err, f.Name(), file, line)
		// 無のタスク一覧とエラーをTasksFuncに返す．(apifuncs/tasks.go)
		return nil, err
	}
	// あとでrows.Close()を行う．（この行は書き忘れないように！）
	defer rows.Close()

	// Taskの配列を宣言
	var Tasks []Task

	// SQLで取得した行数だけ実行
	for rows.Next() {
		// SQLの結果一行分を受け取る変数を宣言
		var id int
		var name string
		var description string
		var submitTime string
		var label int
		// 受け取る
		rows.Scan(&id, &name, &description, &submitTime, &label)
		// 受け取ったものをTasksに追加する．
		Tasks = append(Tasks, Task{ID: id, Name: name, Description: description, SubmitTime: submitTime, Label: label})
	}

	// Tasksとエラー（無）を返す．
	return Tasks, nil
}

// InsertTasks は挿入するべきTaskを全て挿入する関数
func InsertTasks(Tasks []Task) error {
	// 挿入するタスクの数だけ実行
	for _, task := range Tasks {
		// タスクを挿入するSQL
		_, err := db.Exec("insert into tasks values(?, ?, ?, ?, ?)", task.ID, task.Name, task.Description, task.SubmitTime, task.Label)
		// エラー処理
		if err != nil {
			pc, file, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.Printf(errFormat, err, f.Name(), file, line)
			// 今回はエラーのみ返す．
			return err
		}
	}

	// エラーが無いことを返す．
	return nil
}

// DeleteTasks は削除するべきTaskを全て削除する関数
func DeleteTasks(TaskIDs []int) error {
	// 削除するするタスクの数だけ実行
	for _, taskID := range TaskIDs {
		// 指定したIDのタスクを削除するSQL
		_, err := db.Exec("delete from tasks where id = ?", taskID)
		// エラー処理
		if err != nil {
			pc, file, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.Printf(errFormat, err, f.Name(), file, line)
			return err
		}
	}

	return nil
}
