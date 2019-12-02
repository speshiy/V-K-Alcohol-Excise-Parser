package msys

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/Tuvis-Server/_client/models/mclient"
)

//City stucture
type City struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(150);unique_index;not null" json:"Name"`
}

//TableName return new table name
func (City) TableName() string {
	return "sys_l_cities"
}

//GetCities return list of cities
func GetCities(c *gin.Context, db *gorm.DB) ([]City, error) {
	var r *gorm.DB
	var sortOrder string
	var cities []City

	ContextClient, _ := c.Get("Client")
	countryID := ContextClient.(mclient.Client).CountryID

	sortOrder = "IFNULL(500), name"

	r = db.Where("country_id = ?", countryID).Order(sortOrder).Find(&cities)
	if r.Error != nil {
		return cities, r.Error
	}

	return cities, nil
}

//PreFill fill table with new data
func (c City) PreFill(db *gorm.DB) error {
	var r *gorm.DB

	for _, sql := range c.GetRefillSQL() {
		r = db.Exec(sql)
		if r.Error != nil && !strings.Contains(r.Error.Error(), "Duplicate") {
			return r.Error
		}
	}

	return nil
}

//GetRefillSQL sql for clean and fill table with new data
func (City) GetRefillSQL() []string {
	return []string{
		//Казахстан
		"INSERT INTO sys_l_cities(id, name) values (1, 'Алматы')",
		"INSERT INTO sys_l_cities(id, name) values (2, 'Нур-Султан')",
		"INSERT INTO sys_l_cities(id, name) values (3, 'Актау')",
		"INSERT INTO sys_l_cities(id, name) values (4, 'Актобе')",
		"INSERT INTO sys_l_cities(id, name) values (5, 'Атырау')",
		"INSERT INTO sys_l_cities(id, name) values (6, 'Жезказган')",
		"INSERT INTO sys_l_cities(id, name) values (7, 'Караганды')",
		"INSERT INTO sys_l_cities(id, name) values (8, 'Кокшетау')",
		"INSERT INTO sys_l_cities(id, name) values (9, 'Кызылорда')",
		"INSERT INTO sys_l_cities(id, name) values (10, 'Павлодар')",
		"INSERT INTO sys_l_cities(id, name) values (11, 'Петропавловск')",
		"INSERT INTO sys_l_cities(id, name) values (12, 'Риддер')",
		"INSERT INTO sys_l_cities(id, name) values (13, 'Семей')",
		"INSERT INTO sys_l_cities(id, name) values (14, 'Талдыкорган')",
		"INSERT INTO sys_l_cities(id, name) values (15, 'Темиртау')",
		"INSERT INTO sys_l_cities(id, name) values (16, 'Туркестан')",
		"INSERT INTO sys_l_cities(id, name) values (17, 'Уральск')",
		"INSERT INTO sys_l_cities(id, name) values (18, 'Усть-Каменогорск')",
		"INSERT INTO sys_l_cities(id, name) values (19, 'Шымкент')",
	}
}
