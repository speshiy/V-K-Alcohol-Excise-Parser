package database

import (
	"errors"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
)

//CreateUserRTU create rtu user
func CreateUserRTU(host string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "root:"+settings.DBRP+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	//host 192.168.%
	err = db.Exec("CREATE USER IF NOT EXISTS 'rtu'@'%' IDENTIFIED BY '" + settings.DBRTUP + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT ALL PRIVILEGES ON *.* TO 'rtu'@'%' WITH GRANT OPTION").Error
	if err != nil {
		return err
	}

	return nil
}

//CreateUser create user
func CreateUser(isDemo bool, host string, databaseName string, username string, password string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	//host 192.168.%
	err = db.Exec("CREATE USER IF NOT EXISTS '" + username + "'@'192.168.%' IDENTIFIED BY '" + password + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT EXECUTE, CREATE, ALTER, DROP, INDEX, LOCK TABLES, SELECT, INSERT, UPDATE, DELETE ON `" + databaseName + "`.* TO '" + username + "'@'192.168.%'").Error
	if err != nil {
		return err
	}

	//host 127.0.0.1
	err = db.Exec("CREATE USER IF NOT EXISTS '" + username + "'@'127.0.0.1' IDENTIFIED BY '" + password + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT EXECUTE, CREATE, ALTER, DROP, INDEX, LOCK TABLES, SELECT, INSERT, UPDATE, DELETE ON `" + databaseName + "`.* TO '" + username + "'@'127.0.0.1'").Error
	if err != nil {
		return err
	}

	return nil
}

//CreateDatabase if not exist
func CreateDatabase(host string, databaseName string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+host+":3306)/")
	if err != nil {
		if strings.Contains(err.Error(), "1045") {
			err = CreateUserRTU(host)
			if err != nil {
				return err
			}

			db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+host+":3306)/")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer db.Close()

	err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + databaseName + "` CHARACTER SET utf8 COLLATE utf8_general_ci").Error
	if err != nil {
		return err
	}

	return err
}

//DropDatabase drop database
func DropDatabase(host string, databaseName string, username string) error {
	var err error
	var r *gorm.DB
	var db *gorm.DB
	db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	sql := "DROP DATABASE IF EXISTS `" + databaseName + "`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}
	sql = "DROP USER IF EXISTS '" + username + "'@`192.168.%`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}
	sql = "DROP USER IF EXISTS '" + username + "'@`127.0.0.1`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//TryOpenDatabase попытка открыть БД
func TryOpenDatabase(host string, databaseName string, username string, password string, location string) (*gorm.DB, error) {
	if location == "" {
		location = "UTC"
	} else {
		location = strings.Replace(location, "/", "%2F", -1)
	}
	return gorm.Open("mysql", ""+username+":"+password+"@tcp("+host+":3306)/"+databaseName+"?charset=utf8&parseTime=True&loc="+location)
}

//OpenDatabase opens database
func OpenDatabase(databaseName string, username string, password string, host string, location string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	host = settings.DBHostDefault

	if len(host) == 0 {
		return nil, errors.New("Host not set for DB while OpenDatabase " + databaseName)
	}

	db, err = TryOpenDatabase(host, databaseName, username, password, location)
	if err != nil {
		//if user not found then create user and try to open
		if strings.Contains(err.Error(), "1045") {
			err = CreateDatabase(host, databaseName)
			if err != nil {
				return nil, err
			}

			err = CreateUser(false, host, databaseName, username, password)
			if err != nil {
				return nil, err
			}

			db, err = TryOpenDatabase(host, databaseName, username, password, location)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		//if database not found then create DB
		if strings.Contains(err.Error(), "1049") {
			err = CreateDatabase(host, databaseName)
			if err != nil {
				return nil, err
			}

			db, err = TryOpenDatabase(host, databaseName, username, password, location)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		log.Print("Connection to DB in service FAILED.", err.Error())
		return nil, err
	}
	return db, nil
}
