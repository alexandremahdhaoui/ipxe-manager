package controllers

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/adapters"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/google/uuid"
)

// ---------------------------------------------------- INTERFACE --------------------------------------------------- //

type Config interface {
	GetByID(ctx context.Context, profileName string, configID uuid.UUID) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewConfig(profile adapters.Profile, mux ResolveTransformerMux) Config {
	return &config{
		profile: profile,
		mux:     mux,
	}
}

// ---------------------------------------------------- CONFIG ------------------------------------------------- //

type config struct {
	profile adapters.Profile
	mux     ResolveTransformerMux
}

func (c *config) GetByID(ctx context.Context, profileName string, configID uuid.UUID) ([]byte, error) {
	if configID == uuid.Nil {
		return nil, errors.New("TODO") //TODO: err
	}

	profile, err := c.profile.Get(ctx, profileName)
	if err != nil {
		return nil, err //TODO: wrap me
	}

	//TODO(ineffective): change this
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
		return nil, errors.New("TODO") //TODO: err
	}

	res, err := c.mux.ResolveAndTransformBatch(ctx, []types.Content{content})
	if err != nil {
		return nil, err //TODO: wrap me
	}

	return res[content.Name], nil
}
