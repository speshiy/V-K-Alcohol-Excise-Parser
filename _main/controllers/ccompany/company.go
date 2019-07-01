package ccompany

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/models/mcompany"
)

//GetCompanies return list of companies
func GetCompanies(c *gin.Context) {
	var err error
	var companies []mcompany.Company

	err = mcompany.GetCompanies(c, &companies)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": companies})
}
