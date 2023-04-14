package router

import (
	"net/http"
)

func AuthRouter(registrationHandler, authorizationHandler, refreshHandler http.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/registration", registrationHandler)
	router.Handle("/authorization", authorizationHandler)
	router.Handle("/refresh", refreshHandler)

	return router
}
