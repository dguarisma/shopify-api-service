package controller

import (
	"desarrollosmoyan/lambda/src/response"
	"encoding/json"
	"net/http"
)

func FormateBody(data interface{}, successStatus int) response.Result {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return response.NewFailResult(err, http.StatusInternalServerError)
	}
	return response.NewResult(string(bytes), successStatus)
}
