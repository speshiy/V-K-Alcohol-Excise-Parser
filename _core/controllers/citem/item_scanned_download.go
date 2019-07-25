package citem

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mshared"
	"github.com/tealeg/xlsx"
)

//DownloadXLS скачивает отсканированные акцизные кода
func DownloadXLS(c *gin.Context) {
	var err error
	var reportFilter mshared.Filter
	var itemScannedReport []mitem.ItemScannedReport

	if err := c.ShouldBindJSON(&reportFilter); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = mshared.ValidateFilterDates(c, &reportFilter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = mitem.GetItemScannedReport(c, reportFilter, &itemScannedReport)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Формируем xls файл
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	columnsName := []string{"Название алкогольной продукции", "Вид алкогольной продукции",
		"Емкость", "Тип марки", "Серия", "Акцизный номер", "Значение бонусов",
		"Имя", "Фамилия", "Телефон", "Пол", "ИИН", "День рождения", "БИН", "Юр. Адрес", "Отсканирован"}

	//Создаем файл xls
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Отсканированные акцизы")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Заполняем шапку
	row = sheet.AddRow()
	for _, name := range columnsName {
		cell = row.AddCell()
		cell.Value = name
	}

	//Заполняем таблицу
	for _, item := range itemScannedReport {
		row = sheet.AddRow()

		cell = row.AddCell()
		cell.Value = item.ItemName

		cell = row.AddCell()
		cell.Value = item.ItemType

		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", item.ItemVolume)

		cell = row.AddCell()
		cell.Value = item.ItemMarkType

		cell = row.AddCell()
		cell.Value = item.ItemSerial

		cell = row.AddCell()
		cell.Value = strconv.Itoa(int(item.ItemExcise))

		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", item.ItemBonus)

		cell = row.AddCell()
		cell.Value = item.FirstName

		cell = row.AddCell()
		cell.Value = item.LastName

		cell = row.AddCell()
		cell.Value = item.Phone

		cell = row.AddCell()
		cell.Value = item.Gender

		cell = row.AddCell()
		cell.Value = item.DocumentID

		cell = row.AddCell()
		cell.Value = item.DateOfBirth.Format("02-01-2006")

		cell = row.AddCell()
		cell.Value = item.BusinessID

		cell = row.AddCell()
		cell.Value = item.LegalAddress

		cell = row.AddCell()
		cell.Value = item.CreatedAt.Format("02-01-2006 15:04")
	}

	err = os.MkdirAll("resources", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	uid, _ := common.GetNewUUID()
	fileName := "resources/" + uid + ".xlsx"
	err = file.Save(fileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": fileName})
}
