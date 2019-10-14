package citem

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mclient"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

//TuvisManulBonusResponse strunct
type TuvisManulBonusResponse struct {
	common.TuvisResponse
	Client mclient.Client `json:"data"`
}

//ManualBonusType stuct
type ManualBonusType struct {
	CardID          uint    `json:"CardID"`
	ClientID        uint    `json:"ClientID"`
	ClientCode      string  `json:"ClientCode"`
	TransactionType string  `json:"TransactionType"`
	Points          float32 `json:"Points"`
}

//UploadItemXLS загружает JSON с отсканированными/сфотографированными товарами
func UploadItemXLS(c *gin.Context) {
	var err error

	var incomeScannedData []IncomeScannedData

	if err = c.ShouldBindJSON(&incomeScannedData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Валидируем входящие данные
	err = validateScanned(c, &incomeScannedData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Получаем token чтобы обращаться на сервер TuviS
	xTokenAPI := c.Param("token")
	if len(xTokenAPI) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "X-Token-API от TuviS не заполнен"})
		return
	}

	for idx, item := range incomeScannedData {
		//Устанавливает бонус
		err = SetBonus(c, &item, idx, xTokenAPI)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Загрузка данных завершена"})
}

func validateScanned(c *gin.Context, items *[]IncomeScannedData) error {
	var err error

	for idx, item := range *items {
		err = common.Validate.Var(item.ExciseNumber, "required")
		if err != nil {
			return errors.New("Акцизный номер не заполнен - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ClientCode, "required")
		if err != nil {
			return errors.New("Код клиента не заполнен - строка: " + strconv.Itoa(idx+1))
		}
	}

	return nil
}
