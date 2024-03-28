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
	for _, ctt := range profile.AdditionalContent {
		if ctt.ExposedConfigID == content.ExposedConfigID {
			content = ctt
			found = true
			break
		}
	}

	if !found {
		return nil, errors.Join(ErrConfigNotFound, ErrConfigGetById)
	}

	res, err := c.mux.ResolveAndTransformBatch(ctx, []types.Content{content})
	if err != nil {
		return nil, errors.Join(err, ErrConfigGetById)
	}

	return res[content.Name], nil
}
