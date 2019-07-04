package mitem

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ItemInvoice structure
type ItemInvoice struct {
	gorm.Model
	ItemName              string  `gorm:"column:item_name;type:varchar(255);default:null" json:"ItemName"`
	ItemType              string  `gorm:"column:item_type;type:varchar(255);default:null" json:"ItemType"`
	ItemVolume            float32 `gorm:"column:item_volume;type:decimal(19,3);default:null" json:"ItemVolume"`
	ItemMarkType          string  `gorm:"column:item_mark_type;type:varchar(255);default:null" json:"ItemMarkType"`
	ItemSerial            string  `gorm:"column:item_serial;type:varchar(255);default:null" json:"ItemSerial"`
	ItemBeginExciseNumber uint    `gorm:"column:item_begin_excise_number;default:null;unique_index;" json:"ItemBeginExciseNumber"`
	ItemEndExciseNumber   uint    `gorm:"column:item_end_excise_number;default:null;unique_index;" json:"ItemEndExciseNumber"`
	ItemBonus             float32 `gorm:"column:item_bonus;type:decimal(19,3);default:null" json:"ItemBonus"`
}

//TableName return new table name for User struct
func (ItemInvoice) TableName() string {
	return "s_items_invoice"
}

//GetItemInvoices возвращает загруженные накладные
func GetItemInvoices(c *gin.Context, DB *gorm.DB, ii *[]ItemInvoice) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	var r *gorm.DB
	r = DB.Order("id DESC").Find(&ii)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetByExciseRange возвращаем накладную по акцизу в промежутке между цифрами
func (i *ItemInvoice) GetByExciseRange(c *gin.Context, DB *gorm.DB, excise uint) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Where("? BETWEEN item_begin_excise_number AND item_end_excise_number", excise).First(&i)
	if r.Error != nil && !strings.Contains(r.Error.Error(), "not found") {
		return r.Error
	}

	return nil
}

//Post new item into DB
func (i *ItemInvoice) Post(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Create(&i)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put update item into DB
func (i *ItemInvoice) Put(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Model(&i).Where("id = ?", i.ID).Updates(&i)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
