package http

import "moniso-server/internal/config"

type ApiServer struct {
	config *config.Api
}

func New(config *config.Api) *ApiServer {
	return &ApiServer{
		config: config,
	}
}

func (s *ApiServer) Start() error {
	return nil
}
