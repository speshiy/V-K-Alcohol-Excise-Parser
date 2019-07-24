package common

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//GetName return a localized name from DB
func GetName(db *gorm.DB, table string, idField string, id uint, locale string, nameField ...string) string {
	var err error
	var name string
	var fieldName string

	if len(nameField) == 0 {
		fieldName = "name"
	} else {
		fieldName = nameField[0]
	}

	row := db.Raw("SELECT " + fieldName + " " +
		"FROM " + table + " " +
		"WHERE " + idField + " = " + strconv.Itoa(int(id)) + " AND locale = '" + locale + "'").Select("name").Row()
	err = row.Scan(&name)

	if err != nil {
		return "Localize not found for lang " + locale
	}

	return name
}
