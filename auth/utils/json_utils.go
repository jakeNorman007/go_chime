package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
    w.Header().Add("Content-type", "application/json")
    w.WriteHeader(status)

    return json.NewEncoder(w).Encode(v)
}

func ParseJSON(r *http.Request, v any) error {
    if r.Body == nil {
        return fmt.Errorf("Missing request body")
    }

    return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
    WriteJSON(w, status, map[string]string{"Error": err.Error()})
}
