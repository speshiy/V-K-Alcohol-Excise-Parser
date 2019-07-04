package mshared

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

//Filter struct
type Filter struct {
	DateBegin time.Time `json:"DateBegin"`
	DateEnd   time.Time `json:"DateEnd"`
}

//ValidateFilterDates validateincome data
func ValidateFilterDates(c *gin.Context, filter *Filter) error {
	var err error
	//addOne day to dateEnd
	filter.DateEnd = filter.DateEnd.AddDate(0, 0, 1)

	err = common.Validate.Var(filter.DateBegin, "required")
	if err != nil {
		return errors.New("Дата начала периода неверная")
	}

	err = common.Validate.Var(filter.DateEnd, "required")
	if err != nil {
		return errors.New("Дата окончания периода неверная")
	}

	return nil
}
