package apifuncs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"set1.ie.aitech.ac.jp/todo_sample/dbctl"
)

//TasksFunc は/tasksにアクセスされた際に実行される関数
func TasksFunc(w http.ResponseWriter, r *http.Request) {
	// CORSの設定．
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// ヘッダの設定．
	r.Header.Set("Content-Type", "application/json")

	// GETの場合
	// /tasks GET では，クライアントに対してタスク一覧を送信する．
	if r.Method == http.MethodGet {
		// Task一覧を取得．（dbctl/tasks.goを参照．）
		Tasks, err := dbctl.GetTasks()
		// エラー処理
		if err != nil {
			// ヘッダーに失敗したことを書き込む
			w.WriteHeader(http.StatusServiceUnavailable)
			// ついでに失敗したことをフロントがJSONとして認識できるように書き込む
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			// ログにも書いちゃう
			fmt.Println("database error(GetTasks)", err)
			// 終了
			return
		}

		// 取得したタスク一覧をByte列に変換
		jsonBytes, err := json.Marshal(Tasks)
		// エラー処理
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("JSON Marshal error(Tasks)", err)
			return
		}

		// 文字列に変換（送信可能になった）
		jsonString := string(jsonBytes)

		// 成功したことをヘッダーに書き込む
		w.WriteHeader(http.StatusOK)

		// Tasksがもし空なら，空配列を返す．
		if Tasks == nil {
			jsonString = "[]"
		}

		// 結果を書き込んで終了．
		fmt.Fprintln(w, jsonString)

		// POSTの場合．
		// /tasks POSTでは，送られてきたタスク一覧にデータベースを合わせる．
	} else if r.Method == http.MethodPost {
		// リクエストに添付されたbodyを読み込む
		// まずはByte列に変換
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("Can't catch Tasks(io error)", err)
			return
		}

		// bodyを受け取るための変数を宣言．
		var recTasks []dbctl.Task

		// json.Unmarshalでタスク一覧の配列をrecTasksに代入する．
		if err := json.Unmarshal(jsonBytes, &recTasks); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("Can't catch Tasks(JSON Unmarshal error)", err)
			return
		}
		// この時点で，recTasksにbodyを代入する作業は完了している．
		// log.Print(recTasks[0].Name) // は問題なく実行できる．

		// ここからは，クライアントのタスク一覧とデータベースのタスク一覧を一致させていく．
		// データベースのタスク一覧を受け取るための変数を宣言．
		var dbTasks []dbctl.Task

		// そして，データベースから受け取る．（エラー処理も忘れずに．）
		// GetTasksの詳しい処理については，dbctl/tasks.goにて．
		dbTasks, err = dbctl.GetTasks()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("database error(GetTasks)", err)
			return
		}

		// insertTasksには，「クライアントにはあって，データベースに無いタスク」を入れていく．
		// つまり，insertTasksは「データベースにinsertしなければならないTask一覧」となる．
		var insertTasks []dbctl.Task

		// 全てのクライアントのタスクに対して，データベースに存在するかどうかを確かめる．
		for _, recTask := range recTasks {
			var findFlag = false
			for _, dbTask := range dbTasks {
				if recTask.ID == dbTask.ID {
					findFlag = true
					break
				}
			}
			if !findFlag {
				// データベースに存在しなかったら，insertTasksに追加．
				insertTasks = append(insertTasks, recTask)
			}
		}

		// もしinsertしなければならないタスクが無ければ何もしない．
		if insertTasks != nil {
			// あれば，insertする．（エラー処理も忘れずに．）
			// InsertTasksについては，dbctl/tasks.goを参照．
			err := dbctl.InsertTasks(insertTasks)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintln(w, `{"status":"Unavailable"}`)
				fmt.Println("database error(InsertTasks)", err)
				return
			}
		}

		// deleteIDsは，「クライアントには無いが，データベースにはあるタスクのID一覧」である．
		// つまり，削除しなければならないタスクのID一覧ということである．
		var deleteTaskIDs []int

		// 全てのデータベースのタスクについて，クライアントに存在するかどうかを確かめる．
		for _, dbTask := range dbTasks {
			var findFlag = false
			for _, recTask := range recTasks {
				if recTask.ID == dbTask.ID {
					findFlag = true
					break
				}
			}
			if !findFlag {
				// クライアントに存在しなかったら，deleteTaskIDsに追加．
				deleteTaskIDs = append(deleteTaskIDs, dbTask.ID)
			}
		}

		// もしデータベースから消さなければならないタスクが無ければ何もしない．
		if deleteTaskIDs != nil {
			// あれば，消す．（エラー処理も忘れずに．）
			// DeleteTasksについては，dbctl/tasks.goを参照．
			err := dbctl.DeleteTasks(deleteTaskIDs)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintln(w, `{"status":"Unavailable"}`)
				fmt.Println("database error(DeleteTasks)", err)
				return
			}
		}

		// 処理が成功したことをヘッダーに書き込む
		w.WriteHeader(http.StatusOK)

		// フロントが成功したことをJSONとして認識できるようにこっちにも書き込む．
		fmt.Fprintln(w, `{"status":"Available"}`)

		// 余談だが，別に成功したかしていないかをFprintlnでhttp.ResponseWriterに書き込む必要はない．
		// w.WriteHeader(http.StatusOK)で成功失敗の判断はつくからである．

		// じゃあなぜ書いているのかといえば　まあ　前回のチームのルールがそういうルールだったからというだけの理由である．
	}
}

/*
余談だが，長すぎるif文の中身というのは，本来あまり行儀の良い書き方ではない．
r.Method == http.MethodPostにおけるPostの処理は長いため，
別の関数に分けたり別のファイルに分けたりして，よりわかりやすくするべきである．

では，なぜ本サンプルでは分けなかったか．
それは，「ifの中身，コメント抜いたら100行も無かったし，別にまあ良いか．」と思ったからである．

結局のところ，最後はいつだって，書く人の裁量に全て委ねられるのである．（あとレビュアー）
*/
