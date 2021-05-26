package helpers

import (
	"log"

	"gopkg.in/yaml.v2"
)

func ConvertStructToString(in interface{}) string {
	d, err := yaml.Marshal(&in)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}

func ConvertStringToStruct(s string) (out interface{}) {
	err := yaml.Unmarshal([]byte(s), &out)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return out
}
