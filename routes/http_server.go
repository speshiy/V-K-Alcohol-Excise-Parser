package routes

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//HTTPStaticServer return frontend application
func HTTPStaticServer(c *gin.Context) {
	var app string
	var method string
	var uri string

	method = c.Request.Method
	uri = c.Request.RequestURI

	if method != "GET" {
		c.String(http.StatusNotFound, "VKAEP Server allow only GET for root route.")
		return
	}

	dir, file := path.Split(uri)
	ext := filepath.Ext(file)
	//Если запросили корень URL или запрос не файл, то возвращаем index.html
	if (dir == "/" && file == "") || ext == "" {
		dir = "/"
		file = "index.html"
	}

	appPath := path.Join(dir, file)
	app = "manage"

	//Путь к файлам
	filePath := "./frontend/" + path.Join(app, appPath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File "+appPath+" not found.")
		return
	}

	c.File(filePath)
}
