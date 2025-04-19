package routing

import (
	"fmt"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in hello handler")
		w.Write([]byte("hello"))
	})

	router.HandleFunc("GET /organization/{id}/matches", GetOrganizationMatches)

	return router
}
