package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	// Thrift
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/samceena/data-schema-encoding-decoding/gen-go/event"
	"github.com/vmihailenco/msgpack/v5"

	// Protofbuf
	"github.com/samceena/data-schema-encoding-decoding/gen-go/eventpb"
	"google.golang.org/protobuf/proto"

	// Avro
	"github.com/linkedin/goavro/v2"
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

// for thrift
func toThriftEvent(e Event) *event.Event {
	return &event.Event{
		ID:        e.ID,
		Username:  e.Username,
		Action:    e.Action,
		Timestamp: e.Timestamp,
		Metadata:  e.Metadata,
	}
}

func encodeThrift(e Event) []byte {
	thriftEvent := toThriftEvent(e)
	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTBinaryProtocolConf(transport, nil)
	if err := thriftEvent.Write(context.Background(), protocol); err != nil {
		log.Fatal("thrift encode error: %w", err)
	}
	return transport.Bytes()
}

// for protobuff
func toProtoEvent(e Event) *eventpb.Event {
	return &eventpb.Event{
		Id:        e.ID,
		Username:  e.Username,
		Action:    e.Action,
		Timestamp: e.Timestamp,
		Metadata:  e.Metadata,
	}
}

func encodeProtobuff(e Event) []byte {
	protoEvent := toProtoEvent(e)
	data, err := proto.Marshal(protoEvent)
	if err != nil {
		log.Fatalf("Error encoding event to protobuff: %w", err)
	}
	return data
}

func encodeAvro(e Event) []byte {
	schema := `{
		"type": "record",
		"name": "Event",
		"namespace": "event",
		"fields": [
			{"name": "id", "type": "long"},
			{"name": "username", "type": "string"},
			{"name": "action", "type": "string"},
			{"name": "timestamp", "type": "long"},
			{"name": "metadata", "type": {"type": "map", "values": "string"}}
		]
	}`

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		log.Fatalf("Avro schema error : %w", err)
	}

	data, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"id":        e.ID,
		"username":  e.Username,
		"action":    e.Action,
		"timestamp": e.Timestamp,
		"metadata":  e.Metadata,
	})
	if err != nil {
		log.Fatalf("avro encode error: %w", err)
	}
	return data
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
	thrifted := encodeThrift(event)
	fmt.Println("----")

	// Protobuff
	protoBuffEncoded := encodeProtobuff(event)
	fmt.Println("protoBuffEncoded: ", protoBuffEncoded)
	fmt.Println("----")

	// Protobuff
	avro := encodeAvro(event)
	fmt.Println("avroEncoded: ", avro)
	fmt.Println("----")

	fmt.Println("sizes:")
	fmt.Println("Json Data Encoded size: ", len(jsonDataEncoded))
	fmt.Println("MessagePackData ENcoded size: ", len(messagePackDataENcoded))
	fmt.Println("Thrift ENcoded size: ", len(thrifted))
	fmt.Println("ProtoBuff ENcoded size: ", len(protoBuffEncoded))
	fmt.Println("Avro ENcoded size: ", len(avro))

}
