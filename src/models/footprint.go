package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"utility"
)

type Footprint struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	MovieID    int    `json:"movie_id"`
	TimeOnSite string `json:"time_on_site"`
}

func InsertFootprint(footprint Footprint) (res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()

	if footprint.UserID == 0 || footprint.MovieID == 0 || len(footprint.TimeOnSite) == 0 {
		res = REQUIRE_FIELD_EMPTY
		return
	}

	db.Create(&footprint)
	res = SUCCESS
	return
}
