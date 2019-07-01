package cuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/_core/models/muser"
)

//Get return info about user
func Get(c *gin.Context) {
	var err error
	var user muser.User

	err = GetUser(c, &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "", "data": user})
}

func setUserID(c *gin.Context, u *muser.User) error {
	ContextUser, _ := c.Get("User")
	//Set user id
	u.ID = ContextUser.(muser.User).ID
	return nil
}

//GetUser get user
func GetUser(c *gin.Context, u *muser.User) error {
	var err error

	err = setUserID(c, u)
	if err != nil {
		return err
	}

	err = u.GetByID(c, nil)
	if err != nil {
		return err
	}

	return nil
}
