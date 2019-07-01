package cuser

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/models/muser"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/common"
)

//RegistrationInfo struct
type RegistrationInfo struct {
	Password        string `json:"Password"`
	PasswordConfirm string `json:"PasswordConfirm"`
	Email           string `json:"Email"`
}

//GetToken return token for session
func GetToken(c *gin.Context) (string, error) {
	h := c.GetHeader("X-Token")
	if h == "" {
		return "", errors.New("E_USER_TOKEN_CONFIRM_NOT_FOUND")
	}
	return h, nil
}

//SignUp new user
func SignUp(c *gin.Context) {
	var err error
	var registrationInfo RegistrationInfo
	var user muser.User

	if err = c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = common.Validate.Var(registrationInfo.Email, "required,email")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "E_EMAIL_INVALID"})
		return
	}
	err = common.Validate.Var(registrationInfo.Password, "required,min=6")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "E_PASSWORD_INVALID"})
		return
	}
	err = common.Validate.Var(registrationInfo.PasswordConfirm, "required,min=6")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "E_PASSWORD_INVALID"})
		return
	}
	if registrationInfo.Password != registrationInfo.PasswordConfirm {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "E_PASSWORDS_DOES_NOT_MATCH"})
		return
	}

	user.Password = registrationInfo.Password
	user.Email = registrationInfo.Email

	err = CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//OK response to client
	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "SIGN_UP_SUCCESS"})
}

//CreateUser operations for creating user
func CreateUser(c *gin.Context, user *muser.User) error {
	var err error
	err = user.Post(c, nil)
	if err != nil {
		return err
	}

	return nil
}

//Login user
func Login(c *gin.Context) {
	var err error
	var user muser.User
	var userResponse muser.User

	//Auth user by token confirm
	tokenConfirm := c.DefaultQuery("token_confirm", "")

	//Auth user by token restore
	tokenRestorePassword := c.DefaultQuery("token_restore_password", "")

	//determine which token come. Declare slice which contain income token name in base and token value
	tokenIncome := make([]string, 2, 2)

	if tokenConfirm != "" {
		tokenIncome[0] = tokenConfirm
		tokenIncome[1] = "token_confirm"
	}

	if tokenRestorePassword != "" {
		tokenIncome[0] = tokenRestorePassword
		tokenIncome[1] = "token_reset"
	}

	//if TokenIncome not empty and method GET than auth by this token
	if (tokenIncome[0] != "") && (c.Request.Method == "GET") {
		//set to user model income Token
		user.Token = tokenIncome[0]

		//Auth user by confirm or restore token which takes from email
		user, err = AuthToken(c, user, tokenIncome[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"showError": "false", "status": "false", "code": "token_error", "message": "E_USER_TOKEN_CONFIRM_NOT_FOUND"})
			return
		}

		//Creates user TOKEN and put it in Session
		err = CreateUserToken(c, &user, "token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
			return
		}

		userResponse.ID = user.ID
		err = userResponse.GetByID(c, nil)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "true", "message": "", "data": userResponse})
		return
	}

	user.Token, err = GetToken(c)

	//Bind requst to user model
	if err = c.ShouldBindJSON(&user); err != nil {
		if user.Token == "" {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
	}

	//Try to auth by token from LocalStorage if it exists and email is empty
	if (user.Token != "") && (user.Email == "") {
		user, err = AuthToken(c, user, "token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"showError": "false", "status": "false", "code": "token_error", "message": err.Error()})
			return
		}
		//Set new token in header and response OK
		c.Set("token", user.Token)

		userResponse.ID = user.ID
		err = userResponse.GetByID(c, nil)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "true", "message": "", "data": userResponse})
		return
	}

	//If requset GET, that means it's a try to auth by token on start app. If all empty than return
	if (c.Request.Method == "GET") && (user.Token == "") && (user.Email == "") {
		c.JSON(http.StatusUnauthorized, gin.H{"showError": "false", "status": "false", "message": "E_USER_TOKEN_NOT_FOUND"})
		return
	}

	err = common.Validate.Var(user.Email, "required,email")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "E_EMAIL_INVALID"})
		return
	}
	err = common.Validate.Var(user.Password, "required,min=6")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "E_PASSWORD_INVALID"})
		return
	}

	//Auth user by email and password
	err = user.LoginByEmailPassword(c, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
		return
	}

	//Creates user TOKEN and put it in Session
	err = CreateUserToken(c, &user, "token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
		return
	}

	userResponse.ID = user.ID
	err = userResponse.GetByID(c, nil)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "", "data": userResponse})
}

//AuthToken trying to autorize with token, has tokenType
func AuthToken(c *gin.Context, user muser.User, tokenType string) (muser.User, error) {
	var err error
	var tokenClaims jwt.StandardClaims

	if tokenType == "token" {
		tokenClaims, err = common.JWTParse(user.Token)
		if err != nil {
			return user, err
		}

		if time.Now().Unix() >= tokenClaims.ExpiresAt {
			return user, errors.New("E_TOKEN_EXPIRED")
		}

		// userID, err := strconv.Atoi(tokenClaims.Id)
		// user.ID = uint(userID)
		err = user.GetByToken(c)

		//Set new token into header
		setTokenInHeader(c, user.Token)
	} else {
		//Auth user from DB
		err = user.LoginByToken(c, tokenType)
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

//CreateUserToken create user token
func CreateUserToken(c *gin.Context, user *muser.User, tokenType string) error {
	var err error
	var tempToken string

	//Create new token for user
	tempToken, err = user.CreateToken(c, tokenType, nil)
	if err != nil {
		return err
	}

	if tokenType == "token" {
		user.Token = tempToken
		//Set new token into header
		setTokenInHeader(c, user.Token)
	}

	return nil
}

func setTokenInHeader(c *gin.Context, token string) {
	c.Header("X-Token", token)
	c.Request.Header.Set("X-Token", token)
}

func setTokenDemoInHeader(c *gin.Context, token string) {
	c.Header("X-Demo-Token", token)
	c.Request.Header.Set("X-Demo-Token", token)
}
