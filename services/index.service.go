package services

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

    w.Write([]byte("<h1>Hello World!</h1>"))
}