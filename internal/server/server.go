package server

import (
	"context"
	"github.com/alexandremahdhaoui/ipxe-api/internal/service"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"github.com/labstack/echo/v4"
)

func New() ServerInterface {
	return &server{}
}

type server struct {
	ipxe service.IPXE
}

func (s *server) GetBootIpxe(c echo.Context) error {
	// convert into type
	// call service
	// write response

	//TODO implement me
	panic("implement me")
}

func (s *server) GetConfigByID(c echo.Context, configID UUID, _ GetConfigByIDParams) error {
	// convert into type
	// call service
	// write response

	//TODO implement me
	panic("implement me")
}

func (s *server) GetIpxeByLabels(c echo.Context, _ GetIpxeByLabelsParams) error {
	// convert into type
	selectors, err := types.NewIpxeSelectorsFromContext(c)
	if err != nil {
		return err //TODO: wrap
	}

	// call service
	rendered, err := s.ipxe.FindProfileAndRender(context.Background(), selectors)
	if err != nil {
		return err //TODO: wrap
	}

	// write response
	c.Response().Status = 200
	if _, err := c.Response().Write(rendered); err != nil {
		return err //TODO: wrap
	}

	return nil
}
