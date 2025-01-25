package server

import (
	"context"
	"errors"

	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/generated/ipxerserver"
)

var (
	ErrGetConfigByID      = errors.New("getting config by id")
	ErrGetIPXEBySelectors = errors.New("getting ipxe by labels")
)

func New(ipxe controller.IPXE, config controller.Content) ipxerserver.StrictServerInterface {
	return &server{
		ipxe:   ipxe,
		config: config,
	}
}

type server struct {
	ipxe   controller.IPXE
	config controller.Content
}

func (s *server) GetIPXEBootstrap(
	_ context.Context,
	_ ipxerserver.GetIPXEBootstrapRequestObject,
) (ipxerserver.GetIPXEBootstrapResponseObject, error) {
	// call controller
	b := s.ipxe.Boostrap()

	return ipxerserver.GetIPXEBootstrap200TextResponse(b), nil
}

func (s *server) GetContentByID(
	ctx context.Context,
	request ipxerserver.GetContentByIDRequestObject,
) (ipxerserver.GetContentByIDResponseObject, error) {
	// TODO: instantiate child context with correlation ID.

	attributes := types.IPXESelectors{
		Buildarch: string(request.Params.Buildarch),
		UUID:      request.Params.Uuid,
	}

	// call controller
	b, err := s.config.GetByID(ctx, request.ContentID, attributes)
	if err != nil {
		return ipxerserver.GetContentByID500JSONResponse{
			N500JSONResponse: ipxerserver.N500JSONResponse{
				Code:    500,
				Message: errors.Join(err, ErrGetConfigByID).Error(),
			},
		}, nil
	}

	return ipxerserver.GetContentByID200TextResponse(b), nil
}

func (s *server) GetIPXEBySelectors(
	ctx context.Context,
	request ipxerserver.GetIPXEBySelectorsRequestObject,
) (ipxerserver.GetIPXEBySelectorsResponseObject, error) {
	// TODO: create new context with correlation ID.

	// convert into type
	// TODO: use params instead of converting the echo context?
	selectors := types.IPXESelectors{
		Buildarch: string(request.Params.Buildarch),
		UUID:      request.Params.Uuid,
	}

	// call controller
	b, err := s.ipxe.FindProfileAndRender(ctx, selectors)
	if err != nil {
		return ipxerserver.GetIPXEBySelectors500JSONResponse{
			N500JSONResponse: ipxerserver.N500JSONResponse{
				Code:    0,
				Message: errors.Join(err, ErrGetIPXEBySelectors).Error(),
			},
		}, nil
	}

	return ipxerserver.GetIPXEBySelectors200TextResponse(b), nil
}
