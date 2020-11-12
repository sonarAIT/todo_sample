package dbctl

import (
	"database/sql"
	"log"
	"runtime"
)

const errFormat = "%v\nfunction:%v file:%v line:%v\n"

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("mysql", "gopher:password@tcp(todo_sample_mysql:3306)/prod_db")
	if err != nil {
		pc, file, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.Printf(errFormat, err, f.Name(), file, line)

		panic("Can't Open database.")
	}
}
