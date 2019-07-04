package mitem

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ItemScanned structure
type ItemScanned struct {
	gorm.Model
	ClientID     uint    `gorm:"column:client_id;not_null;" json:"ClientID"`
	ItemName     string  `gorm:"column:item_name;type:varchar(255);not_null;" json:"ItemName"`
	ItemType     string  `gorm:"column:item_type;type:varchar(255);not_null;" json:"ItemType"`
	ItemVolume   float32 `gorm:"column:item_volume;type:decimal(19,3);not_null;" json:"ItemVolume"`
	ItemMarkType string  `gorm:"column:item_mark_type;type:varchar(255);not_null;" json:"ItemMarkType"`
	ItemSerial   string  `gorm:"column:item_serial;type:varchar(255);not_null;" json:"ItemSerial"`
	ItemExcise   uint    `gorm:"column:item_excise;not_null;;unique_index;" json:"ItemExcise"`
	ItemBonus    float32 `gorm:"column:item_bonus;type:decimal(19,3);not_null;" json:"ItemBonus"`
}

//TableName return new table name for User struct
func (ItemScanned) TableName() string {
	return "d_items_scanned"
}

//Post new item into DB
func (is *ItemScanned) Post(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Create(&is)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
