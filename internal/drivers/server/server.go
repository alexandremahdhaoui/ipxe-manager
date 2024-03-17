package server

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/controllers"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/labstack/echo/v4"
)

func New(ipxe controllers.IPXE, config controllers.Config) ServerInterface {
	return &server{
		ipxe:   ipxe,
		config: config,
	}
}

type server struct {
	ipxe   controllers.IPXE
	config controllers.Config
}

func (s *server) GetBootIpxe(c echo.Context) error {
	// call controllers
	b := s.ipxe.Boostrap()

	// write response
	if _, err := c.Response().Write(b); err != nil {
		return err //TODO: wrap me
	}
	c.Response().Status = 200

	return nil
}

func (s *server) GetConfigByID(c echo.Context, profileName string, configID UUID, _ GetConfigByIDParams) error {
	// call controllers
	b, err := s.config.GetByID(context.Background(), profileName, configID)
	if err != nil {
		return err //TODO: wrap me
	}

	// write response
	if _, err := c.Response().Write(b); err != nil {
		return err //TODO: wrap me
	}
	c.Response().Status = 200

	return nil
}

func (s *server) GetIpxeByLabels(c echo.Context, _ GetIpxeByLabelsParams) error {
	// convert into type
	selectors, err := types.NewIpxeSelectorsFromContext(c)
	if err != nil {
		return err //TODO: wrap
	}

	// call controllers
	b, err := s.ipxe.FindProfileAndRender(context.Background(), selectors)
	if err != nil {
		return err //TODO: wrap
	}

	// write response
	if _, err := c.Response().Write(b); err != nil {
		return err //TODO: wrap
	}
	c.Response().Status = 200

	return nil
}
