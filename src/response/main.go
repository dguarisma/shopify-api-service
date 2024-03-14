package response

import "encoding/json"

type Result struct {
	Body   string
	Status int
	Errors []struct {
		Error string
	} `json:"Errors,omitempty"`
}

type ErrMsg struct {
	Error string
}

func NewResult(body string, status int) Result {
	return Result{
		Body:   body,
		Status: status,
	}
}

func NewFailResult(err error, status int) Result {
	e := ErrMsg{err.Error()}
	msgErr, _ := json.Marshal(e)
	return Result{
		Body:   string(msgErr),
		Status: status,
	}
}
