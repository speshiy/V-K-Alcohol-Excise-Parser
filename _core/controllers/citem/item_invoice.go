package citem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
)

//GetItemInvoices возращает список акцизных накладных
func GetItemInvoices(c *gin.Context) {
	var err error
	var itemInvoices []mitem.ItemInvoice

	err = mitem.GetItemInvoices(c, nil, &itemInvoices)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": itemInvoices})
}
