package mscript

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ScriptBlock structure
type ScriptBlock struct {
	gorm.Model
	ScriptTypeID  uint   `gorm:"column:script_type_id;" json:"ScriptTypeID"`
	ScriptElement string `gorm:"column:script_element;type:enum('header','block');" json:"ScriptElement"`
	HREF          string `gorm:"column:href;type:varchar(500);" json:"HREF"`
	Title         string `gorm:"column:title;type:varchar(500);" json:"Title"`
	Content       string `gorm:"column:content;type:text;" json:"Content"`
	Doubt         string `gorm:"column:doubt;type:text;" json:"Doubt"`
	Position      uint   `gorm:"column:position;" json:"Position"`
}

//TableName return new table name
func (ScriptBlock) TableName() string {
	return "s_script_block"
}

//GetBlocks return list of blocks
func GetBlocks(c *gin.Context, typeID uint, blocks *[]ScriptBlock) error {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB

	r = DB.Order("position").Where("script_type_id = ?", typeID).Find(&blocks)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Post record
func (sb *ScriptBlock) Post(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Create(&sb)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put record
func (sb *ScriptBlock) Put(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Model(&sb).Where("id = ?", sb.ID).Updates(&sb)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Delete delete card into DB
func (sb *ScriptBlock) Delete(c *gin.Context, DB *gorm.DB) error {
	var r *gorm.DB
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r = DB.Delete(&sb)

	if r.Error != nil {
		return r.Error
	}

	return nil
}
