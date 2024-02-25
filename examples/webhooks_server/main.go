package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeTimer EventType = "Timer"
	//EventTypeNewMessage EventType = "NewMessage"
)

type Subscription struct {
	Url       string
	EventType EventType
}

func main() {
	subscriptions := make(map[EventType][]Subscription)

	http.HandleFunc("/createWebHook", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "POST" {
			http.Error(w, "NOT POST!", http.StatusBadRequest)
			return
		}
		var s Subscription
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Cannot parse data!", http.StatusBadRequest)
			return
		}

		fmt.Printf("New subscription: %s\n", s)

		subscriptions[s.EventType] = append(subscriptions[s.EventType], s)

		w.WriteHeader(http.StatusOK)
	})

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for t := range ticker.C {
			for _, s := range subscriptions[EventTypeTimer] {
				resp, err := http.Post(s.Url, "text/plain", strings.NewReader(t.String()))
				if err != nil {
					fmt.Printf("Can't send event to %s: %s\n", s, err)
					continue
				}
				resp.Body.Close()
			}
		}
	}()

	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		return
	}
}
