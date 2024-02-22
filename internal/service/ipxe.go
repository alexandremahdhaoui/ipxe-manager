package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxe-api/internal/adapter"
	"text/template"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type IPXE interface {
	FindProfileAndRender(ctx context.Context, selectors adapter.IpxeSelectors) ([]byte, error)
}

type ResolveTransformerMux interface {
	ResolveAndTransformBatch(ctx context.Context, batch []adapter.Content) (map[string][]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewIPXE(profile adapter.Profile, mux ResolveTransformerMux) IPXE {
	return &ipxe{
		profile: profile,
		mux:     mux,
	}
}
func NewResolveTransformerMux(
	resolvers map[adapter.ContentResolverKind]adapter.Resolver,
	transformers map[adapter.TransformerKind]adapter.Transformer,
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

func (svc *ipxe) FindProfileAndRender(ctx context.Context, selectors adapter.IpxeSelectors) ([]byte, error) {
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
		return nil, err
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
	resolvers    map[adapter.ContentResolverKind]adapter.Resolver
	transformers map[adapter.TransformerKind]adapter.Transformer
}

func (r *resolveTransformerMux) ResolveAndTransformBatch(
	ctx context.Context,
	batch []adapter.Content,
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

//TODO:
// When recursively templating:
//   - First check if we have DAG.
//   - If any cycles is spotted we should not allow such an operation to be performed.
// Additionally, we should ensure that no v1alpha1.Profile Custom Resource can be created if there are cycles:
//   - We need to create a Validating webhook that checks for cycles.
// Finally, we also need a form of runtime invalidation mechanism for dynamic non-DAGs:
//   - Indeed, to ensure that the content of "v1alpha1.ArbitraryResources"- or "v1alpha1.WebhookContent"-
//     v1alpha1.AdditionalContent does not contain cycles.
//   - We can create a DAG upon requesting all those information however, we should use BFS in order to avoid infinite
//     cycles. A max depth when running BFS might be a good solution.

// Template
//
// There is a body containing references to other configs that can themselves contain references to other configs.
//
// Templating happens in 3 phases:
//   1. find references
//   2. resolve references
//   3. render the template
//
// After rendering the template, we can recursively search for any references to a template.
