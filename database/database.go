package database

import (
	"Mirror/tools"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func Connect(login string, password string, host string, name string) *gorm.DB {
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", login, password, host, name)
	db, err := gorm.Open("mysql", dbLink)
	tools.CheckErr(err)
	db.Exec(`SET NAMES UTF8`)

	return db
}
