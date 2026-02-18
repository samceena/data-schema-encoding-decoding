 namespace go event

struct Event {
  1: i64 id
  2: string username
  3: string action
  4: i64 timestamp
  5: map<string, string> metadata
}