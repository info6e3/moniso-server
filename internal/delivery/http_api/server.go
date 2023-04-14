package http_api

import (
	"log"
	"moniso-server/internal/config"
	"moniso-server/internal/delivery/http_api/v1/handler"
	"moniso-server/internal/delivery/http_api/v1/middleware"
	"moniso-server/internal/delivery/http_api/v1/router"
	"moniso-server/internal/services"
	"net/http"
)

type ApiServer struct {
	config *config.Api
}

func NewApiServer(config *config.Api) *ApiServer {
	return &ApiServer{
		config: config,
	}
}

func (s *ApiServer) Start(services *services.Services) error {

	apiRouter := http.NewServeMux()

	authRouter := router.AuthRouter(
		handler.Registration{
			Service: services.Auth,
		},
		handler.Authorization{
			Service: services.Auth,
		},
		handler.Refresh{
			Service: services.Auth,
		})
	apiRouter.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	apiRouter.Handle("/point-types", &middleware.Auth{
		Service: services.Auth,
		Next: handler.PointType{
			Service: services.PointType,
		},
	})

	apiRouter.Handle("/points", &middleware.Auth{
		Service: services.Auth,
		Next: handler.Point{
			ServicePoint:     services.Point,
			ServicePointType: services.PointType,
		},
	})

	apiRouter.Handle("/events", &middleware.Auth{
		Service: services.Auth,
		Next: handler.Event{
			Service: services.Event,
		},
	})

	err := http.ListenAndServe(":80", apiRouter)
	if err != nil {
		return err
	}
	log.Println("Сервер запущен на http://localhost:80")
	return nil
}
