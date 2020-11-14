package dbctl

import (
	"log"
	"runtime"
)

// Label はラベル1つ分の構造体
type Label struct {
	ID        int    `json:"id"`
	LabelText string `json:"labelText"`
}

// GetLabels は全ラベルを取得する関数
func GetLabels() ([]Label, error) {
	// SQLを実行し，ラベルの数だけIDとlabel_textを取得する．
	rows, err := db.Query("select id, label_text from labels")
	// エラー処理も行う．
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		// フォーマットに合わせてエラー，関数の名前，ファイル名，何行目かを表示している．
		// errFormatについては，dbctl/util.goを参照．
		log.Printf(errFormat, err, f.Name(), file, line)
		// 無のラベル一覧とエラーをLabelsFuncに返す．(apifuncs/labels.go)
		return nil, err
	}
	// あとでrows.Close()を行う．（この行は書き忘れないように！）
	defer rows.Close()

	// Labelの配列を宣言．
	var Labels []Label

	// SQLで取得した行数だけ実行
	for rows.Next() {
		// SQLの結果一行分を受け取る変数を宣言
		var id int
		var labelText string
		// 受け取る
		rows.Scan(&id, &labelText)
		// 受け取ったものをLabelsに追加する．
		Labels = append(Labels, Label{ID: id, LabelText: labelText})
	}

	// Labelsとエラー(無)を返す．
	return Labels, nil
}
