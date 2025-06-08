package routing

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	// USTA "proxy" API routes
	router.HandleFunc("GET /usta/organization/{id}/teams", GetUSTAOrganizationTeams)
	router.HandleFunc("GET /usta/organization/{id}/matches", GetUSTAOrganizationMatches)

	// App API routes

	return router
}
