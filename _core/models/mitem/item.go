package mitem

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Item structure
type Item struct {
	gorm.Model
	ItemName             string    `gorm:"column:item_name;type:varchar(255);default:null" json:"ItemName"`
	ItemType             string    `gorm:"column:item_type;type:varchar(255);default:null" json:"ItemType"`
	ItemVolume           float32   `gorm:"column:item_volume;type:decimal(19,3);default:null" json:"ItemVolume"`
	ItemMarkType         string    `gorm:"column:item_mark_type;type:varchar(255);default:null" json:"ItemMarkType"`
	ItemSerial           string    `gorm:"column:item_serial;type:varchar(255);default:null" json:"ItemSerial"`
	ItemExcise           uint      `gorm:"column:item_excise;default:null;unique_index;" json:"ItemExcise"`
	ItemBonus            float32   `gorm:"column:item_bonus;type:decimal(19,3);default:null" json:"ItemBonus"`
	IsUsed               bool      `gorm:"column:is_used;default:0" json:"IsUsed"`
	IsUsedAt             time.Time `gorm:"column:is_used_at;type:datetime;default:null" json:"IsUsedAt"`
	RecipientName        string    `gorm:"column:recipient_name;type:varchar(255);default:null" json:"RecipientName"`
	RecipientDocumentID  string    `gorm:"column:recipient_document_id;type:varchar(255);default:null" json:"RecipientDocumentID"`
	RecipientPhone       string    `gorm:"column:recipient_phone;type:varchar(255);default:null" json:"RecipientPhone"`
	RecipientDateOfBirth time.Time `gorm:"column:recipient_date_of_birth;type:datetime;default:null" json:"RecipientDateOfBirth"`
}

//TableName return new table name for User struct
func (Item) TableName() string {
	return "s_items"
}

//Post new item into DB
func (i *Item) Post(c *gin.Context, DB *gorm.DB) error {
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
func (i *Item) Put(c *gin.Context, DB *gorm.DB) error {
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
