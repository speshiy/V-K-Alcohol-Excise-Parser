package citem

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

type incomeData struct {
	Excise string `json:"Excise"`
	Code   string `json:"Code"`
}

//UploadExciseXLS загружает JSON с алкогольной продукцией и вставляет их в БД
func UploadExciseXLS(c *gin.Context) {
	var err error
	var items []incomeData

	if err = c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Валидируем входящие данные
	err = validate(c, &items)
	if err != nil {
		return
	}

	for index, item := range items {
		newItem := mitem.ItemInvoice{}
		newItem.ItemName = item.ItemName
		newItem.ItemType = item.ItemType
		newItem.ItemVolume = item.ItemVolume
		newItem.ItemMarkType = item.ItemMarkType
		newItem.ItemSerial = item.ItemSerial
		newItem.ItemBeginExciseNumber = item.ItemBeginExciseNumber
		newItem.ItemEndExciseNumber = item.ItemEndExciseNumber
		newItem.ItemBonus = item.ItemBonus

		err = newItem.Post(c, nil)
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Ошибка при вставке товара в БД, cтрока: " + strconv.Itoa(index+1) +
					" - попытка вставить существующий диапазон акцизов"})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Ошибка при вставке товара в БД, cтрока: " + strconv.Itoa(index+1)})
			}
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Загрузка данных завершена"})
}

func validateScanned(c *gin.Context, items *[]incomeData) error {
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
