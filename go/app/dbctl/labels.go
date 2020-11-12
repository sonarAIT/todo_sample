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
	rows, err := db.Query("select id, label_text from labels")
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.Printf(errFormat, err, f.Name(), file, line)
		return nil, err
	}
	defer rows.Close()

	var Labels []Label

	for rows.Next() {
		var id int
		var labelText string
		rows.Scan(&id, &labelText)
		Labels = append(Labels, Label{ID: id, LabelText: labelText})
	}

	return Labels, nil
}
