package router

import (
	"moniso-server/internal/delivery/http_api/v1/handler"
	"net/http"
)

func AuthRouter(handler *handler.AuthHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/registration", handler.Registration)
	router.HandleFunc("/authorization", handler.Authorization)

	return router
}
