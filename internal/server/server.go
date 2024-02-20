package server

import (
	"github.com/alexandremahdhaoui/ipxe-api/internal/dsa"
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

func (s *server) GetConfigByID(c echo.Context, configID UUID, _ GetConfigByIDParams) error {

}

func (s *server) GetIpxeByLabels(c echo.Context, _ GetIpxeByLabelsParams) error {
	selectors, err := dsa.NewIpxeParamsFromContext(c)
	if err != nil {
		return err // TODO: wrap this err
	}

	return nil
}
