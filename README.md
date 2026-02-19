# Data Encoding & Decoding Compared

A small Go project that encodes the same message in multiple serialization formats and compares them across the dimensions that actually matter in distributed systems.

## Introduction

When two services in the backend need to communicate with each other, what encoding format should they use? JSON? MessagePack? Avro? Protobuf? The choice matters.

### Why?

The encoding format you pick directly affects:

- **Payload size**: smaller encodings mean less data over the wire and less storage cost.
- **Encode/decode speed**: some formats are significantly faster to serialize and deserialize.
- **Forward and backward compatibility**: when you roll out new deployments, old and new versions of your services will coexist. Forward compatibility means an older reader can still read data written by a newer writer (it ignores unknown fields). Backward compatibility means a newer reader can still read data written by an older writer (it handles missing fields with defaults). Not all formats support this equally.
- **Human readability**: text formats like JSON are easy to inspect and debug; binary formats are not.

### Example: a distributed job runner

Consider a distributed job runner with many workers. When you deploy a new version, not all workers update at the same time. During the rollout:

- New workers may produce messages that old workers need to read (requires **forward compatibility**).
- Old workers may have produced messages that new workers need to read (requires **backward compatibility**).

If your encoding format does not support schema evolution, adding or removing a field can break the entire system during deployment. This is why the encoding format matters, it determines how safely you can evolve your service contracts over time.

## Formats compared

- **JSON**: text-based, no formal schema, ubiquitous & verbose.
- **MessagePack**: binary JSON, more compact but offers no schema evolution guarantees beyond what JSON gives you.
- **Thrift**: binary with a required schema (.thrift), supports schema evolution via field tags.
- **Protobuf**: binary with a required schema (.proto), compact, strong schema evolution guarantees.
- **Avro**: binary with a required schema (.avsc), most compact because field names are omitted from the encoding entirely. Schema must be present at read time.


_____


## How to run it:
`go run main.go`
or 
`go run main.go --username jone`


to encode and decode thrift:
-  you need to install the thrift compiler, for mac i used: `brew install thrift`
- Create the thrift schema file event: `touch event.thrift` and add code:
```
 namespace go event

  struct Event {
    1: i64 id
    2: string username
    3: string action
    4: i64 timestamp
    5: map<string, string> metadata
  }
```

- Generate the golang code from the file above with: `thrift -r --gen go event.thrift`: 
gen-go/event/ will be created with the proper Go code.

- Add the code to the project:
```
        "github.com/apache/thrift/lib/go/thrift"
        "github.com/samceena/data-schema-encoding-decoding/gen-go/event"
```

Encode and decode Protobuff:
Protobuff works in 2 steps:
- define a schema in .proto file
- Generate go code from the schema using protoc compiler
the compiler reads the .proto file and generates Go code with serialization/deserialization methods built in.
You can also do all these manually, type it out, but it's quite a lot of work to do that for this demo.
Steps:
1. install
  ```
  brew install protobuf
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```
2. Check it's working:
```
  protoc --version
  protoc-gen-go --version
```
3. Gen go code `protoc --go_out=gen-go/eventpb --go_opt=paths=source_relative event.proto`



Encode and decode Avro:
- Avro does not require code generation — the `goavro` library handles encoding/decoding dynamically using the schema.
- Define the schema in a `.avsc` file (JSON format).
- The schema is passed inline to the codec at runtime.


Schemaless:
- JSON
- MessagePack

Schema based:
- Thrift
- Avro
- Protobuf

___

### Output / Results:
```
sizes:
Json Data Encoded size:  95
MessagePackData ENcoded size:  81
Thrift ENcoded size:  73
ProtoBuff ENcoded size:  37
Avro ENcoded size:  31

```