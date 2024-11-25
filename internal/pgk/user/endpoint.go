package user

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/handlers"
)

// endpoint interface
type Endpoint interface {
	Find(c fiber.Ctx) error
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

func NewEndpoint() Endpoint {
	return &endpoint{
		config:  config.CF,
		result:  config.RR,
		service: NewService(),
	}
}

// @Tags User
// @Summary Find User
// @Description Find User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.User
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /.well-known/nostr.json [get]
func (ep *endpoint) Find(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Find, &Request{})
}
