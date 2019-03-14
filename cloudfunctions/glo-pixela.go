package cloudfunctions

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func GloToPixela(w http.ResponseWriter, r *http.Request) {
	header := parseHeader(&r.Header)
	var payload GloWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println(fmt.Sprintf("Bad Request: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	et, err := decideExecType(header, payload)
	if err != nil {
		log.Println(fmt.Sprintf("Bad Request: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if et == none {
		return
	}

	req, err := generateRequest(et)
	if err != nil {
		log.Println(fmt.Sprintf("Internal Server Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := postRequest(req)
	if err != nil {
		log.Println(fmt.Sprintf("Internal Server Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(res)
	_, _ = fmt.Fprint(w, html.EscapeString(res))
}

func parseHeader(h *http.Header) GloWebhookHeader {
	return GloWebhookHeader{strings.ToLower(h.Get("x-gk-event"))}
}

func decideExecType(header GloWebhookHeader, payload GloWebhookPayload) (ExecType, error) {
	if header.Event != "cards" {
		return none, nil
	}
	boardId := os.Getenv("GLO_BOARD_ID")
	if boardId != "" && boardId != payload.Board.Id {
		return none, fmt.Errorf("id of board is not match(should: %v, get: %v)", boardId, payload.Board.Id)
	}
	switch strings.ToLower(payload.Action) {
	case "added", "copied", "unarchived", "moved_from_board":
		return increment, nil
	case "archived", "deleted", "moved_to_board":
		return decrement, nil
	default:
		// updated, reordered, moved_column, labels_updated, assignees_updated
		return none, nil
	}
}

func generateRequest(et ExecType) (*http.Request, error) {
	username := os.Getenv("PIXELA_USERNAME")
	var webhookHash string
	if et == increment {
		webhookHash = os.Getenv("PIXELA_INCREMENT_WEBHOOK_HASH")
	} else {
		webhookHash = os.Getenv("PIXELA_DECREMENT_WEBHOOK_HASH")
	}

	if username == "" || webhookHash == "" {
		return nil, fmt.Errorf("environment variables are not defined(username: %v, webhook hash: %v)", username, webhookHash)
	}

	url := fmt.Sprintf("https://pixe.la/v1/users/%s/webhooks/%s", username, webhookHash)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "0")

	return req, err
}

func postRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get response body : %s", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode > 299 {
		return "", fmt.Errorf("%s", string(b))
	}
	return string(b), nil
}
