package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxe-api/internal/adapter"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"text/template"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type IPXE interface {
	FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error)
}

type ResolveTransformerMux interface {
	ResolveAndTransformBatch(ctx context.Context, batch []types.Content) (map[string][]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewIPXE(profile adapter.Profile, mux ResolveTransformerMux) IPXE {
	return &ipxe{
		profile: profile,
		mux:     mux,
	}
}
func NewResolveTransformerMux(
	resolvers map[types.ResolverKind]adapter.Resolver,
	transformers map[types.TransformerKind]adapter.Transformer,
) ResolveTransformerMux {
	return &resolveTransformerMux{
		resolvers:    resolvers,
		transformers: transformers,
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

// ---------------------------------------------------- MULTIPLEXER ------------------------------------------------- //

type resolveTransformerMux struct {
	resolvers    map[types.ResolverKind]adapter.Resolver
	transformers map[types.TransformerKind]adapter.Transformer
}

func (r *resolveTransformerMux) ResolveAndTransformBatch(
	ctx context.Context,
	batch []types.Content,
) (map[string][]byte, error) {
	output := make(map[string][]byte)

	for _, c := range batch {
		resolver, ok := r.resolvers[c.ResolverKind]
		if !ok {
			return nil, errors.New("TODO") //TODO: err
		}

		result, err := resolver.Resolve(ctx, c)
		if err != nil {
			return nil, err //TODO: wrap
		}

		for _, t := range c.PostTransformers {
			transformer, ok := r.transformers[t.Kind]
			if !ok {
				return nil, errors.New("TODO") //TODO: err
			}

			result, err = transformer.Transform(ctx, result, t)
		}

		output[c.Name] = result
	}

	return output, nil
}
