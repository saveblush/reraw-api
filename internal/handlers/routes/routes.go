package routes

import (
	swagger "github.com/saveblush/gofiber3-swagger"

	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/handlers/middlewares"
	"github.com/saveblush/reraw-api/internal/pgk/healthcheck"
	"github.com/saveblush/reraw-api/internal/pgk/system"
	"github.com/saveblush/reraw-api/internal/pgk/user"
)

// InitRouter init router
func (s *server) InitRouter() {
	// api
	api := s.Group(config.CF.App.ApiBaseUrl)

	// system
	systemEndpoint := system.NewEndpoint()
	systemRoute := api.Group("/system")
	systemRoute.Post("/action", systemEndpoint.Action, middlewares.AuthorizationAdminRequired())
	systemRoute.Get("/maintenance", middlewares.Maintenance())

	api.Use(
		middlewares.Available(), // ปิด/เปิด ระบบ
		middlewares.AcceptLanguage(),
	)

	// healthcheck endpoint
	healthCheckEndpoint := healthcheck.NewEndpoint()
	api.Get("/healthcheck", healthCheckEndpoint.HealthCheck)

	// api v1
	v1 := api.Group("/v1")

	// swagger
	if config.CF.Swagger.Enable {
		v1.Get("/swagger/*", swagger.HandlerDefault)
	}

	v1.Get("/healthcheck", healthCheckEndpoint.HealthCheck)

	// user nostr
	userEndpoint := user.NewEndpoint()
	userRoute := s
	userRoute.Get(".well-known/nostr.json", userEndpoint.FindWellKnownName)
	userRoute.Get(".well-known/lnurlp/:name", userEndpoint.FindWellKnownLNURL)

	// not found
	s.Use(middlewares.Notfound())
}
