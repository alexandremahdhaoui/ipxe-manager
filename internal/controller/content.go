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
	GetByID(
		ctx context.Context,
		contentID uuid.UUID,
		attributes types.IPXESelectors,
	) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewContent(profile adapter.Profile, mux ResolveTransformerMux) Content {
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
	attributes types.IPXESelectors,
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
	// NB: mux.ResolveAndTransform will always render the content. Please call ResolveAndTransformBatch
	// with the mux.ReturnExposedContentURL option to return a URL instead.
	out, err := c.mux.ResolveAndTransform(ctx, cont, types.IPXESelectors{
		UUID:      contentID, // the contentID takes precedence, thus should always overwrite the attribute uuid.
		Buildarch: attributes.Buildarch,
	})
	if err != nil {
		return nil, errors.Join(err, ErrContentGetById)
	}

	return out, nil
}
