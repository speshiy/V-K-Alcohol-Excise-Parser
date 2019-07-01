package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_core/models/mcompany"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_core/models/mscript"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_core/models/muser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/database"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/settings"
)

//MigrateScript migrate all user models
func MigrateScript(c *gin.Context) {
	db, err := database.OpenDatabase("script", "script", settings.DBRTUP, "", "UTC")
	if err != nil {
		c.String(http.StatusOK, "Connection to DB in service FAILED. %s", err.Error())
		return
	}
	defer db.Close()

	//CREATING TABLES
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").AutoMigrate(
		&muser.User{},
		&mcompany.Company{},
		&mscript.ScriptType{},
		&mscript.ScriptBlock{},
	)
	log.Println("Models in DB script created")

	db.Model(&mscript.ScriptBlock{}).AddForeignKey("script_type_id", "s_script_type(id)", "RESTRICT", "RESTRICT")

	log.Println("Foreign key in DB script created")

	var company mcompany.Company
	err = company.PreFill(db)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Company in DB script filled")

	log.Println("Migration script done")

	c.String(http.StatusOK, "Migration script done")
}
