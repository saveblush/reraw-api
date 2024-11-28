package user

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/saveblush/reraw-api/internal/core/breaker"
	"github.com/saveblush/reraw-api/internal/core/cctx"
	"github.com/saveblush/reraw-api/internal/core/config"
	"github.com/saveblush/reraw-api/internal/core/connection/client"
	"github.com/saveblush/reraw-api/internal/core/generic"
	"github.com/saveblush/reraw-api/internal/core/utils/logger"
	"github.com/saveblush/reraw-api/internal/models"
)

// service interface
type Service interface {
	FindWellKnownName(c *cctx.Context, req *RequestWellKnownName) (interface{}, error)
	FindWellKnownLNURL(c *cctx.Context, req *RequestWellKnownName) (interface{}, error)
}

type service struct {
	config     *config.Configs
	repository Repository
	client     client.Client
}

func NewService() Service {
	return &service{
		config:     config.CF,
		repository: NewRepository(),
		client:     client.New(),
	}
}

// FindWellKnownName find well known name nostr username
func (s *service) FindWellKnownName(c *cctx.Context, req *RequestWellKnownName) (interface{}, error) {
	resNotfound := map[string]interface{}{
		"status": "error",
	}

	if generic.IsEmpty(req.Name) {
		res := resNotfound
		res["message"] = "field validation for 'name'"
		return res, nil
	}

	var relays []string
	if !generic.IsEmpty(s.config.App.LazyRelays) {
		relays = s.config.App.LazyRelays
	}

	fetch := &models.User{}
	err := s.repository.FindByIDString(c.GetRelayDatabase(), "name", req.Name, fetch)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if generic.IsEmpty(fetch.Name) || generic.IsEmpty(fetch.Pubkey) {
		res := resNotfound
		res["message"] = fmt.Sprintf("%s is not found", req.Name)
		return res, nil
	}

	res := map[string]interface{}{
		"names": map[string]interface{}{
			fetch.Name: fetch.Pubkey,
		},
		"relays": map[string]interface{}{
			fetch.Pubkey: relays,
		},
	}

	return res, nil
}

func (s *service) FindWellKnownLNURL(c *cctx.Context, req *RequestWellKnownName) (interface{}, error) {
	resNotfound := map[string]interface{}{
		"status": "error",
	}

	if generic.IsEmpty(req.Name) {
		res := resNotfound
		res["message"] = "field validation for 'name'"
		return res, nil
	}

	fetch := &models.User{}
	err := s.repository.FindByIDString(c.GetRelayDatabase(), "name", req.Name, fetch)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if generic.IsEmpty(fetch.LightningURL) {
		res := resNotfound
		res["message"] = fmt.Sprintf("%s is not found", req.Name)
		return res, nil
	}

	lnDomain := strings.Split(fetch.LightningURL, "@")
	if len(lnDomain) != 2 {
		res := resNotfound
		res["message"] = "invalid LNURL"
		return res, nil
	}

	var res map[string]interface{}
	url := fmt.Sprintf("https://%s/.well-known/lnurlp/%s", lnDomain[1], lnDomain[0])
	_, err = s.client.Get(url, nil, nil, &res, breaker.BreakerName)
	if err != nil {
		logger.Log.Errorf("get lnurl error: %s", err)
		return nil, err
	}

	return res, nil
}
