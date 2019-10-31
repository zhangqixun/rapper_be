package controllers

import (
	"io/ioutil"
	"net/http"
	"utility"
)

func QuickStart(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)
	body, _ := ioutil.ReadAll(r.Body)
}
