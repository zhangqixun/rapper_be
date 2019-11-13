package controllers

import (
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"strconv"
	"utility"
)

type MovieIDQueryReq struct {
	ID int `json:"id"`
}

type MovieIDQueryRes struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    models.Movie `json:"data"`
}

type MovieKeywordQueryReq struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Kind    string `json:"kind"`
}

type MovieTypeQueryReq struct {
	Genre   string `json:"genre"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type MovieQueryData struct {
	Page     int                           `json:"page"`
	PerPage  int                           `json:"per_page"`
	Projects [models.PER_PAGE]models.Movie `json:"projects"`
	Total    int                           `json:"total"`
}

type MovieQueryRes struct {
	Code    string         `json:"code"`
	Data    MovieQueryData `json:"data"`
	Message string         `json:"message"`
}

type MovieSimilarityRes struct {
	Code    string   `json:"code"`
	Data    []string `json:"data"`
	Message string   `json:"message"`
}

func MovieIDQuery(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	vars := r.URL.Query()

	var req MovieIDQueryReq
	id, ok := vars["id"]

	var info MovieIDQueryRes

	if !ok {
		info.Code = models.REQUIRE_FIELD_EMPTY_CODE
		info.Message = models.REQUIRE_FIELD_EMPTY_MESS
	} else {
		var convErr error
		req.ID, convErr = strconv.Atoi(id[0])
		if convErr != nil {
			info.Code = models.TYPE_ERROR_CODE
			info.Message = models.TYPE_ERROR_MESS
		} else {
			var res int
			info.Data, res = models.IDQuery(req.ID)
			if res == models.DB_ERROR {
				info.Code = models.DB_ERROR_CODE
				info.Message = models.DB_ERROR_MESS
			} else if res == models.NO_DATA {
				info.Code = models.NO_DATA_CODE
				info.Message = models.NO_DATA_MESS
			} else if res == models.SUCCESS {
				info.Code = models.SUCCESS_CODE
				info.Message = models.SUCCESS_MESS
			}
		}
	}

	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func MovieKeywordQuery(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	vars := r.URL.Query()

	var req MovieKeywordQueryReq
	keyword, ok1 := vars["keyword"]
	page, ok2 := vars["page"]
	perPage, ok3 := vars["per_page"]
	kind, ok4 := vars["kind"]

	var info MovieQueryRes

	if !ok1 || !ok2 || !ok3 || !ok4 {
		info.Code = models.REQUIRE_FIELD_EMPTY_CODE
		info.Message = models.REQUIRE_FIELD_EMPTY_MESS
	} else {
		req.Keyword = keyword[0]
		req.Kind = kind[0]
		var convErr1 error
		req.Page, convErr1 = strconv.Atoi(page[0])
		var convErr2 error
		req.PerPage, convErr2 = strconv.Atoi(perPage[0])
		if convErr1 != nil || convErr2 != nil {
			info.Code = models.TYPE_ERROR_CODE
			info.Message = models.TYPE_ERROR_MESS
		} else {
			info.Data.Page = req.Page
			info.Data.PerPage = req.PerPage
			var movies []models.Movie
			var res int
			movies, info.Data.Total, res = models.KeywordQuery(req.Keyword, req.Kind)
			if res == models.DB_ERROR {
				info.Code = models.DB_ERROR_CODE
				info.Message = models.DB_ERROR_MESS
			} else if res == models.SUCCESS {
				if info.Data.Total == 0 {
					info.Code = models.NO_DATA_CODE
					info.Message = models.NO_DATA_MESS
				} else if (info.Data.Page-1)*info.Data.PerPage >= info.Data.Total { //查询越界
					info.Code = models.QUERY_CROSS_BORDER_CODE
					info.Message = models.QUERY_CROSS_BORDER_MESS
				} else {
					for i := 0; i < info.Data.PerPage; i++ {
						if (info.Data.Page-1)*info.Data.PerPage+i < info.Data.Total {
							info.Data.Projects[i] = movies[(info.Data.Page-1)*info.Data.PerPage+i]
						} else {
							break
						}
					}
					info.Code = models.SUCCESS_CODE
					info.Message = models.SUCCESS_MESS
				}
			}
		}
	}

	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func MovieTypeQuery(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	vars := r.URL.Query()

	var req MovieTypeQueryReq
	genre, ok1 := vars["genre"]
	page, ok2 := vars["page"]
	perPage, ok3 := vars["per_page"]

	var info MovieQueryRes

	if !ok1 || !ok2 || !ok3 {
		info.Code = models.REQUIRE_FIELD_EMPTY_CODE
		info.Message = models.REQUIRE_FIELD_EMPTY_MESS
	} else {
		req.Genre = genre[0]
		var convErr1 error
		req.Page, convErr1 = strconv.Atoi(page[0])
		var convErr2 error
		req.PerPage, convErr2 = strconv.Atoi(perPage[0])
		if convErr1 != nil || convErr2 != nil {
			info.Code = models.TYPE_ERROR_CODE
			info.Message = models.TYPE_ERROR_MESS
		} else {
			info.Data.Page = req.Page
			info.Data.PerPage = req.PerPage
			var movies []models.Movie
			var res int
			movies, info.Data.Total, res = models.TypeQuery(req.Genre)
			if res == models.DB_ERROR {
				info.Code = models.DB_ERROR_CODE
				info.Message = models.DB_ERROR_MESS
			} else if res == models.SUCCESS {
				if info.Data.Total == 0 {
					info.Code = models.NO_DATA_CODE
					info.Message = models.NO_DATA_MESS
				} else if (info.Data.Page-1)*info.Data.PerPage >= info.Data.Total { //查询越界
					info.Code = models.QUERY_CROSS_BORDER_CODE
					info.Message = models.QUERY_CROSS_BORDER_MESS
				} else {
					for i := 0; i < info.Data.PerPage; i++ {
						if (info.Data.Page-1)*info.Data.PerPage+i < info.Data.Total {
							info.Data.Projects[i] = movies[(info.Data.Page-1)*info.Data.PerPage+i]
						} else {
							break
						}
					}
					info.Code = models.SUCCESS_CODE
					info.Message = models.SUCCESS_MESS
				}
			}
		}
	}

	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func MovieSimilarityQuery(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)
	vars := r.URL.Query()
	movie, exists := vars["movie"]
	info := MovieSimilarityRes{}
	if !exists {
		info.Code = models.REQUIRE_FIELD_EMPTY_CODE
		info.Message = models.REQUIRE_FIELD_EMPTY_MESS
	} else {
		movies, code := models.SimilarityQuery(movie[0])
		if code != models.SUCCESS {
			info.Code = models.DB_ERROR_CODE
			info.Message = models.DB_ERROR_MESS
		} else {
			info.Code = models.SUCCESS_CODE
			info.Message = models.SUCCESS_MESS
			info.Data = movies
		}
	}
	ret, _ := json.Marshal(info)
	fmt.Fprint(w, string(ret))
}
