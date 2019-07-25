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
	app = "vkaep"

	//Путь к файлам
	var fileNotFound bool
	var fileForResponse string

	filePath := []string{"./frontend/" + path.Join(app, appPath), "/var/www/tuvis.world/frontend/" + path.Join(app, appPath)}

	//Перебираем пути где может лежать frontend
	for _, f := range filePath {
		//Обнуляем флаг fileNotFound для текущего пути
		fileNotFound = false
		//Если файл не найден то переходим к следующему пути, если это последний путь то флаг будет в TRUE и пользователь увидит ошибку
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fileNotFound = true
			continue
		} else {
			//Если файл найден то возвращаем первый попавшийся
			fileForResponse = f
			break
		}
	}

	if fileNotFound {
		c.String(http.StatusNotFound, "File "+appPath+" not found.")
		return
	}

	c.File(fileForResponse)
}
