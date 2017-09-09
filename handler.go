package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type LineNotifyHandler struct {
	AccessToken string
}

type PingdomReqBody struct {
	CheckID               int    `json:"check_id"`
	CheckName             string `json:"check_name"`
	CheckType             string `json:"check_type"`
	PreviousState         string `json:"previous_state"`
	CurrentState          string `json:"current_state"`
	ImportanceLevel       string `json:"importance_level"`
	StateChangedTimestamp int    `json:"state_changed_timestamp"`
	LongDescription       string `json:"long_description"`
	Description           string `json:"description"`
}

const LineNotifyApiURL = "https://notify-api.line.me/api/notify"

func (l LineNotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := new(PingdomReqBody)
	if err := json.Unmarshal(buf, body); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := NotifyLine(l.AccessToken, body); err != nil {
		log.Printf("[ERROR] Failed to LineNotify Post Request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NotifyLine(token string, p *PingdomReqBody) error {
	message := "\n" + p.CheckName + "\n\n" +
		p.LongDescription + "\n\n" +
		p.Description + "\n\n" +
		"状態: " + p.PreviousState + " -> " + p.CurrentState

	values := url.Values{"message": {message}}
	req, err := http.NewRequest("POST", LineNotifyApiURL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	fmt.Println(resp)
	return nil
}
