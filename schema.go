package schema

import (
	"reflect"
	"strings"
)

//example:
//type TestType struct {
//	Location *types.Location
//	Images   []*types.File
//}
//
//func TestGetSchema(t *testing.T) {
//	s := TestType{
//		Location: &types.Location{},
//		Images:   []*types.File{{}},
//	}
//
//	schemaItem := GetSchema(reflect.ValueOf(s), "", "")
//	schemaItemBytes, _ := json.Marshal(schemaItem)
//	fmt.Println(string(schemaItemBytes))
//}

type Schema struct {
	Type           string             `json:"type,omitempty"`
	Label          string             `json:"label,omitempty"`
	Properties     map[string]*Schema `json:"properties,omitempty"`
	Items          *Schema            `json:"items,omitempty"`
	Required       bool               `json:"required,omitempty"`
	WidgetSettings *WidgetSettings    `json:"widgetSettings,omitempty"`
}

type WidgetSettings struct {
	Name       string            `json:"name,omitempty"`
	Options    map[string]string `json:"options,omitempty"`
	URLPrefix  string            `json:"URLPrefix,omitempty"`
	Storage    string            `json:"storage,omitempty"`
	Images     bool              `json:"images,omitempty"`
	Vocabulary string            `json:"vocabulary,omitempty"`
}

func GetSchema(v reflect.Value) *Schema {
	return getSchema(v, "", "")
}

func getSchema(v reflect.Value, label, vocabulary string) *Schema {
	if v.IsValid() {
		typeOfS := v.Type()
		kind := v.Kind()

		if kind == reflect.Int64 || kind == reflect.Float64 || kind == reflect.String || kind == reflect.Bool {
			typeName := typeOfS.String()
			return &Schema{
				Type:  typeName,
				Label: label,
				WidgetSettings: &WidgetSettings{
					Vocabulary: vocabulary,
				},
			}
		} else if kind == reflect.Map {
		} else if kind == reflect.Ptr {
			return getSchema(v.Elem(), label, vocabulary)
		} else if kind == reflect.Array || kind == reflect.Slice {
			if v.Len() > 0 {
				schema := &Schema{
					Type:  "array",
					Label: label,
					Items: getSchema(v.Index(0), "", ""),
				}
				return schema
			}
		} else if kind == reflect.Struct {
			fieldsCount := v.NumField()
			schema := &Schema{
				Type:       "object",
				Label:      label,
				Properties: map[string]*Schema{},
			}
			for i := 0; i < fieldsCount; i++ {
				f := typeOfS.Field(i)
				fieldName := f.Name
				fieldNameTag := strings.Split(f.Tag.Get("json"), ",")
				if len(fieldNameTag) > 0 && fieldNameTag[0] != "" {
					fieldName = fieldNameTag[0]
				}
				label := f.Tag.Get("label")
				vocabulary := f.Tag.Get("vocabulary")
				schema.Properties[fieldName] = getSchema(v.Field(i), label, vocabulary)
			}
			return schema
		}
	}
	return nil
}
