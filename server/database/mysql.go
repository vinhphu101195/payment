package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DBInfo ...
type DBInfo struct {
	username string
	password string
	url      string
	port     string
	database string
}

func (db *DBInfo) setDB() {
	db.username = os.Getenv("MYSQL_USERNAME")
	db.password = os.Getenv("MYSQL_PASSWORD")
	db.url = os.Getenv("MYSQL_URL")
	db.port = os.Getenv("MYSQL_PORT")
	db.database = os.Getenv("MYSQL_DATABASE")
}

//GetDB ...
func GetDB() (*gorm.DB, error) {
	var db DBInfo
	db.setDB()
	return gorm.Open("mysql", db.username+":"+db.password+"@("+db.url+":"+db.port+")/"+db.database+"?charset=utf8&parseTime=True&loc=Local")
}
