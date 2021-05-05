package schema

import (
	"reflect"
	"strconv"
	"strings"
)

type Schema struct {
	Type           string             `json:"type,omitempty"`
	Label          string             `json:"label,omitempty"`
	Properties     map[string]*Schema `json:"properties,omitempty"`
	Items          *Schema            `json:"items,omitempty"`
	Required       bool               `json:"required,omitempty"`
	Weight         int64              `json:"weight,omitempty"`
	WidgetSettings *WidgetSettings    `json:"widgetSettings,omitempty"` //for Object.assign to component
}

type WidgetSettings struct {
	Name       string            `json:"name,omitempty"`
	Options    map[string]string `json:"options,omitempty"`
	Storage    string            `json:"storage,omitempty"`
	URLPrefix  string            `json:"URLPrefix,omitempty"`
	Vocabulary string            `json:"vocabulary,omitempty"`
	Images     bool              `json:"images,omitempty"`
	Cols       int64             `json:"cols,omitempty"`
}

func Get(v reflect.Value) *Schema {
	return getSchema(v, map[string]string{})
}

func getSchema(v reflect.Value, tags map[string]string) *Schema {
	if v.IsValid() {
		typeOfS := v.Type()
		kind := v.Kind()

		var weight int64
		if tags["weight"] != "" {
			weight, _ = strconv.ParseInt(tags["weight"], 10, 64)
		}
		var required bool
		if tags["validate"] != "" {
			validations := strings.Split(tags["validate"], ",")
			for _, str := range validations {
				if str == "required" {
					required = true
				}
			}
		}

		widgetSettingsSplitted := strings.Split(tags["widget"], ",")
		widgetSettingsFromStrFlags := map[string]bool{}
		widgetSettingsFromStrKV := map[string]string{}
		for _, setting := range widgetSettingsSplitted { //for now only flags
			settingKV := strings.Split(setting, "=")
			if len(settingKV) == 1 { //flag
				widgetSettingsFromStrFlags[settingKV[0]] = true
			}
			if len(settingKV) == 2 { //key-value
				widgetSettingsFromStrKV[settingKV[0]] = settingKV[1]
			}
		}
		widgetSettings := &WidgetSettings{
			Name:       widgetSettingsSplitted[0],
			Vocabulary: tags["vocabulary"],
		}
		if widgetSettingsFromStrFlags["images"] {
			widgetSettings.Images = true
		}
		if widgetSettingsFromStrKV["URLPrefix"] != "" {
			widgetSettings.URLPrefix = widgetSettingsFromStrKV["URLPrefix"]
		}
		if widgetSettingsFromStrKV["vocabulary"] != "" {
			widgetSettings.Vocabulary = widgetSettingsFromStrKV["vocabulary"]
		}
		if widgetSettingsFromStrKV["cols"] != "" {
			cols, err := strconv.ParseInt(widgetSettingsFromStrKV["cols"], 10, 64)
			if err == nil {
				widgetSettings.Cols = cols
			}
		}

		if kind == reflect.Int64 || kind == reflect.Float64 || kind == reflect.String || kind == reflect.Bool {
			typeName := typeOfS.String()
			return &Schema{
				Type:           typeName,
				Label:          tags["label"],
				Required:       required,
				Weight:         weight,
				WidgetSettings: widgetSettings,
			}
		} else if kind == reflect.Map {
			keys := v.MapKeys()
			if len(keys) > 0 {
				return &Schema{
					Type:           "map",
					Label:          tags["label"],
					Required:       required,
					Weight:         weight,
					WidgetSettings: widgetSettings,
					Items:          getSchema(v.MapIndex(keys[0]), map[string]string{}),
				}
			}
		} else if kind == reflect.Ptr {
			return getSchema(v.Elem(), tags)
		} else if kind == reflect.Array || kind == reflect.Slice {
			if v.Len() > 0 {
				schema := &Schema{
					Type:           "array",
					Label:          tags["label"],
					Required:       required,
					Weight:         weight,
					WidgetSettings: widgetSettings,
					Items:          getSchema(v.Index(0), map[string]string{}),
				}
				return schema
			}
		} else if kind == reflect.Struct {
			fieldsCount := v.NumField()
			schema := &Schema{
				Type:           "object",
				Label:          tags["label"],
				Required:       required,
				Weight:         weight,
				WidgetSettings: widgetSettings,
				Properties:     map[string]*Schema{},
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
				widget := f.Tag.Get("widget")
				weight := f.Tag.Get("weight")
				validate := f.Tag.Get("validate")
				schema.Properties[fieldName] = getSchema(v.Field(i), map[string]string{
					"label":      label,
					"vocabulary": vocabulary,
					"widget":     widget,
					"weight":     weight,
					"validate":   validate,
				})
			}
			return schema
		}
	}
	return nil
}
