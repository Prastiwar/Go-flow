package httpf

import (
	"encoding/json"
	"net/http"
)

// Json marshals the data and writes it to http.ResponseWriter with given status code.
// "Content-Type" header is set to "application/json".
func Json(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Add(ContentTypeHeader, ApplicationJsonType)

	v, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write(v)
	return err
}
