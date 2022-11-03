package kcschema

import "encoding/json"

type Value struct {
	Schema  Schema                 `json:"schema"`
	Payload map[string]interface{} `json:"payload"`
}

type SchemaField struct {
	Field    string        `json:"field,omitempty"`
	Type     string        `json:"type,omitempty"`
	Optional bool          `json:"optional"`
	Fields   []SchemaField `json:"fields,omitempty"`
}

type Schema struct {
	Type     string        `json:"type,omitempty"`
	Name     string        `json:"name,omitempty"`
	Optional bool          `json:"optional"`
	Fields   []SchemaField `json:"fields,omitempty"`
}

func (p Payload) ParseAsKCSchema() (StructuredPayload, error) {
	var m Value
	err := json.Unmarshal(p, &m)
	if err != nil {
		return nil, err
	}

	sp := make(map[string]Field)
	for _, s := range m.Schema.Fields {
		sp[s.Field] = *parseKCField(m.Payload[s.Field], s.Type)
	}

	return sp, nil
}

func parseKCField(v interface{}, t string) *Field {
	ft := mapKCTypeToType(t)
	switch ft {
	case MapField:
		return &Field{
			Type:  MapField,
			Value: parseField(v),
		}
	default:
		return &Field{
			Type:  ft,
			Value: v,
		}
	}
}

func mapKCTypeToType(t string) FieldType {
	switch t {
	case "struct":
		return MapField
	case "string":
		return StringField
	case "int", "int32", "int64":
		return IntField
	case "float32", "float64":
		return FloatField
	case "boolean":
		return BoolField
	default:
		return UnknownField
	}
}

func mapTypeToKCType(t FieldType) string {
	switch t {
	case MapField:
		return "struct"
	case StringField:
		return "string"
	case IntField:
		return "integer"
	case FloatField:
		return "double"
	case BoolField:
		return "boolean"
	default:
		return "unknown"
	}
}
func (sp StructuredPayload) AsKCSchemaJSON(name string) ([]byte, error) {
	sch := Schema{
		Type:     "struct",
		Name:     name,
		Optional: false,
		Fields:   []SchemaField{},
	}

	pay := map[string]interface{}{}
	for f, v := range sp {
		sch.Fields = append(sch.Fields, SchemaField{
			Field:    f,
			Type:     mapTypeToKCType(v.Type),
			Optional: false,
		})
		pay[f] = v.Value
	}

	val := Value{
		Schema:  sch,
		Payload: pay,
	}

	return json.Marshal(val)
}
