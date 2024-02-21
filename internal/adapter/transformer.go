package adapter

import (
	"github.com/coreos/butane/config"
	_ "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	_ "github.com/coreos/butane/config/common"
)

type TransformerType int

const (
	ButaneTransformerType TransformerType = iota
	WebhookTransformerType
)

type TransformerConfig struct {
	Type TransformerType

	Webhook *WebhookConfig
}

type Transformer interface {
	Transform(content []byte) ([]byte, error)
}

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

func NewWebhookTransformer() Transformer {
	return &webhookTransformer{}
}

type webhookTransformer struct{}

func (t *webhookTransformer) Transform(content []byte) ([]byte, error) {
	//TODO: implement me
	panic("implement me")
}
