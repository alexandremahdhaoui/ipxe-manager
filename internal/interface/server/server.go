package server

import (
	"github.com/labstack/echo/v4"
)

func New() ServerInterface {
	return &server{}
}

type server struct{}

func (s *server) GetBootIpxe(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetConfigByID(c echo.Context, configID UUID, params GetConfigByIDParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetIpxeByLabels(c echo.Context, params GetIpxeByLabelsParams) error {
	//TODO implement me
	panic("implement me")
}
