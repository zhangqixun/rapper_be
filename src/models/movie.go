package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
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

type Rec struct {
	User int    `json:"user"`
	Recs string `json:"recs"`
}

func IDQuery(id int) (movie Movie, res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()

	db.Where("id = ?", id).First(&movie)
	if movie.ID == 0 {
		res = NO_DATA
	} else {
		res = SUCCESS
	}
	return
}

func KeywordQuery(keyword string, kind string) (movies []Movie, count int, res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()

	if kind == "title" {
		db.Where("title like ?", "%"+keyword+"%").Order("imdb_rating desc").Find(&movies).Count(&count)
	} else if kind == "actor" {
		db.Where("actors like ?", "%"+keyword+"%").Order("imdb_rating desc").Find(&movies).Count(&count)
	}
	res = SUCCESS
	return
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

func SimilarityQuery(movie string) (movies []Movie, res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()
	moviestr, err := utility.QueryTopSimilarities(movie)
	if err != nil {
		res = DB_ERROR
		return
	}
	res = SUCCESS
	for _, v := range moviestr {
		var movie Movie
		db.Where("id = ?", v).First(&movie)
		movies = append(movies, movie)
	}
	return
}

func GetRecommendationViaFootprint(user string) (movies []Movie, res int) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()
	var recs Rec
	db.Where("user = ?", user).First(&recs)
	ids := strings.Split(recs.Recs, " ")
	res = SUCCESS
	for _, v := range ids {
		var movie Movie
		db.Where("id = ?", v).First(&movie)
		movies = append(movies, movie)
	}
	return
}
