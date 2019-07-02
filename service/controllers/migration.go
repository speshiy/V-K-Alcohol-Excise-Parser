package controllers

import (
	"github.com/gin-gonic/gin"
)

//Migrate all bases
func Migrate(c *gin.Context) {
	MigrateVkaep(c)
}
