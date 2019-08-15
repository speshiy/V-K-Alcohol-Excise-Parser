package citem

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mclient"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

//IncomeScannedData структура
type IncomeScannedData struct {
	Excise uint `json:"Excise"`
	Code   uint `json:"Code"`
}

//UploadItemXLS загружает JSON с отсканированными/сфотографированными товарами
func UploadItemXLS(c *gin.Context) {
	var err error
	var DB *gorm.DB
	var incomeScannedData []IncomeScannedData
	var itemScanned mitem.ItemScanned
	var itemInvoice mitem.ItemInvoice
	var client mclient.Client

	if err = c.ShouldBindJSON(&incomeScannedData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Валидируем входящие данные
	err = validateScanned(c, &incomeScannedData)
	if err != nil {
		return
	}

	//Получаем подключание к БД
	DB = c.MustGet("DB").(*gorm.DB)

	//Запускаем транзакцию
	tx := DB.Begin()

	for idx, item := range incomeScannedData {

		//Ищем такой же акциз в отсканированных акцизах для проверки
		itemScanned.ItemExcise = item.Excise
		err = itemScanned.GetByExcise(c, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		//Если акциз уже использован то возвращаем ошибку
		if itemScanned.ID > 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акциз уже использован"})
			return
		}

		//Получаем бонусы и товар из акцизных накладных
		err = itemInvoice.GetByExciseRange(c, tx, item.Excise)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		if itemInvoice.ID == 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акциз не найден среди загруженных накладных"})
			return
		}

		//Получаем клиента по ClientID
		client.ClientID = client.ClientID
		err = client.GetByClientID(c, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		//Заполняем клиента новой информацией
		client.ClientID = client.ClientID
		client.FirstName = client.FirstName
		client.LastName = client.LastName
		client.Phone = client.Phone
		client.Gender = client.Gender
		client.DocumentID = client.DocumentID
		client.DateOfBirth = client.DateOfBirth
		client.Comment = client.Comment
		client.IsLegalEntity = client.IsLegalEntity
		client.BusinessID = client.BusinessID
		client.LegalAddress = client.LegalAddress

		if client.ID > 0 {
			err = client.Put(c, tx)
		} else {
			err = client.Post(c, tx)
		}
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		//Вставляем в отсканированные акцизы новую запись
		itemScanned.ClientID = client.ID
		itemScanned.ItemName = itemInvoice.ItemName
		itemScanned.ItemType = itemInvoice.ItemType
		itemScanned.ItemVolume = itemInvoice.ItemVolume
		itemScanned.ItemMarkType = itemInvoice.ItemMarkType
		itemScanned.ItemSerial = itemInvoice.ItemSerial
		itemScanned.ItemExcise = incomeScannedData.Excise
		itemScanned.ItemBonus = itemInvoice.ItemBonus

		err = itemScanned.Post(c, nil)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Загрузка данных завершена"})
}

func validateScanned(c *gin.Context, items *[]IncomeScannedData) error {
	var err error

	for idx, item := range *items {
		err = common.Validate.Var(item.Excise, "required")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акцизный номер не заполнен - строка: " + strconv.Itoa(idx+1)})
			return err
		}

		err = common.Validate.Var(item.Code, "required")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Код клиента не заполнен - строка: " + strconv.Itoa(idx+1)})
			return err
		}

	}

	return nil
}
