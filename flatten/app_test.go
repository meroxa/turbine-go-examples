package main

import (
	"github.com/meroxa/turbine-go"
	"testing"
)

func TestFlattenTransform(t *testing.T) {
	r := turbine.Record{
		Key:     "1",
		Payload: []byte(`{"id": 1, "user": {"id": 100, "name": "alice", "email": "alice@example.com"}, "actions": ["register", "purchase"]}`),
	}

	out := Flatten{}.Process([]turbine.Record{r})

	payload, err := out[0].Payload.Map()
	if err != nil {
		t.Fatalf("want no error, got %s", err.Error())
	}

	if _, ok := payload["user.id"]; !ok {
		t.Fatalf("want user.id to exist, got missing: %s", string(out[0].Payload))
	}

	if _, ok := payload["actions.1"]; !ok {
		t.Fatalf("want actions.1 to exist, got missing: %s", string(out[0].Payload))
	}
}
