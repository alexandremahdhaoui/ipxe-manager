package adapters

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	butaneconfig "github.com/coreos/butane/config"
	butanecommon "github.com/coreos/butane/config/common"
)

var (
	ErrTransformerTransform = errors.New("transforming content")
)

// --------------------------------------------------- INTERFACE ---------------------------------------------------- //

type Transformer interface {
	Transform(ctx context.Context, content []byte, cfg types.TransformerConfig) ([]byte, error)
}

// ----------------------------------------------- BUTANE TRANSFORMER ----------------------------------------------- //

func NewButaneTransformer() Transformer {
	return &butaneTransformer{}
}

type butaneTransformer struct{}

func (t *butaneTransformer) Transform(
	_ context.Context,
	content []byte,
	_ types.TransformerConfig,
) ([]byte, error) {
	b, _, err := butaneconfig.TranslateBytes(content, butanecommon.TranslateBytesOptions{Raw: true})
	if err != nil {
		return nil, errors.Join(err, ErrTransformerTransform)
	}

	return b, nil
}

// ---------------------------------------------- WEBHOOK TRANSFORMER ----------------------------------------------- //

func NewWebhookTransformer(resolver WebhookResolver) Transformer {
	return &webhookTransformer{webhook: resolver}
}

type webhookTransformer struct {
	webhook WebhookResolver
}

func (t *webhookTransformer) Transform(
	ctx context.Context,
	content []byte,
	cfg types.TransformerConfig,
) ([]byte, error) {
	//TODO: implement me
	panic("implement me")
}
