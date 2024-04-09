package controller

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
)

var (
	ErrResolveAndTransformBatch = errors.New("batch resolving and transforming contents")

	ErrResolverUnknown    = errors.New("unknown resolver")
	ErrTransformerUnknown = errors.New("unknown transformer")
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type ResolveTransformerMux interface {
	ResolveAndTransformBatch(
		ctx context.Context,
		batch []types.Content,
		selectors types.IpxeSelectors,
	) (map[string][]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewResolveTransformerMux(
	resolvers map[types.ResolverKind]adapter.Resolver,
	transformers map[types.TransformerKind]adapter.Transformer,
) ResolveTransformerMux {
	return &resolveTransformerMux{
		resolvers:    resolvers,
		transformers: transformers,
	}
}

// ---------------------------------------------------- MULTIPLEXER ------------------------------------------------- //

type resolveTransformerMux struct {
	resolvers    map[types.ResolverKind]adapter.Resolver
	transformers map[types.TransformerKind]adapter.Transformer
}

//TODO: ResolveAndTransformBatch should return the URL corresponding to the ConfigID of the content if the content has
//      ExposedConfigID set to true. (only in the case that the func is called by controller.IPXE)
//      !!! Otherwise create a special func for controller.Config called ResolveAndTransform which only takes a
//          types.Content as an argument and fully compute the Resolve/Transformation.
//      !!! Then ResolveAndTransformBatch will only resolve and transform if types.Content.ExposedConfigID != true.

func (r *resolveTransformerMux) ResolveAndTransformBatch(
	ctx context.Context,
	batch []types.Content,
	selectors types.IpxeSelectors,
) (map[string][]byte, error) {
	output := make(map[string][]byte)

	for _, content := range batch {
		resolver, ok := r.resolvers[content.ResolverKind]
		if !ok {
			return nil, errors.Join(ErrResolverUnknown, ErrResolveAndTransformBatch)
		}

		result, err := resolver.Resolve(ctx, content, selectors)
		if err != nil {
			return nil, errors.Join(err, ErrResolveAndTransformBatch)
		}

		for _, transformerConfig := range content.PostTransformers {
			transformer, ok := r.transformers[transformerConfig.Kind]
			if !ok {
				return nil, errors.Join(ErrTransformerUnknown, ErrResolveAndTransformBatch)
			}

			result, err = transformer.Transform(ctx, transformerConfig, result, selectors)
			if err != nil {
				return nil, errors.Join(err, ErrResolveAndTransformBatch)
			}
		}

		output[content.Name] = result
	}

	return output, nil
}
