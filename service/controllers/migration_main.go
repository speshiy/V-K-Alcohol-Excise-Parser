package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mclient"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/mitem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/msys"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/muser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/database"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
)

//MigrateVkaep migrate all user models
func MigrateVkaep(c *gin.Context) {
	db, err := database.OpenDatabase("vkaep", "vkaep", settings.DBRTUP, "", "UTC")
	if err != nil {
		c.String(http.StatusOK, "Connection to DB in service FAILED. %s", err.Error())
		return
	}
	defer db.Close()

	//CREATING TABLES
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").AutoMigrate(
		&muser.User{},
		&mclient.Client{},
		&mitem.ItemScanned{},
		&mitem.ItemInvoice{},
		&msys.City{},
	)
	log.Println("Models in DB vkaep created")

	db.Model(&mitem.ItemScanned{}).AddForeignKey("client_id", "s_clients(id)", "RESTRICT", "RESTRICT")

	log.Println("Foreign key in DB vkaep created")

	var user muser.User
	err = user.PreFill(db)
	if err != nil {
		log.Println(err.Error())
	}

	var cities msys.City
	err = cities.PreFill(db)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Migration vkaep done")

	c.String(http.StatusOK, "Migration vkaep done")
}
