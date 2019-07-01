package muser

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
)

//CreateToken create or update token in USER struct. Token Type: token, token_confirm, token_reset
func (u *User) CreateToken(c *gin.Context, tokenType string, DB *gorm.DB) (string, error) {
	var r *gorm.DB

	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	var err error
	var newToken string
	if tokenType == "token" {
		newToken, err = common.JWTCreate(u.ID, "global", 720)
	} else {
		newToken, err = common.GetNewUUID()
	}
	if err != nil {
		return "", err
	}
	r = DB.Model(&u).Where("id = ?", u.ID).Update(tokenType, newToken)
	if r.Error != nil {
		return "", r.Error
	}

	return newToken, nil
}

//LoginByEmailPassword auth user by username and password
func (u *User) LoginByEmailPassword(c *gin.Context, password string) error {
	var DB = c.MustGet("DB").(*gorm.DB)

	var r *gorm.DB
	r = DB.Where("email = ?", u.Email).First(&u)
	if r.Error != nil {
		if r.Error.Error() == "record not found" {
			return errors.New("E_USER_WRONG_EMAIL")
		}
		return r.Error
	}

	b, err := common.ComparePasswords(u.Password, []byte(password))
	if err != nil {
		return err
	}
	if !b {
		return errors.New("E_USER_WRONG_PASSWORD")
	}

	return nil
}

//LoginByToken is authorizing user by token
func (u *User) LoginByToken(c *gin.Context, tokenType string) error {
	var DB = c.MustGet("DB").(*gorm.DB)

	var r *gorm.DB
	r = DB.Where(tokenType+" = ?", u.Token).First(&u)
	if r.Error != nil {
		return r.Error
	}

	//Clear token confirm after success use
	if tokenType == "token_confirm" {
		r = DB.Model(&u).Where("id = ?", u.ID).Update("token_confirm", gorm.Expr("NULL"))
		if r.Error != nil {
			return r.Error
		}
	}

	//Clear token reset after success use
	if tokenType == "token_reset" {
		r = DB.Model(&u).Where("id = ?", u.ID).Update("token_reset", gorm.Expr("NULL"))
		if r.Error != nil {
			return r.Error
		}
	}

	return nil
}

//Logout user
func (u *User) Logout(c *gin.Context) error {
	var DB = c.MustGet("DB").(*gorm.DB)

	var r *gorm.DB
	r = DB.Model(&u).Where("id = ?", u.ID).Update("token", "null")
	if r.Error != nil {
		return r.Error
	}
	return nil
}
