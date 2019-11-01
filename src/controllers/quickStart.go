package controllers

import (
	"net/http"
)

func QuickStart(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("QuickStart"))
}
