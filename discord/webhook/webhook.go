package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type request struct {
	url     string
	message Message
}

var webhook = make(chan request)

func init() {
	go func() {
		for {
			select {
			case w := <-webhook:
				send(w)
			}
		}
	}()
}

func send(w request) {
	enc, err := json.Marshal(w.message)
	if err != nil {
		fmt.Println("[Webhook API]: Error marshaling json: " + err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.url, bytes.NewReader(enc))
	if err != nil {
		fmt.Println("[Webhook API]: Error making post request: " + err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[Webhook API]: Error sending post request: " + err.Error())
		return
	}
	_ = resp.Body.Close()
}
