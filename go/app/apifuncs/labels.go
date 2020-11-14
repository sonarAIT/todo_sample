package apifuncs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"set1.ie.aitech.ac.jp/todo_sample/dbctl"
)

//LabelsFunc は/labelsにアクセスされた際に実行される関数
func LabelsFunc(w http.ResponseWriter, r *http.Request) {
	// CORSの設定．
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// ヘッダの設定．
	r.Header.Set("Content-Type", "application/json")
	// GETの場合
	// /labels GET では，クライアントに対してラベル一覧を送信する．
	if r.Method == http.MethodGet {
		// Label一覧を取得．(dbctl/labels.goを参照．)
		Labels, err := dbctl.GetLabels()
		// エラー処理
		if err != nil {
			// ヘッダーに失敗したことを書き込む
			w.WriteHeader(http.StatusServiceUnavailable)
			// ついでに失敗したことをフロントがJSONとして認識できるように書き込む
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			// ログにも書いちゃう
			fmt.Println("database error(GetLabels)", err)
			// 終了
			return
		}

		// 取得したラベル一覧をByte列に変換
		jsonBytes, err := json.Marshal(Labels)
		// エラー処理
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{"status":"Unavailable"}`)
			fmt.Println("JSON Marshal error(Labels)", err)
			return
		}

		// 文字列に変換（送信可能になった）
		jsonString := string(jsonBytes)

		// 成功したことをヘッダーに書き込む
		w.WriteHeader(http.StatusOK)

		// Labelsがもし空なら，空配列を返す．
		if Labels == nil {
			jsonString = "[]"
		}

		// 結果を書き込んで終了．
		fmt.Fprintln(w, jsonString)
	}

}
