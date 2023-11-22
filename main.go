package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.000"
)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{WebhookURL: webhookURL}
}

func (sn *SlackNotifier) Notify(message string) error {
	if sn.WebhookURL == "" {
		return nil
	}

	payload := map[string]string{"text": message}
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(sn.WebhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func logMessage(message string) {
	currentTime := time.Now().Format(timeFormat)
	fmt.Printf("%s: %s\n----\n", currentTime, message)
}

func parsePort() int {
	portNum, _ := strconv.Atoi(os.Getenv("PORT"))
	if portNum == 0 {
		portNum = 8080
	}
	return portNum
}

func getClientIP(r *http.Request) string {
	ip := r.RemoteAddr
	if ipPortSeparator := strings.LastIndex(ip, ":"); ipPortSeparator != -1 {
		ip = ip[:ipPortSeparator]
	}
	return ip
}

func getRealClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	return getClientIP(r)
}

func echoHandler(notifier *SlackNotifier, authKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		response := struct {
			ClientIP string      `json:"client_ip"`
			Method   string      `json:"method"`
			Host     string      `json:"host"`
			Path     string      `json:"path"`
			Body     string      `json:"body,omitempty"`
			Params   url.Values  `json:"params,omitempty"`
		}{
			ClientIP: getRealClientIP(r),
			Method:   r.Method,
			Host:     r.Host,
			Path:     r.URL.Path,
			Params:   r.URL.Query(),
		}

		if authKey != "" {
			if response.Params.Get("authKey") != authKey {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("403 - Forbidden"))
				logMessage(fmt.Sprintf("Forbidden request from %s", response.ClientIP))
				return
			}
		}

		if r.Method == "POST" {
			response.Body = string(body)
		}

		responseJSON, _ := json.Marshal(response)
		prettyJSON, _ := json.MarshalIndent(response, "", "  ")
		logMessage(string(prettyJSON))

		format := "Web request received at *%s* from *%s* \n```%s```"
		if err := notifier.Notify(fmt.Sprintf(format, response.Host, response.ClientIP, string(prettyJSON))); err != nil {
			logMessage(fmt.Sprintf("Error posting to Slack: %s", err))
		} else {
			logMessage("Posted to Slack")
		}

		w.Write(responseJSON)
	}
}

func main() {
	authKey := flag.String("authKey", os.Getenv("AUTH_KEY"), "add authentication key parameter to requests")
	port := flag.Int("port", parsePort(), "port to listen on")
	flag.Parse()

	notifier := NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL"))
	http.HandleFunc("/", echoHandler(notifier, *authKey))

	address := fmt.Sprintf(":%d", *port)
	logMessage(fmt.Sprintf("Starting echo-server at %s with the authKey '%s'", address, *authKey))
	if err := http.ListenAndServe(address, nil); err != nil {
		logMessage(fmt.Sprintf("Server failed to start: %s", err))
	}
}
