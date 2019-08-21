package citem

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mclient"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
)

//IncomeScannedData структура
type IncomeScannedData struct {
	ExciseNumber uint   `json:"ExciseNumber"`
	ClientCode   string `json:"ClientCode"`
}

//SetBonus начисляет бонус в TuviS
func SetBonus(c *gin.Context, item *IncomeScannedData, idx int, xTokenAPI string) error {
	var err error
	var client mclient.Client
	var itemInvoice mitem.ItemInvoice
	var itemScanned mitem.ItemScanned

	//Ищем такой же акциз в отсканированных акцизах для проверки
	itemScanned.ItemExcise = item.ExciseNumber
	err = itemScanned.GetByExcise(c, nil)
	if err != nil {
		return err
	}

	//Если за акциз было начисление то переходим к следующему акцизу
	if itemScanned.ID > 0 {
		return nil
	}

	//Получаем бонусы и товар из акцизных накладных
	err = itemInvoice.GetByExciseRange(c, nil, item.ExciseNumber)
	if err != nil {
		return err
	}

	if itemInvoice.ID == 0 {
		return errors.New("Акциз не найден среди загруженных накладных - строка: " + strconv.Itoa(idx+1))
	}

	//Post запрос в tuvis.world для начисления бонусов и получения информации о потребителе
	responseData := TuvisManulBonusResponse{}

	request := resty.New()
	resp, err := request.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Token-API", xTokenAPI).
		SetBody(ManualBonusType{CardID: 1, ClientCode: item.ClientCode, Points: itemInvoice.ItemBonus, TransactionType: "sale"}).
		SetResult(&responseData).
		SetError(&responseData).
		Post(settings.TuviSHost + "/card/bonus/transaction/vkaep")

	if err != nil {
		return err
	}

	//Если запрос верулся с ошибкой
	if !resp.IsSuccess() {
		return errors.New("TuviS: " + responseData.TuvisResponse.Message + " - строка: " + strconv.Itoa(idx+1))
	}

	//Если при начислении бонусов произошла ошибка
	if responseData.TuvisResponse.Status == "false" {
		return errors.New("TuviS: " + responseData.TuvisResponse.Message + " - строка: " + strconv.Itoa(idx+1))
	}

	//Получаем клиента по ClientID
	client.ClientID = responseData.Client.ID
	err = client.GetByClientID(c, nil)
	if err != nil {
		return err
	}

	//Заполняем клиента новой информацией
	client.ClientID = responseData.Client.ID
	client.FirstName = responseData.Client.FirstName
	client.LastName = responseData.Client.LastName
	client.Phone = responseData.Client.Phone
	client.Gender = responseData.Client.Gender
	client.DocumentID = responseData.Client.DocumentID
	client.DateOfBirth = responseData.Client.DateOfBirth
	client.Comment = responseData.Client.Comment
	client.IsLegalEntity = responseData.Client.IsLegalEntity
	client.BusinessID = responseData.Client.BusinessID
	client.LegalAddress = responseData.Client.LegalAddress

	if client.ID > 0 {
		err = client.Put(c, nil)
	} else {
		err = client.Post(c, nil)
	}
	if err != nil {
		return err
	}

	//Вставляем в отсканированные акцизы новую запись
	itemScanned.ClientID = client.ID
	itemScanned.ItemName = itemInvoice.ItemName
	itemScanned.ItemType = itemInvoice.ItemType
	itemScanned.ItemVolume = itemInvoice.ItemVolume
	itemScanned.ItemMarkType = itemInvoice.ItemMarkType
	itemScanned.ItemSerial = itemInvoice.ItemSerial
	itemScanned.ItemExcise = item.ExciseNumber
	itemScanned.ItemBonus = itemInvoice.ItemBonus

	err = itemScanned.Post(c, nil)
	if err != nil {
		return err
	}

	return nil
}

//IncomeScanData данные которые приходят от TuviS после сканирования
type IncomeScanData struct {
	Client mclient.Client `json:"Client"`
	Excise uint           `json:"Excise"`
}

//GetItemBonus находим акцизный номер в загруженных накладных и возвращаем бонусы в TuviS для начисления
func GetItemBonus(c *gin.Context) {
	var err error
	var incomeScannedData IncomeScanData
	var itemScanned mitem.ItemScanned
	var itemInvoice mitem.ItemInvoice
	var client mclient.Client

	if err = c.ShouldBindJSON(&incomeScannedData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Валидируем обязательные данные
	err = common.Validate.Var(incomeScannedData.Excise, "required")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акциз не заполнен"})
		return
	}

	err = common.Validate.Var(incomeScannedData.Client.ClientID, "required")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "ClientID не заполнен"})
		return
	}

	err = common.Validate.Var(incomeScannedData.Client.Phone, "required")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Телефон не заполнен"})
		return
	}

	if !common.IsValidPhone(&incomeScannedData.Client.Phone) {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Телефон неверный"})
		return
	}

	//Ищем такой же акциз в отсканированных акцизах для проверки
	itemScanned.ItemExcise = incomeScannedData.Excise
	err = itemScanned.GetByExcise(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Если акциз уже использован то возвращаем ошибку
	if itemScanned.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акциз уже использован"})
		return
	}

	//Получаем клиента по ClientID
	client.ClientID = incomeScannedData.Client.ClientID
	err = client.GetByClientID(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Заполняем клиента новой информацией
	client.ClientID = incomeScannedData.Client.ClientID
	client.FirstName = incomeScannedData.Client.FirstName
	client.LastName = incomeScannedData.Client.LastName
	client.Phone = incomeScannedData.Client.Phone
	client.Gender = incomeScannedData.Client.Gender
	client.DocumentID = incomeScannedData.Client.DocumentID
	client.DateOfBirth = incomeScannedData.Client.DateOfBirth
	client.Comment = incomeScannedData.Client.Comment
	client.IsLegalEntity = incomeScannedData.Client.IsLegalEntity
	client.BusinessID = incomeScannedData.Client.BusinessID
	client.LegalAddress = incomeScannedData.Client.LegalAddress

	if client.ID > 0 {
		err = client.Put(c, nil)
	} else {
		err = client.Post(c, nil)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Получаем бонусы и товар из акцизных накладных
	err = itemInvoice.GetByExciseRange(c, nil, incomeScannedData.Excise)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if itemInvoice.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Акциз не найден среди загруженных накладных"})
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
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": itemScanned.ItemBonus})
}
