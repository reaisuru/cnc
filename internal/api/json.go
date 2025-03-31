package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, code int, v map[string]interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	http.Error(w, string(data), code)
}
