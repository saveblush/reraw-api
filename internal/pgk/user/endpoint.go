package user

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/handlers"
)

// endpoint interface
type Endpoint interface {
	FindWellKnownName(c fiber.Ctx) error
	FindWellKnownLNURL(c fiber.Ctx) error
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
// @Summary FindWellKnownName
// @Description FindWellKnownName
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
func (ep *endpoint) FindWellKnownName(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindWellKnownName, &RequestWellKnownName{})
}

// @Tags User
// @Summary FindWellKnownLNURL
// @Description FindWellKnownLNURL
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.User
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /.well-known/lnurlp/{id} [get]
func (ep *endpoint) FindWellKnownLNURL(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindWellKnownLNURL, &RequestWellKnownName{})
}
