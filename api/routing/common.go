package routing

import (
	"encoding/json"
	"net/http"
)

const ustaOrganizationURL = "https://leagues.ustanorcal.com/organization.asp?id=%s"

const HeaderContentType = "Content-Type"
const ContentTypeApplicationJson = "application/json"

func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)

	var errorResponse = struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	data, _ := json.Marshal(errorResponse)
	w.Write(data)
	http.Error(w, string(data), statusCode)
}
