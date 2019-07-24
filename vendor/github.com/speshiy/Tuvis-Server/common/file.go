package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/Tuvis-Server/settings"
)

//SaveImageBase64 in the OS and return new generated name. Save depents of entityType. entityType set the path on to OS
func SaveImageBase64(c *gin.Context, userID uint, data string, entityType string) (string, error) {
	var (
		err             error
		filePath        string
		newFileName     string
		newFileURLPath  string
		ErrInvalidImage = errors.New("File is not an image")
	)

	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", ErrInvalidImage
	}

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return "", err
	}

	//Check file size, if size > 2MB than return
	if (binary.Size(unbased) / 1024 / 1024) > 3 {
		return "", errors.New(GetError(c, "E_FILE_SIZE"))
	}

	r := bytes.NewReader(unbased)

	imageType := data[11:idx]

	//Set new file name
	newFileName, _ = GetNewUUID()
	newFileName = newFileName + "." + imageType
	newFileURLPath = newFileName

	//Determine entity-type and set path for entity image
	switch entityType {
	case "user_company":
		filePath = filepath.Join(settings.ResourcesPath, "company_logo", "user", strconv.Itoa(int(userID)))
	case "card_logo":
		filePath = filepath.Join(settings.ResourcesPath, "card_logo", "user", strconv.Itoa(int(userID)))
	case "item_logo":
		filePath = filepath.Join(settings.ResourcesPath, "item_logo", "user", strconv.Itoa(int(userID)))
	}

	//Create directory
	err = os.MkdirAll("/"+filePath, 0777)
	if err != nil {
		return "", err
	}

	newFileName = filepath.Join(filePath, newFileName)
	log.Println("New file path = ", newFileName)

	switch imageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			return "", err
		}

		err = checkImageDimension(c, im)
		if err != nil {
			return "", err
		}

		f, err := os.OpenFile("/"+newFileName, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", err
		}

		err = png.Encode(f, im)
		if err != nil {
			return "", err
		}

	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			return "", err
		}

		err = checkImageDimension(c, im)
		if err != nil {
			return "", err
		}

		f, err := os.OpenFile("/"+newFileName, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", err
		}

		err = jpeg.Encode(f, im, nil)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New(GetError(c, "E_FILE_TYPE"))
	}

	switch entityType {
	case "user_company":
		newFileURLPath = filepath.Join("resources", "company_logo", "user", strconv.Itoa(int(userID)), newFileURLPath)
	case "card_logo":
		newFileURLPath = filepath.Join("resources", "card_logo", "user", strconv.Itoa(int(userID)), newFileURLPath)
	case "item_logo":
		newFileURLPath = filepath.Join("resources", "item_logo", "user", strconv.Itoa(int(userID)), newFileURLPath)
	}

	// log.Println("New file url path = ", newFileURLPath)
	return newFileURLPath, nil
}

//UploadImageBase64 upload image to firebas storage
func UploadImageBase64(c *gin.Context, id uint, data string, entityType string) (string, error) {
	var (
		err             error
		filePath        string
		newFileName     string
		ErrInvalidImage = errors.New("File is not an image")
	)

	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", ErrInvalidImage
	}

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return "", err
	}

	//Check file size, if size > 2MB than return
	if (binary.Size(unbased) / 1024 / 1024) > 3 {
		return "", errors.New(GetError(c, "E_FILE_SIZE"))
	}

	r := bytes.NewReader(unbased)
	imageType := data[11:idx]

	switch imageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			return "", err
		}

		err = checkImageDimension(c, im)
		if err != nil {
			return "", err
		}
	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			return "", err
		}

		err = checkImageDimension(c, im)
		if err != nil {
			return "", err
		}
	}

	//Set new file name
	newFileName, _ = GetNewUUID()
	newFileName = newFileName + "." + imageType

	//Determine entity-type and set path for entity image
	switch entityType {
	case "user_company":
		filePath = filepath.Join("company_logo", "user", strconv.Itoa(int(id)))
	case "card_logo":
		filePath = filepath.Join("card_logo", "user", strconv.Itoa(int(id)))
	case "item_logo":
		filePath = filepath.Join("item_logo", "user", strconv.Itoa(int(id)))
	case "client_plastic_card":
		filePath = filepath.Join("card_plastic_logo", "client", strconv.Itoa(int(id)))
	}

	newFileName = filepath.Join(filePath, newFileName)

	//Set unix path to name
	obj := GetUnixFileName(newFileName)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[idx+8:]))

	//Upload image to Google Storage
	url, err := UploadFileToFirebaseStorageBucket(obj, &reader)
	if err != nil {
		return "", err
	}

	return url, nil
}

//GetUnixFileName change all back-slashes on slashes
func GetUnixFileName(filename string) string {
	return filepath.ToSlash(filename)
}

//DeleteFile delete file from OS
func DeleteFile(path string) error {
	var err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func getImageWH(img image.Image) (int, int) {
	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y

	return width, height
}

func checkImageDimension(c *gin.Context, img image.Image) error {
	width, height := getImageWH(img)
	if width < 100 || height < 100 {
		return errors.New(GetError(c, "E_IMG_DIMENSION") + " 300x200")
	}
	return nil
}
