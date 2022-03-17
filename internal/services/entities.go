package services

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	writeResponse(w http.ResponseWriter) error
}

type errResponse struct {
	Error string `json:"error"`
}

func (e *errResponse) writeResponse(w http.ResponseWriter) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

/*func (er errResponse) MarshalJSON() ([]byte, error) {

}*/
