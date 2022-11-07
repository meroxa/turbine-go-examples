package kcschema

import (
	"encoding/json"
	"errors"
	"github.com/meroxa/turbine-go"
)

type FieldType string

const (
	StringField  FieldType = "string"
	IntField     FieldType = "int"
	FloatField   FieldType = "float"
	MapField     FieldType = "map"
	BoolField    FieldType = "bool"
	UnknownField FieldType = "unknown"
)

type FieldInterface interface {
	Get() interface{}
	Set(value FieldInterface) error
	Type() FieldType
}

type PayloadInterface interface {
	Get(field string) FieldInterface
	Set(field FieldInterface) error
	Type(field string) FieldType
}

type Payload turbine.Payload

type Field struct {
	Type  FieldType
	Value interface{}
}

type StructuredPayload map[string]Field

type PayloadType string

const (
	DebeziumCDCType      PayloadType = "debezium"
	KCJSONWithSchemaType PayloadType = "KCJsonWithSchema"
	JSONType             PayloadType = "json"
	RawType              PayloadType = "raw"
)

func (p Payload) Type() PayloadType {
	var m map[string]interface{}
	err := json.Unmarshal(p, &m)
	if err != nil {
		return RawType
	}

	if p, ok := m["payload"]; ok {
		if _, ok := p.(map[string]interface{})["before"]; ok {
			return DebeziumCDCType
		}
		if _, ok := p.(map[string]interface{})["after"]; ok {
			return DebeziumCDCType
		}
	}

	if _, ok := m["schema"]; ok {
		return KCJSONWithSchemaType
	}

	return JSONType
}

func Parse(p Payload) (StructuredPayload, error) {
	switch p.Type() {
	case JSONType:
		return p.ParseAsJSON()
	case KCJSONWithSchemaType:
		return p.ParseAsKCSchema()
	case DebeziumCDCType:
		return p.ParseAsDBZSchema()
	default:
		return nil, errors.New("unable to parse payload")
	}
}
