package muser

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/common"
)

//User structure
type User struct {
	gorm.Model
	UserType string `gorm:"column:user_type;type:enum('manager','staff');" json:"UserType"`
	Email    string `gorm:"column:email;type:varchar(100);unique_index;not null" json:"Email"`
	Password string `gorm:"column:password;type:varchar(100);not null" json:"Password"`
	Token    string `gorm:"column:token;type:varchar(1000);DEFAULT:null"`
}

//TableName return new table name for User struct
func (User) TableName() string {
	return "s_users"
}

//GetByID return information about user by id
func (u *User) GetByID(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}
	var r *gorm.DB

	r = DB.First(&u)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetByToken return information about user by id
func (u *User) GetByToken(c *gin.Context) error {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB

	r = DB.Where("token = ?", u.Token).First(&u)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetByEmail return user by email
func (u *User) GetByEmail(c *gin.Context) error {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB

	r = DB.Where("email = ?", u.Email).First(&u)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetPartnerChildsQuantityByID return user by partner code
func (u *User) GetPartnerChildsQuantityByID(c *gin.Context) (float32, error) {
	var DB = c.MustGet("DB").(*gorm.DB)
	var r *gorm.DB
	var quantity float32

	r = DB.Model(&User{}).Where("partner_id = ?", u.ID).Count(&quantity)
	if r.Error != nil && r.Error.Error() != "record not found" {
		return 0, r.Error
	}

	return quantity, nil
}

//Post create new user in main DB
func (u *User) Post(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	var err error
	var userSysObjName string

	userSysObjName = strings.Split(u.Email, "@")[0] + "_"
	userSysObjName = userSysObjName + strings.Split(strings.Split(u.Email, "@")[1], ".")[0]

	//Hash user password
	u.Password, err = common.HashAndSaltPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	//Create new USER in Main DB
	var r *gorm.DB
	r = DB.Where("email = ?", u.Email).First(&u)

	//Check if user already Exists
	if u.ID != 0 {
		return err
	}

	r = DB.Create(&u)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//PutPassword update user password
func (u *User) PutPassword(c *gin.Context) error {
	var DB = c.MustGet("DB").(*gorm.DB)

	var r *gorm.DB
	var err error
	var password string

	//Hash user password
	password, err = common.HashAndSaltPassword([]byte(u.Password))
	if err != nil {
		return err
	}

	r = DB.Model(&u).Where("id = ?", u.ID).Update("password", password)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
