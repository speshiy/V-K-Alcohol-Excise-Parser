package citem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mshared"
)

//GetItemScannedReport возращает отчет об отсканированных акцизах
func GetItemScannedReport(c *gin.Context) {
	var itemScannedReport []mitem.ItemScannedReport
	var err error
	var reportFilter mshared.Filter

	if err := c.ShouldBindJSON(&reportFilter); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = mitem.GetItemScannedReport(c, reportFilter, &itemScannedReport)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": itemScannedReport})
}
