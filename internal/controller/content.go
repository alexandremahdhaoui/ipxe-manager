package controller

import (
	"context"
	"errors"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/google/uuid"
)

var (
	ErrContentNotFound = errors.New("content cannot be found")
	ErrContentGetById  = errors.New("getting content by id")

	errUUIDCannotBeNil = errors.New("uuid cannot be nil")
)

// ---------------------------------------------------- INTERFACE --------------------------------------------------- //

type Content interface {
	GetByID(ctx context.Context, contentID uuid.UUID, attributes types.IpxeSelectors) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewConfig(profile adapter.Profile, mux ResolveTransformerMux) Content {
	return &content{
		profile: profile,
		mux:     mux,
	}
}

// ---------------------------------------------------- CONTENT ----------------------------------------------------- //

type content struct {
	profile adapter.Profile
	mux     ResolveTransformerMux
}

func (c *content) GetByID(
	ctx context.Context,
	contentID uuid.UUID,
	attributes types.IpxeSelectors,
) ([]byte, error) {
	if contentID == uuid.Nil {
		return nil, errors.Join(errUUIDCannotBeNil, ErrContentGetById)
	}

	list, err := c.profile.ListByContentID(ctx, contentID)
	if errors.Is(err, adapter.ErrProfileNotFound) || len(list) == 0 {
		return nil, errors.Join(err, ErrContentNotFound, ErrContentGetById)
	}

	contentName := list[0].ContentIDToNameMap[contentID]
	cont := list[0].AdditionalContent[contentName]
	// TODO: create `mux.ResolveAndTransform()`.
	// TODO: to choose b/w template exposed-content as a URL with ID or resolving+transforming:
	//       - add a boolean param to `mux` to either return a URL or a template.
	//       - add a parameter to `mux` for the baseURL.
	out, err := c.mux.ResolveAndTransform(ctx, cont, types.IpxeSelectors{
		UUID:      contentID, // the contentID takes precedence, thus should always overwrite the attribute uuid.
		Buildarch: attributes.Buildarch,
	})
	if err != nil {
		return nil, errors.Join(err, ErrContentGetById)
	}

	return out, nil
}
