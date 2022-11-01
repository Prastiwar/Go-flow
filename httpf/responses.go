package httpf

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, status int, data interface{}) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write(v)
	return err
}
