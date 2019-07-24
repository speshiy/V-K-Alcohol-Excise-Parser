package common

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ErrorTable struct of errors table
type ErrorTable struct {
	gorm.Model
	Func    string `gorm:"column:func;type:varchar(100);not null"`
	Message string `gorm:"column:message;type:varchar(1000);not null"`
	Line    uint   `gorm:"column:line;not null"`
}

//TableName return new table name for User struct
func (ErrorTable) TableName() string {
	return "sys_errors"
}

//Post error into table
func (e *ErrorTable) Post(c *gin.Context, db *gorm.DB) {
	r := db.Create(&e)
	if r.Error != nil {
		log.Println(r.Error.Error())
	}
}

//TError type
type TError map[string]string

//GetError return error message by key
func GetError(c *gin.Context, key string) string {
	result := Translate(key, nil, GetLocale(c))

	return result
}

//GetNormalErrorDB return normalized SQL error
func GetNormalErrorDB(c *gin.Context, err error, model string, fieldName string) string {
	//if error contains "not found" then return msg + not found
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		if model == "user" {
			switch fieldName {
			case "email":
				return GetError(c, "E_USER_WRONG_EMAIL")
			case "phone":
				return GetError(c, "E_USER_WRONG_PHONE")
			default:
				return "Unknown DB error"
			}
		}
	}
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		if model == "all" {
			switch fieldName {
			case "name":
				return GetError(c, "E_NAME_EXISTS")
			case "phone":
				return GetError(c, "E_PHONE_EXISTS")
			case "email":
				return GetError(c, "E_EMAIL_EXISTS")
			case "phoneORemail":
				return GetError(c, "E_EMAIL_OR_PHONE_EXISTS")
			default:
				return "Unknown DB error"
			}
		}
	}

	return err.Error()
}
