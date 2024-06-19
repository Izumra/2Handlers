package server

import (
	_ "github.com/Izumra/2Handlers/docs"
	"github.com/Izumra/2Handlers/internal/server/handlers/auth"
	"github.com/Izumra/2Handlers/internal/server/handlers/news"
	serviceAuth "github.com/Izumra/2Handlers/internal/service/auth"
	serviceNews "github.com/Izumra/2Handlers/internal/service/news"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type app struct {
	srv         *fiber.App
	serviceNews *serviceNews.Service
	serviceAuth *serviceAuth.Service
}

func New(
	serviceNews *serviceNews.Service,
	serviceAuth *serviceAuth.Service,
) *app {

	srv := fiber.New()

	return &app{
		srv,
		serviceNews,
		serviceAuth,
	}
}

func (s *app) RegHandlers() {
	s.srv.Get("/swagger/*", swagger.HandlerDefault)
	auth.MountAuthHandlers(s.srv, s.serviceAuth)
	news.MountNewsHandlers(s.srv, s.serviceNews)
}

func (s *app) Start(addr string, servChan chan<- error) {
	servChan <- s.srv.Listen(addr)
}

func (s *app) Stop() error {
	return s.srv.Shutdown()
}
