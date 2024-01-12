package models

import (
	"database/sql/driver"
	"gorm.io/gorm"
)

var ItemIndex = "item"

type Color string

const (
	Red    Color = "Red"
	Blue   Color = "Blue"
	Green  Color = "Green"
	Yellow Color = "Yellow"
	Purple Color = "Purple"
	Orange Color = "Orange"
	Pink   Color = "Pink"
)

func (c *Color) Scan(v interface{}) error {
	*c = Color(v.([]byte))
	return nil
}

func (c Color) Value() (driver.Value, error) {
	return string(c), nil
}

type Item struct {
	gorm.Model
	Title string `gorm:"unique"`
	//Color       Color  `gorm:"type:enum('Red', 'Blue', 'Green', 'Yellow', 'Purple', 'Orange', 'Pink') default:'Red'"`
	Price       uint `gorm:"default:0"`
	Description string
}
