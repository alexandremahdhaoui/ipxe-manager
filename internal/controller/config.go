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

type Config interface {
	GetByID(ctx context.Context, profileName string, configID uuid.UUID) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewConfig(profile adapter.Profile, mux ResolveTransformerMux) Config {
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

func (c *config) GetByID(ctx context.Context, profileName string, configID uuid.UUID) ([]byte, error) {
	if configID == uuid.Nil {
		return nil, errors.Join(errUUIDCannotBeNil, ErrConfigGetById)
	}

	profile, err := c.profile.Get(ctx, profileName)
	if err != nil {
		return nil, errors.Join(err, ErrConfigGetById)
	}

	var (
		content types.Content
		found   bool
	)

	// This can be constant time by building a map if we let users name the additional content.
	// Constant time can also be achieved by building that map during the v1alpha1-to-types conversion.
	// NB: these comments are pointless because we don't expect len(profile.AdditionalContent) to be big.
	for _, ctt := range profile.AdditionalContent {
		if ctt.ExposedConfigID == configID {
			content = ctt
			found = true
			break
		}
	}

	if !found {
		return nil, errors.Join(ErrConfigNotFound, ErrConfigGetById)
	}

	res, err := c.mux.ResolveAndTransformBatch(ctx, []types.Content{content}, types.IpxeSelectors{
		UUID: configID,
		//TODO: allow config.GetByID to accept types.IPXESelectors as argument.
		// Buildarch: "",
	})
	if err != nil {
		return nil, errors.Join(err, ErrConfigGetById)
	}

	return res[content.Name], nil
}
