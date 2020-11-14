package dbctl

import (
	"database/sql"
	"log"
	"runtime"
)

// errFormatはデータベースのエラー表示に使用されるフォーマット
const errFormat = "%v\nfunction:%v file:%v line:%v\n"

var db *sql.DB

// init関数は特殊な関数
// dbctlパッケージが読み込まれた瞬間に実行される．
func init() {
	var err error

	// データベースを開く
	db, err = sql.Open("mysql", "gopher:password@tcp(todo_sample_mysql:3306)/prod_db")
	// データベースを開けなかった時のエラー
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.Printf(errFormat, err, f.Name(), file, line)

		// データベースが使えなかったら処理もクソも無いので，panicで強制終了．
		panic("Can't Open database.")
	}
}
