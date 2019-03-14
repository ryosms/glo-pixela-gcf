package cloudfunctions

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	var d GloWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		_, _ = fmt.Fprint(w, "Hello World!")
		return
	}
	if d.Action == "" {
		_, _ = fmt.Fprint(w, "Hello World!")
		return
	}
	_, _ = fmt.Fprint(w, html.EscapeString(d.Board.Name))
}

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
