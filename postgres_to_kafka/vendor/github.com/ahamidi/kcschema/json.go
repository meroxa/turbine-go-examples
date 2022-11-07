package kcschema

import (
	"encoding/json"
	"log"
	"reflect"
)

func (p Payload) ParseAsJSON() (StructuredPayload, error) {
	sp := make(map[string]Field)
	var m map[string]interface{}
	err := json.Unmarshal(p, &m)

	for f, v := range m {
		log.Printf("f: %s, v: %+v", f, v)
		log.Printf("parsed: %+v", parseField(v))
		sp[f] = *parseField(v)
	}
	return sp, err
}

func parseField(v interface{}) *Field {
	switch v.(type) {
	case map[string]interface{}:
		return &Field{
			Type:  MapField,
			Value: parseField(v),
		}
	case string:
		return &Field{
			Type:  StringField,
			Value: v.(string),
		}
	case int, int32, int64:
		return &Field{
			Type:  IntField,
			Value: v.(int),
		}
	case float32, float64:
		return &Field{
			Type:  FloatField,
			Value: v.(float64),
		}
	default:
		log.Printf("v: %+v, type: %s", v, reflect.TypeOf(v).String())
	}
	return nil
}
