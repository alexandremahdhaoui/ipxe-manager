package controllers

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/adapters"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type ResolveTransformerMux interface {
	ResolveAndTransformBatch(ctx context.Context, batch []types.Content) (map[string][]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewResolveTransformerMux(
	resolvers map[types.ResolverKind]adapters.Resolver,
	transformers map[types.TransformerKind]adapters.Transformer,
) ResolveTransformerMux {
	return &resolveTransformerMux{
		resolvers:    resolvers,
		transformers: transformers,
	}
}

// ---------------------------------------------------- MULTIPLEXER ------------------------------------------------- //

type resolveTransformerMux struct {
	resolvers    map[types.ResolverKind]adapters.Resolver
	transformers map[types.TransformerKind]adapters.Transformer
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
