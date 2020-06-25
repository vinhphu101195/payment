package database

import (
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
	db.username = "root"
	db.password = "123456"
	db.url = "167.71.195.5"
	db.port = "3306"
	db.database = "VirtualGiftSystem"
}

//GetDB ...
func GetDB() (*gorm.DB, error) {
	var db DBInfo
	db.setDB()
	return gorm.Open("mysql", db.username+":"+db.password+"@("+db.url+":"+db.port+")/"+db.database+"?charset=utf8&parseTime=True&loc=Local")
}
