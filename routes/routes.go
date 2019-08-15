package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/controllers/citem"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/controllers/cuser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/muser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/database"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
)

//InitRoutes инициализирует пути
func InitRoutes(router *gin.Engine) *gin.Engine {
	//THIS IS HTTP SERVER IMPLEMENTATION FOR FRONTEND
	router.NoRoute(HTTPStaticServer)

	router.StaticFS("/resources", http.Dir("resources"))

	g := router.Group("/api/v1/")

	g.GET("/", func(c *gin.Context) {
		c.JSON(200, "VKAEP API works in normal mode")
	})

	//Незащищенные пути
	g.Use(MainMiddleware())
	{
		g.POST("/user/login", cuser.Login)
		g.GET("/user/login", cuser.Login)
	}

	//Пути защищенные авторизацией
	g.Use(MainMiddleware())
	g.Use(AuthMiddleware())
	{
		g.POST("/items/invoices/upload-xls", citem.UploadExciseXLS)
		g.GET("/items/invoices", citem.GetItemInvoices)
		g.POST("/items/scanned/upload-xls", citem.GetItemInvoices)
		g.POST("/items/scanned/report", citem.GetItemScannedReport)
		g.POST("/items/scanned/download-xls", citem.DownloadXLS)
		g.POST("/item/bonus", citem.GetItemBonus)
	}

	return router
}

//MainMiddleware Open GLOBAL connection to Main Database, LOGIN and SIGN UP uses MAIN base for authentication
func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var DB *gorm.DB

		DB, err = database.OpenDatabase("vkaep", "vkaep", settings.DBRTUP, "", "UTC")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Ошибка при подключении к БД"})
			return
		}
		defer DB.Close()
		//Set DB into context
		c.Set("DB", DB)

		//Set date in header
		currentDate := time.Now()
		c.Header("X-Server-Date", currentDate.Format("2006-01-02"))
		c.Request.Header.Set("X-Server-Date", currentDate.Format("2006-01-02"))

		c.Next()
	}
}

//AuthMiddleware set guard on routes wich needs the authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var user muser.User

		//Получение токена из Headers
		user.Token, _ = cuser.GetToken(c)
		if user.Token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": "Пользователь не авторизован"})
			return
		}

		//Попытка авторизации по X-Token
		_, err = cuser.AuthToken(c, user, "token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": err.Error()})
			return
		}

		c.Next()
	}
}
