package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/controllers/ccompany"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/controllers/cscript"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/controllers/cuser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/models/muser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/common"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/database"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/settings"
)

//InitRoutes инициализирует пути
func InitRoutes(router *gin.Engine) *gin.Engine {
	g := router.Group("/api/v1/")
	g.Use(MainMiddleware())
	{
		g.POST("/", func(c *gin.Context) {
			var config common.Config
			config.SetConfigLocal()
			c.JSON(200, config)
		})

		g.GET("/companies", ccompany.GetCompanies)

		g.GET("/script-types/:CompanyID", cscript.GetScriptTypes)
		g.POST("/script-type", cscript.PostScriptType)
		g.PUT("/script-type", cscript.PutScriptType)
		g.DELETE("/script-type/:ID", cscript.DeleteScriptType)

		g.GET("/blocks/:TypeID", cscript.GetScriptBlocks)
		g.POST("/block", cscript.PostScriptBlock)
		g.PUT("/block", cscript.PutScriptBlock)
		g.DELETE("/block/:ID", cscript.DeleteScriptBlock)
	}

	return router
}

//MainMiddleware Open GLOBAL connection to Main Database, LOGIN and SIGN UP uses MAIN base for authentication
func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var DB *gorm.DB

		//Set locale from header
		DB, err = database.OpenDatabase("script", "script", settings.DBRTUP, "", "UTC")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "E_CONNECT_MAIN_DB"})
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

		//Getting token from session
		user.Token, _ = cuser.GetToken(c)
		if user.Token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": "E_USER_NOT_AUTHORIZED"})
			return
		}

		//Trying to auth with token
		user, err = cuser.AuthToken(c, user, "token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": err.Error()})
			return
		}

		//Set main info about user into context
		c.Set("User", user)

		c.Next()
	}
}

//UserMiddleware set guard on routes wich needs the authorization
func UserMiddleware(checkIsSuper bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var user muser.User

		//Getting token from session
		user.Token, _ = cuser.GetToken(c)
		if user.Token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": "E_USER_NOT_AUTHORIZED"})
			return
		}

		//Trying to auth with token
		user, err = cuser.AuthToken(c, user, "token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "false", "code": "token_error", "message": err.Error()})
			return
		}

		//Set main info about user into context
		c.Set("User", user)

		c.Next()
	}
}
