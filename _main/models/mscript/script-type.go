package mscript

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ScriptType structure
type ScriptType struct {
	gorm.Model
	CompanyID uint   `gorm:"column:company_id;" json:"CompanyID"`
	Name      string `gorm:"column:name;" json:"Name"`
	Code      string `gorm:"column:code;default: null;" json:"Code"`
}

//TableName return new table name
func (ScriptType) TableName() string {
	return "s_script_type"
}

//GetScriptTypes return list of blocks
func GetScriptTypes(c *gin.Context, companyID uint, scriptTypes *[]ScriptType) error {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB

	r = DB.Where("company_id = ?", companyID).Find(&scriptTypes)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Post record
func (st *ScriptType) Post(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Create(&st)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put record
func (st *ScriptType) Put(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Model(&st).Where("id = ?", st.ID).Updates(&st)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Delete delete card into DB
func (st *ScriptType) Delete(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Delete(&st)

	if r.Error != nil {
		return r.Error
	}

	return nil
}
