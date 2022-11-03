package httpf

import (
	"encoding/json"
	"net/http"
)

// Json marshals the data and writes it to http.ResponseWriter with given status code.
func Json(w http.ResponseWriter, status int, data interface{}) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}

	v, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write(v)
	return err
}
