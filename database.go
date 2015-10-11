package htmlhouse

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (app *app) initDatabase() error {
	var err error

	app.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true", app.cfg.DBUser, app.cfg.DBPassword, app.cfg.DBHost, app.cfg.DBName))
	if err != nil {
		return err
	}

	return nil
}
