package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type ResponseFormatter struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, p interface{}, status int) {
	encodedData, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(encodedData)
}

func APIResponse(message string, code int, status string, data interface{}) ResponseFormatter {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	responseJSON := ResponseFormatter{
		Meta: meta,
		Data: data,
	}

	return responseJSON
}
