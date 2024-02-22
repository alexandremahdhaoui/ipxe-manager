package adapter

import (
	"context"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// --------------------------------------------------- INTERFACE ---------------------------------------------------- //

type Resolver interface {
	Resolve(ctx context.Context, c types.Content) ([]byte, error)
}

// ------------------------------------------------- INLINE RESOLVER ------------------------------------------------ //

type inlineResolver struct{}

func (r *inlineResolver) Resolve(ctx context.Context, c types.Content) ([]byte, error) {
	return []byte(c.Inline), nil
}

func NewInlineResolver() Resolver {
	return &inlineResolver{}
}

// ---------------------------------------------- OBJECT REF RESOLVER ----------------------------------------------- //

type objectRefResolver struct {
	client client.Client
}

func (r *objectRefResolver) Resolve(ctx context.Context, c types.Content) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func NewObjectRefResolver(c client.Client) Resolver {
	return &objectRefResolver{client: c}
}

// ------------------------------------------------ WEBHOOK RESOLVER ------------------------------------------------ //

type webhookResolver struct {
	client http.Client
}

func (r *webhookResolver) Resolve(ctx context.Context, c types.Content) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func NewWebhookResolver(c http.Client) Resolver {
	return &webhookResolver{client: c}
}
