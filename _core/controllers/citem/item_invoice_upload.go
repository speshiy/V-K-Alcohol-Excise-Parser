package citem

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

//IncomeInvoiceData структура
type IncomeInvoiceData struct {
	ItemName              string  `json:"ItemName"`
	ItemType              string  `json:"ItemType"`
	ItemVolume            float32 `json:"ItemVolume"`
	ItemMarkType          string  `json:"ItemMarkType"`
	ItemSerial            string  `json:"ItemSerial"`
	ItemBeginExciseNumber string  `json:"ItemBeginExciseNumber"`
	ItemEndExciseNumber   string  `json:"ItemEndExciseNumber"`
	ItemBonus             float32 `json:"ItemBonus"`
}

//UploadExciseXLS загружает JSON с алкогольной продукцией и вставляет их в БД
func UploadExciseXLS(c *gin.Context) {
	var err error
	var items []IncomeInvoiceData

	if err = c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Валидируем входящие данные
	err = validate(c, &items)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	for idx, item := range items {
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
				// c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Ошибка при вставке товара в БД, cтрока: " + strconv.Itoa(idx+1) +
				// 	" - попытка вставить существующий диапазон акцизов"})
				continue
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Ошибка при вставке товара в БД, cтрока: " + strconv.Itoa(idx+1)})
			}
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Загрузка данных завершена"})
}

func validate(c *gin.Context, items *[]IncomeInvoiceData) error {
	var err error

	for idx, item := range *items {
		err = common.Validate.Var(item.ItemName, "required")
		if err != nil {
			return errors.New("Название продукции не заполнено - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemType, "required")
		if err != nil {
			return errors.New("Вид продукции не заполнен - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemVolume, "required")
		if err != nil {
			return errors.New("Емкость продукции не заполнена - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemMarkType, "required")
		if err != nil {
			return errors.New("Тип марки не заполнен - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemSerial, "required")
		if err != nil {
			return errors.New("Серия не заполнена - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemBeginExciseNumber, "required")
		if err != nil {
			return errors.New("Начальный номер не заполнен - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemEndExciseNumber, "required")
		if err != nil {
			return errors.New("Конечный номер не заполнен - строка: " + strconv.Itoa(idx+1))
		}

		n1, _ := strconv.Atoi(item.ItemBeginExciseNumber)
		n2, _ := strconv.Atoi(item.ItemEndExciseNumber)
		if n1 > n2 {
			return errors.New("Начальный номер акцизов больше чем конечный - строка: " + strconv.Itoa(idx+1))
		}

		err = common.Validate.Var(item.ItemBonus, "required")
		if err != nil {
			return errors.New("Бонусы не заполнены - строка: " + strconv.Itoa(idx+1))
		}

		if item.ItemBonus < 0 {
			return errors.New("Бонус не может быть меньше нуля - строка: " + strconv.Itoa(idx+1))
		}
	}

	return nil
}
