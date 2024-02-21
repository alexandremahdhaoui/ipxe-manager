package adapter

import (
	"github.com/coreos/butane/config"
	_ "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	_ "github.com/coreos/butane/config/common"
)

type TransformerKind int

const (
	ButaneTransformerKind TransformerKind = iota
	WebhookTransformerKind
)

// ----------------------------------------------------- TYPES ------------------------------------------------------ //

type TransformerConfig struct {
	Kind TransformerKind

	Webhook *WebhookConfig
}

// --------------------------------------------------- INTERFACE ---------------------------------------------------- //

type Transformer interface {
	Transform(content []byte) ([]byte, error)
}

// ----------------------------------------------- BUTANE TRANSFORMER ----------------------------------------------- //

func NewButaneTransformer() Transformer {
	return &butaneTransformer{}
}

type butaneTransformer struct{}

func (t *butaneTransformer) Transform(content []byte) ([]byte, error) {
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

func (t *webhookTransformer) Transform(content []byte) ([]byte, error) {
	//TODO: implement me
	panic("implement me")
}
