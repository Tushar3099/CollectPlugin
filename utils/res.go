package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

type response struct {
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data"`
}

func WriteError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = errorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, err := json.Marshal(response)

	if err != nil {
		panic("Error marshalling ErrorResponse")
	}

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func WriteResponse(res interface{}, w http.ResponseWriter) {

	// log.Fatal(err.Error())
	var response = response{
		StatusCode: http.StatusOK,
		Data:       res,
	}

	message, err := json.Marshal(response)

	if err != nil {
		panic("Error marshalling Response")
	}

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
