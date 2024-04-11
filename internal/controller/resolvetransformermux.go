package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
)

const (
	ipxerAPIConfigPath = "config"
)

var (
	ErrResolveAndTransform      = errors.New("resolve and transform content")
	ErrResolveAndTransformBatch = errors.New("resolve and transform batch")

	ErrResolverUnknown    = errors.New("unknown resolver")
	ErrTransformerUnknown = errors.New("unknown transformer")
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type ResolveTransformerMux interface {
	ResolveAndTransform(ctx context.Context, content types.Content, selectors types.IpxeSelectors) ([]byte, error)

	ResolveAndTransformBatch(
		ctx context.Context,
		batch map[string]types.Content,
		selectors types.IpxeSelectors,
		options ...resolveTransformBatchOption,
	) (map[string][]byte, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewResolveTransformerMux(
	ipxerBaseURL string,
	resolvers map[types.ResolverKind]adapter.Resolver,
	transformers map[types.TransformerKind]adapter.Transformer,
) ResolveTransformerMux {
	return &resolveTransformerMux{
		ipxerBaseURL: ipxerBaseURL,
		resolvers:    resolvers,
		transformers: transformers,
	}
}

// ---------------------------------------------------- MULTIPLEXER ------------------------------------------------- //

type resolveTransformerMux struct {
	resolvers    map[types.ResolverKind]adapter.Resolver
	transformers map[types.TransformerKind]adapter.Transformer

	ipxerBaseURL string
}

func (r *resolveTransformerMux) ResolveAndTransform(
	ctx context.Context,
	content types.Content,
	selectors types.IpxeSelectors,
) ([]byte, error) {
	resolver, ok := r.resolvers[content.ResolverKind]
	if !ok {
		return nil, errors.Join(ErrResolverUnknown, ErrResolveAndTransform)
	}

	out, err := resolver.Resolve(ctx, content, selectors)
	if err != nil {
		return nil, errors.Join(err, ErrResolveAndTransform)
	}

	for _, transformerConfig := range content.PostTransformers {
		transformer, ok := r.transformers[transformerConfig.Kind]
		if !ok {
			return nil, errors.Join(ErrTransformerUnknown, ErrResolveAndTransform)
		}

		out, err = transformer.Transform(ctx, transformerConfig, out, selectors)
		if err != nil {
			return nil, errors.Join(err, ErrResolveAndTransform)
		}
	}

	return out, nil
}

// -------------------------------------------------- ResolveAndTransformBatch -------------------------------------- //

//TODO: ResolveAndTransformBatch should return the URL corresponding to the ConfigID of the content if the content has
//      ExposedConfigID set to true. (only in the case that the func is called by controller.IPXE)
//      !!! Otherwise create a special func for controller.Content called ResolveAndTransform which only takes a
//          types.Content as an argument and fully compute the Resolve/Transformation.
//      !!! Then ResolveAndTransformBatch will only resolve and transform if types.Content.ExposedConfigID != true.

func (r *resolveTransformerMux) ResolveAndTransformBatch(
	ctx context.Context,
	batch map[string]types.Content,
	selectors types.IpxeSelectors,
	options ...resolveTransformBatchOption,
) (map[string][]byte, error) {
	opts := new(resolveTransformBatchOptions).apply(options...)

	output := make(map[string][]byte)

	for name, content := range batch {
		if opts.returnURLInsteadOfResolveAndTransform && content.Exposed {
			output[name] = []byte(fmt.Sprintf(
				"%s/%s/%s", r.ipxerBaseURL, ipxerAPIConfigPath, content.ExposedUUID.String()))
			continue
		}

		result, err := r.ResolveAndTransform(ctx, content, selectors)
		if err != nil {
			return nil, errors.Join(err, ErrResolveAndTransformBatch)
		}

		output[name] = result
	}

	return output, nil
}

type (
	resolveTransformBatchOptions struct {
		returnURLInsteadOfResolveAndTransform bool
	}

	resolveTransformBatchOption func(options *resolveTransformBatchOptions)
)

func (o *resolveTransformBatchOptions) apply(options ...resolveTransformBatchOption) *resolveTransformBatchOptions {
	for _, f := range options {
		f(o)
	}

	return o
}

// ReturnExposedContentURL will ensure resolvetransformermux.ResolveAndTransformBatch does not resolve and transform the
// content but return a URL to that content.
func ReturnExposedContentURL(options *resolveTransformBatchOptions) {
	options.returnURLInsteadOfResolveAndTransform = true
}
