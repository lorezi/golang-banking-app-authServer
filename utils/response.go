package utils

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}

func NotAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func AuthorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
