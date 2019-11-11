package mclient

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Client stucture
type Client struct {
	gorm.Model
	ClientID      uint      `gorm:"column:client_id;not null;unique_index;" json:"ClientID"`
	FirstName     string    `gorm:"column:first_name;type:varchar(100);" json:"FirstName"`
	LastName      string    `gorm:"column:last_name;type:varchar(100);" json:"LastName"`
	Phone         string    `gorm:"column:phone;type:varchar(50); DEFAULT:null" json:"Phone"`
	Gender        string    `gorm:"column:gender;type:enum('male','female'); DEFAULT:null" json:"Gender"`
	DocumentID    string    `gorm:"column:document_id;type:varchar(255); DEFAULT:null" json:"DocumentID"`
	DateOfBirth   time.Time `gorm:"column:date_of_birth;type:datetime; DEFAULT:null" json:"DateOfBirth"`
	Comment       string    `gorm:"column:comment;type:varchar(255);" json:"Comment"`
	IsLegalEntity bool      `gorm:"column:is_legal_entity;not null; default:0" json:"IsLegalEntity"`
	BusinessID    string    `gorm:"column:business_id;type:varchar(255); DEFAULT:null" json:"BusinessID"`
	LegalAddress  string    `gorm:"column:legal_address;type:varchar(255); DEFAULT:null" json:"LegalAddress"`
	Timezone      string    `gorm:"column:timezone;type:varchar(255); DEFAULT:null" json:"Timezone"`
}

//TableName return new table name
func (Client) TableName() string {
	return "s_clients"
}

//GetByClientID возвращаем клинета по ClientID
func (cl *Client) GetByClientID(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Where("client_id = ?", cl.ClientID).First(&cl)
	if r.Error != nil && !strings.Contains(r.Error.Error(), "not found") {
		return r.Error
	}

	return nil
}

//Post создание нового клиента
func (cl *Client) Post(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	var r *gorm.DB
	r = DB.Create(&cl)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put обновление информации о клиенте
func (cl *Client) Put(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	//Create main record
	var r *gorm.DB
	r = DB.Model(&cl).Where("id = ?", cl.ID).Updates(&cl)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
