package controller

import (
	"encoding/json"
	"fmt"
)

type Delete struct {
	ID uint
}

func HandleDelete(body []byte) (string, error) {
	adapt := Delete{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		return "", ErrFormatForID
	}
	return fmt.Sprint(adapt.ID), nil
}
