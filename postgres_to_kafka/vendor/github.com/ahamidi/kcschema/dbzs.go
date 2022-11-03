package kcschema

import (
	"encoding/json"
	"time"
)

type Operation string

const (
	Read   Operation = "r"
	Create           = "c"
	Update           = "u"
	Delete           = "d"
)

type DBZValue struct {
	Schema  Schema `json:"schema"`
	Payload struct {
		Before map[string]interface{} `json:"before,omitempty"`
		After  map[string]interface{} `json:"after,omitempty"`
	} `json:"payload,omitempty"`
}
type DBZPayload struct {
	Op          Operation   `json:"op,omitempty"`
	Patch       string      `json:"patch,omitempty"`
	Filter      string      `json:"filter,omitempty"`
	Source      interface{} `json:"source,omitempty"`
	Transaction interface{} `json:"transaction,omitempty"`
	TimestampMS time.Time   `json:"ts_ms"`
}

func (p Payload) ParseAsDBZSchema() (StructuredPayload, error) {
	var m DBZValue
	err := json.Unmarshal(p, &m)
	if err != nil {
		return nil, err
	}

	// typically only interested in the After field values
	val := m.Payload.After

	sp := make(map[string]Field)
	// need to find the "after" schema fields
	var afterSchema []SchemaField
	for _, schema := range m.Schema.Fields {
		if schema.Field == "after" {
			afterSchema = schema.Fields
		}
	}
	for _, s := range afterSchema {
		sp[s.Field] = *parseKCField(val[s.Field], s.Type)
	}

	return sp, nil
}
