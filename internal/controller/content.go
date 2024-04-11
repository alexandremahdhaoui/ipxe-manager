package controller

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/google/uuid"
)

var (
	ErrConfigNotFound = errors.New("config cannot be found")
	ErrConfigGetById  = errors.New("getting config by id")

	errUUIDCannotBeNil = errors.New("uuid cannot be nil")
)

// ---------------------------------------------------- INTERFACE --------------------------------------------------- //

// TODO: RENAME CONFIG TO CONTENT.

type Content interface {
	GetByID(ctx context.Context, contentID uuid.UUID, attributes types.IpxeSelectors) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewConfig(profile adapter.Profile, mux ResolveTransformerMux) Content {
	return &config{
		profile: profile,
		mux:     mux,
	}
}

// ---------------------------------------------------- CONFIG ------------------------------------------------- //

type config struct {
	profile adapter.Profile
	mux     ResolveTransformerMux
}

func (c *config) GetByID(
	ctx context.Context,
	contentID uuid.UUID,
	attributes types.IpxeSelectors,
) ([]byte, error) {
	if contentID == uuid.Nil {
		return nil, errors.Join(errUUIDCannotBeNil, ErrConfigGetById)
	}

	list, err := c.profile.ListByConfigID(ctx, contentID)
	if errors.Is(err, adapter.ErrProfileNotFound) || len(list) == 0 {
		return nil, errors.Join(err, ErrConfigNotFound, ErrConfigGetById)
	}

	contentName := list[0].ContentIDToNameMap[contentID]
	content := list[0].AdditionalContent[contentName]
	// TODO: create `mux.ResolveAndTransform()`.
	// TODO: to choose b/w template exposed-content as a URL with ID or resolving+transforming:
	//       - add a boolean param to `mux` to either return a URL or a template.
	//       - add a parameter to `mux` for the baseURL.
	out, err := c.mux.ResolveAndTransform(ctx, content, types.IpxeSelectors{
		UUID:      contentID, // the contentID is authoritative! It should always overwrite the attribute uuid.
		Buildarch: attributes.Buildarch,
	})
	if err != nil {
		return nil, errors.Join(err, ErrConfigGetById)
	}

	return out, nil
}
