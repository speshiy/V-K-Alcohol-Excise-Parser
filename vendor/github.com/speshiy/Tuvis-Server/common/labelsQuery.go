package common

import (
	"github.com/gin-gonic/gin"
)

//GetCountLabel return
func GetCountLabel(c *gin.Context) string {
	return Translate("QUANTITY_LABEL", nil, GetLocale(c))
}

//GetAmountLabel return
func GetAmountLabel(c *gin.Context) string {
	return Translate("TOTAL_LABEL", nil, GetLocale(c))
}

//GetMaleLabel return
func GetMaleLabel(c *gin.Context) string {
	return Translate("MALE_LABEL", nil, GetLocale(c))
}

//GetFemaleLabel return
func GetFemaleLabel(c *gin.Context) string {
	return Translate("FEMALE_LABEL", nil, GetLocale(c))
}
