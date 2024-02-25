package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type EventType string

const (
	EventTypeTimer EventType = "Timer"
)

type Subscription struct {
	Url       string
	EventType EventType
}

func main() {
	req, err := json.Marshal(Subscription{
		Url:       "http://localhost:8086/events/timer",
		EventType: EventTypeTimer,
	})
	if err != nil {
		fmt.Printf("Can't marshal subscription %s\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8085/createWebHook", "application/json", bytes.NewReader(req))
	if err != nil {
		fmt.Printf("Can't send subscription %s\n", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		return
	}

	http.HandleFunc("/events/timer", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buf := new(strings.Builder)
		_, err = io.Copy(buf, r.Body)
		if err != nil {
			fmt.Printf("Can't decode body %s\n", err)
		}
		fmt.Println(buf.String())
		w.WriteHeader(http.StatusOK)
	})

	err = http.ListenAndServe(":8086", nil)
	if err != nil {
		return
	}
}
