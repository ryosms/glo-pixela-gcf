package cloudfunctions

type ExecType int

const (
	increment ExecType = iota
	decrement
	none
)

type WebhookItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GloWebhookPayload struct {
	Action   string      `json:"action"`
	Board    WebhookItem `json:"board"`
	Sender   WebhookItem `json:"sender"`
	Card     WebhookItem `json:"card"`
	Sequence int         `json:"sequence"`
}

type GloWebhookHeader struct {
	Event string
}
