package routing

import "net/http"

const HeaderContentType = "Content-Type"
const ContentTypeApplicationJson = "application/json"

func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
