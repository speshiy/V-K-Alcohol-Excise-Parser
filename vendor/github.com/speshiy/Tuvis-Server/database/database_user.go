package database

import (
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/speshiy/Tuvis-Server/common"
)

//DBUserConnect connect
type DBUserConnect struct {
	Host     string
	DBName   string
	User     string
	Password string
	Connect  *gorm.DB
}

//Open connection
func (db *DBUserConnect) Open(host string, dbname string, user string, password string, location string) error {
	var err error

	if location == "" {
		location = "UTC"
	} else {
		location = strings.Replace(location, "/", "%2F", -1)
	}

	db.Host = host
	db.DBName = dbname
	db.User = user
	db.Password = password
	db.Connect, err = OpenDatabase(db.DBName, db.User, db.Password, db.Host, location)
	//db.Connect, err =  gorm.Open("mysql", db.User+":"+db.Password+"@tcp("+db.Host+":3306)/"+db.DBName+"?charset=utf8&parseTime=True&loc="+location)
	if err != nil {
		log.Print("Connection to user DB FAILED.", err.Error())
		return err
	}

	return nil
}

//Close opens Main Database
func (db *DBUserConnect) Close() error {
	err := db.Connect.Close()
	if err != nil {
		log.Print("Closing user DB FAILED.", err.Error())
		return err
	}

	return nil
}

//Active return true if connection is open
func (db *DBUserConnect) Active(err *common.ErrorTable) bool {
	return db.Connect.DB().Ping() != nil
}

//LogError log error in DB
func (db *DBUserConnect) LogError(err common.ErrorTable) {
	db.Connect.Create(&err)
}
