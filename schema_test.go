package schema

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

//example:
type TestType struct {
	Location   [2]float64         `label:"Location label" json:"location" weight:"123" validate:"required"`
	ExampleMap map[string]float64 `label:"exampleMap label" json:"exampleMap" widget:"custom-map-input,images,sortable,URLPrefix=/asdasd/,cols=5"`
}

func TestGet(t *testing.T) {
	s := TestType{
		Location:   [2]float64{1, 2},
		ExampleMap: map[string]float64{"z": 123},
	}

	schemaItem := Get(reflect.ValueOf(s))
	schemaItemBytes, _ := json.Marshal(schemaItem)
	log.Println(string(schemaItemBytes))
	//t.Fatal(string(schemaItemBytes))
}
