package adapter

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
	//TODO: implement me
	panic("implement me")
}

func NewWebhookTransformer() Transformer {
	return &webhookTransformer{}
}

type webhookTransformer struct{}

func (t *webhookTransformer) Transform(content []byte) ([]byte, error) {
	//TODO: implement me
	panic("implement me")
}
