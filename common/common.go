package common

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

//Config fields
type Config struct {
	Type        string `json:"Type"`
	DBHD        string `json:"DBHD"`
	Port        string `json:"Port"`
	PortService string `json:"PortService"`
	DBRP        string `json:"DBRP"`
	DBRTUP      string `json:"DBRTUP"`
	CertPath    string `json:"CertPath"`
	KeyCertPath string `json:"KeyCertPath"`
}

//SetConfigLocal set config for local development
func (c *Config) SetConfigLocal() {
	c.Type = "local"
	c.DBHD = "127.0.0.1"
	c.Port = "49777"
	c.PortService = "49778"
	c.DBRP = "1"
	c.DBRTUP = "1"
	c.CertPath = ""
	c.KeyCertPath = ""
}

//SetConfigProd set config for prod development
func (c *Config) SetConfigProd() {
	c.Type = "prod"
	c.DBHD = "127.0.0.1"
	c.Port = "49001"
	c.PortService = "49002"
	c.DBRP = "UWW4ghrj#$skjerk32ejlwq"
	c.DBRTUP = "dfgadrtglOu8#$43uuhfdjnJS"
	c.CertPath = "/etc/letsencrypt/live/vkaep.tuvis.world/fullchain.pem"
	c.KeyCertPath = "/etc/letsencrypt/live/vkaep.tuvis.world/privkey.pem"
}

//GetConfigByType return config by type
func GetConfigByType(t string, configs *[]Config) Config {
	config := Config{}

	for _, c := range *configs {
		if c.Type == t {
			config = c
			break
		}
	}

	return config
}

//PrintBinPath print bin path
func PrintBinPath() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	log.Println("Bin = ", exPath)
}

//InitGlobalVars init global variables
func InitGlobalVars() error {
	var config Config

	serverConfigType := flag.String("sct", "local", "")
	releaseFlag := flag.Bool("release", false, "release")
	sslFlag := flag.Bool("ssl", false, "ssl")
	flag.Parse()

	fmt.Println("Server config type: ", *serverConfigType)

	if *serverConfigType == "local" {
		config.SetConfigLocal()
	} else {
		config.SetConfigProd()
	}

	//this command MUST be

	settings.IsRelease = *releaseFlag

	settings.DBHostDefault = config.DBHD
	settings.Port = config.Port
	settings.PortService = config.PortService
	settings.DBRP = config.DBRP
	settings.DBRTUP = config.DBRTUP
	settings.IsSSL = *sslFlag
	settings.CertPath = config.CertPath
	settings.KeyPath = config.KeyCertPath

	log.Println("Port = ", settings.Port)
	log.Println("PortService = ", settings.PortService)

	PrintBinPath()

	initValidator()

	return nil
}

//GetNewUUID returns new UUID
func GetNewUUID() (string, error) {
	uid := uuid.NewV4()

	return uid.String(), nil
}

func random(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn((max - min)) + min
}

//GetNewRandomValue return new random value from start to end
func GetNewRandomValue(start int, finish int) string {
	return strconv.FormatUint(uint64(random(start, finish)), 10)
}

//GetNewCardBarcode return new pasword
func GetNewCardBarcode() string {
	return GetNewRandomValue(100000, 999999)
}

//GetNewClientPassword return new pasword
func GetNewClientPassword() string {
	return GetNewRandomValue(100000, 999999)
}

//GetNewVerificationCode return new verification code
func GetNewVerificationCode() string {
	return GetNewRandomValue(1000, 9999)
}

//GetNewClientCode return new client code
func GetNewClientCode() string {
	return GetNewRandomValue(100000, 999999)
}

//GetNewPartnerCode return new partner code
func GetNewPartnerCode() string {
	return GetNewRandomValue(100000, 999999)
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

//GetCardBarcodeFormatted return barcode for card in format 0000-0000-000000
func GetCardBarcodeFormatted(barcode string) string {
	splittedBarcode := strings.Split(barcode, "-")
	splittedBarcode[0] = leftPad2Len(splittedBarcode[0], "0", 5)
	splittedBarcode[1] = leftPad2Len(splittedBarcode[1], "0", 5)
	return strings.Join(splittedBarcode, "-")
}

//FloatToString convert float to string
func FloatToString(input float32) string {
	return strconv.FormatFloat(float64(input), 'f', 2, 32)
}

//FloatToStringPushMsg convert float to string
func FloatToStringPushMsg(input float32) string {
	return strconv.FormatFloat(float64(input), 'f', 0, 32)
}

/******************************************/
/*CAPTCHA*/
/******************************************/

//CreateCaptchaNew create new captha and return capthca ID
func CreateCaptchaNew(c *gin.Context) {
	d := struct {
		CaptchaID string
	}{
		captcha.New(),
	}
	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "", "data": d})
}

/******************************************/

/******************************************/
/*PASSWORD*/
/******************************************/

//HashAndSaltPassword is hashin incoming password
func HashAndSaltPassword(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

//ComparePasswords compare passwords
func ComparePasswords(hashedPwd string, plainPwd []byte) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

/******************************************/
/*VALIDATOR*/
/******************************************/

//Validate validator
var Validate *validator.Validate

//InitValidator create new global validator
func initValidator() {
	Validate = validator.New()
}

//GetValidationNewSimpleError return a simple error from validation error
func GetValidationNewSimpleError(err error) error {
	var errorString string
	for _, err := range err.(validator.ValidationErrors) {
		errorString += "Field " + err.Namespace() + " has validation Error on '" + err.ActualTag() + "' = '" + err.Param() + "'"
	}
	return errors.New(errorString)
}

//ParseParam parse param
func ParseParam(c *gin.Context, key string) (uint, error) {
	parsedIDuint64, err := strconv.ParseUint(c.Param(key), 0, 64)
	if err != nil {
		return 0, err
	}

	return uint(parsedIDuint64), nil
}
