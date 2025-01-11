package healthcheck

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/handlers/render"
	"github.com/saveblush/reraw-api/internal/models"
)

// endpoint interface
type Endpoint interface {
	HealthCheck(c fiber.Ctx) error
}

type endpoint struct {
	config *config.Configs
	result *config.ReturnResult
}

func NewEndpoint() Endpoint {
	return &endpoint{
		config: config.CF,
		result: config.RR,
	}
}

// HealthCheck health check
func (ep *endpoint) HealthCheck(c fiber.Ctx) error {
	return render.JSON(c, models.NewSuccessMessage())
}
