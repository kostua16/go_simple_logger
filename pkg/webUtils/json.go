package webUtils

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data interface{}) error {
	response, fmtErr := json.Marshal(data)
	if fmtErr != nil {
		return fmtErr
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	_, err := w.Write(response)
	if err != nil {
		return err
	}
	return nil
}

func ReadJson(request *http.Request, data interface{}) error {
	var result error = nil
	defer func(Body io.ReadCloser) {
		if result == nil {
			err := Body.Close()
			if err != nil {
				result = err
			}
		}
	}(request.Body)
	result = json.NewDecoder(request.Body).Decode(data)
	return result
}
