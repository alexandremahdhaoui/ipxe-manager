package adapter

import (
	"context"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"github.com/coreos/butane/config"
	_ "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	_ "github.com/coreos/butane/config/common"
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

func (t *butaneTransformer) Transform(ctx context.Context, content []byte, cfg types.TransformerConfig) ([]byte, error) {
	b, _, err := config.TranslateBytes(content, common.TranslateBytesOptions{Raw: true})
	if err != nil {
		return nil, err //TODO: wrap me
	}

	return b, nil
}

// ---------------------------------------------- WEBHOOK TRANSFORMER ----------------------------------------------- //

func NewWebhookTransformer() Transformer {
	return &webhookTransformer{}
}

type webhookTransformer struct{}

func (t *webhookTransformer) Transform(ctx context.Context, content []byte, cfg types.TransformerConfig) ([]byte, error) {
	//TODO: implement me
	panic("implement me")
}
