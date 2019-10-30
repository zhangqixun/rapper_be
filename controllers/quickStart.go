package controllers

import (
	"io/ioutil"
	"net/http"
)

func QuickStart(w http.ResponseWriter, r *http.Request)  {
	PreprocessXHR(&w,r)
	body, _ := ioutil.ReadAll(r.Body)
}
