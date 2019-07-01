package mcompany

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Company structure
type Company struct {
	gorm.Model
	Name string `gorm:"column:name;unique_index" json:"Name"`
}

//TableName return new table name
func (Company) TableName() string {
	return "s_company"
}

//GetCompanies return list of companies
func GetCompanies(c *gin.Context, companies *[]Company) error {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB

	r = DB.Find(&companies)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Post record
func (com *Company) Post(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Create(&com)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put record
func (com *Company) Put(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Model(&com).Where("id = ?", com.ID).Updates(&com)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Delete delete card into DB
func (com *Company) Delete(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Delete(&com)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//PreFill table with new data
func (com Company) PreFill(db *gorm.DB) error {
	for _, sql := range com.GetRefillSQL() {
		r := db.Exec(sql)
		if r.Error != nil && !strings.Contains(r.Error.Error(), "Duplicate") {
			return r.Error
		}
	}

	return nil
}

//GetRefillSQL sql for clean and fill table with new data
func (Company) GetRefillSQL() []string {
	return []string{
		"INSERT INTO s_company(id, name) values (1, 'TuviS')",
		"INSERT INTO s_company(id, name) values (2, 'Paloma365')",
	}
}
