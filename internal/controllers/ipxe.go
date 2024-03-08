package controllers

import (
	"bytes"
	"context"
	"github.com/alexandremahdhaoui/ipxe-api/internal/adapter"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"text/template"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type IPXE interface {
	FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewIPXE(profile adapter.Profile, mux ResolveTransformerMux) IPXE {
	return &ipxe{
		profile: profile,
		mux:     mux,
	}
}

// -------------------------------------------------------- IPXE ---------------------------------------------------- //

type ipxe struct {
	profile adapter.Profile
	mux     ResolveTransformerMux
}

func (svc *ipxe) FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error) {
	p, err := svc.profile.FindBySelectors(ctx, selectors)
	if err != nil {
		return nil, err //TODO: wrap
	}

	data, err := svc.mux.ResolveAndTransformBatch(ctx, p.AdditionalContent)
	if err != nil {
		return nil, err //TODO: wrap
	}

	output, err := templateIPXEProfile(p.IPXETemplate, data)
	if err != nil {
		return nil, err //TODO: wrap
	}

	return output, nil
}

func templateIPXEProfile(ipxeTemplate string, data map[string][]byte) ([]byte, error) {
	tpl, err := template.New("").Parse(ipxeTemplate)
	if err != nil {
		return nil, err //TODO: wrap
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	if err := tpl.Execute(buf, data); err != nil {
		return nil, err //TODO: wrap
	}

	return buf.Bytes(), nil
}
