package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Event struct {
	ID        int64             `json:"id"`
	Username  string            `json:"username"`
	Action    string            `json:"action"`
	Timestamp int64             `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

func encodeJSON(event Event) []byte {
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("json encode error: %w", err)
	}
	return data
}

func encodeMessagePack(event Event) []byte {
	data, err := msgpack.Marshal(event)
	if err != nil {
		log.Fatalf("messagepack encode error: %w", err)
	}
	return data
}

func defaultEvent() Event {
	return Event{
		ID:        1,
		Username:  "samceena",
		Action:    "login",
		Timestamp: time.Now().Unix(),
		Metadata: map[string]string{
			"foo": "bar",
		},
	}
}

func main() {

	username := flag.String("username", "", "override default username")
	action := flag.String("action", "", "override default action")
	flag.Parse()

	event := defaultEvent()
	if *username != "" {
		event.Username = *username
	}
	if *action != "" {
		event.Action = *action
	}

	fmt.Printf("Event to encode: \n%+v\n", event)

	jsonDataEncoded := encodeJSON(event)
	messagePackDataENcoded := encodeMessagePack((event))

	fmt.Println("Encoded Data:")
	fmt.Println("jsonDataEncoded: ", jsonDataEncoded)
	fmt.Println("messagePackDataENcoded: ", messagePackDataENcoded)
	fmt.Println("----")

	fmt.Println("sizes:")
	fmt.Println("Json Data Encoded size: ", len(jsonDataEncoded))
	fmt.Println("MessagePackData ENcoded size: ", len(messagePackDataENcoded))

}
