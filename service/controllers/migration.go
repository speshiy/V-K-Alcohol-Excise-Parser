package controllers

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/database"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
)

//AutoMigrate автоматическая миграция БД при запуске
func AutoMigrate() {
	var c gin.Context

	DB, err := database.TryOpenDatabase(settings.DBHostDefault, "information_schema", "remote_user", settings.DBRTUP, "UTC")
	if err != nil {
		log.Println("Connection to DBMain in service FAILED. ", err.Error())
		return
	}
	defer DB.Close()

	databases := []string{"vkaep"}

	for _, database := range databases {
		schemaName := ""
		row := DB.Raw("SELECT schema_name as SchemaName FROM schemata WHERE schema_name = ?", database).Select("SchemaName").Row()
		err = row.Scan(&schemaName)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			log.Println("Automigrate schema select", err.Error())
			return
		}

		if schemaName == "" {
			switch database {
			case "vkaep":
				MigrateVkaep(&c)
			}
		}
	}
}

//Migrate all bases
func Migrate(c *gin.Context) {
	MigrateVkaep(c)
}
