package mitem

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mshared"
)

//ItemScannedReport отчет об использованных акцизных кодах
type ItemScannedReport struct {
	CreatedAt     time.Time `gorm:"column:CreatedAt;" json:"CreatedAt"`
	UpdatedAt     string    `gorm:"column:UpdatedAt;" json:"UpdatedAt"`
	ItemName      string    `gorm:"column:ItemName" json:"ItemName"`
	ItemType      string    `gorm:"column:ItemType" json:"ItemType"`
	ItemVolume    float32   `gorm:"column:ItemVolume" json:"ItemVolume"`
	ItemMarkType  string    `gorm:"column:ItemMarkType" json:"ItemMarkType"`
	ItemSerial    string    `gorm:"column:ItemSerial" json:"ItemSerial"`
	ItemExcise    string    `gorm:"column:ItemExcise" json:"ItemExcise"`
	ItemBonus     float32   `gorm:"column:ItemBonus" json:"ItemBonus"`
	FirstName     string    `gorm:"column:FirstName" json:"FirstName"`
	LastName      string    `gorm:"column:LastName" json:"LastName"`
	Phone         string    `gorm:"column:Phone" json:"Phone"`
	Gender        string    `gorm:"column:Gender" json:"Gender"`
	DocumentID    string    `gorm:"column:DocumentID" json:"DocumentID"`
	DateOfBirth   time.Time `gorm:"column:DateOfBirth" json:"DateOfBirth"`
	Comment       string    `gorm:"column:Comment" json:"Comment"`
	IsLegalEntity bool      `gorm:"column:IsLegalEntity" json:"IsLegalEntity"`
	BusinessID    string    `gorm:"column:BusinessID" json:"BusinessID"`
	LegalAddress  string    `gorm:"column:LegalAddress" json:"LegalAddress"`
	Timezone      string    `gorm:"column:Timezone" json:"Timezone"`
}

//GetItemScannedReport возращает отчет об отсканированных акцизах
func GetItemScannedReport(c *gin.Context, filter mshared.Filter, isr *[]ItemScannedReport) error {
	var r *gorm.DB
	var sql string
	DB := c.MustGet("DB").(*gorm.DB)

	sql = `SELECT
				isc.id as ID,
				isc.created_at as CreatedAt,
				isc.updated_at as UpdatedAt,
				isc.item_name as ItemName,
				isc.item_type as ItemType,
				isc.item_volume as ItemVolume,
				isc.item_mark_type as ItemMarkType,
				isc.item_serial as ItemSerial,
				isc.item_excise as ItemExcise,
				isc.item_bonus as ItemBonus,
				cl.first_name as FirstName,
				cl.last_name as LastName,
				cl.phone as Phone,
				CASE cl.gender WHEN 'male' THEN 'Мужской' WHEN 'female' THEN 'Женский' END  as Gender,
				cl.document_id as DocumentID,
				cl.date_of_birth as DateOfBirth,
				cl.comment as Comment, 
				cl.is_legal_entity as IsLegalEntity,
				cl.business_id as BusinessID,
				cl.legal_address as LegalAddress,
				cl.timezone as Timezone
			FROM
				d_items_scanned AS isc
			JOIN s_clients as cl ON isc.client_id = cl.id
			WHERE
				isc.created_at BETWEEN ? AND ?`

	r = DB.Raw(sql, filter.DateBegin, filter.DateEnd).Scan(&isr)

	if r.Error != nil && !strings.Contains(r.Error.Error(), "not found") {
		return r.Error
	}

	return nil
}
