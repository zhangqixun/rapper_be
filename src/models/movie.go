package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"utility"
)

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Year       string `json:"year"`
	Released   string `json:"released"`
	Runtime    string `json:"runtime"`
	Genre      string `json:"genre"`
	Director   string `json:"director"`
	Writer     string `json:"writer"`
	Actors     string `json:"actors"`
	Plot       string `json:"plot"`
	Language   string `json:"language"`
	Country    string `json:"country"`
	Awards     string `json:"awards"`
	Poster     string `json:"poster"`
	ImdbRating string `json:"imdb_rating"`
	ImdbID     string `json:"imdb_id"`
}

func TypeQuery(movieType string) (movies []Movie, count int, res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()

	db.Where("genre like ?", "%"+movieType+"%").Order("imdb_rating desc").Find(&movies).Count(&count)
	res = SUCCESS
	return
}
