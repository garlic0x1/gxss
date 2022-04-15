package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func parsePayloads(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]map[string]map[string][]string
	if err := yaml.Unmarshal(file, &data); err != nil {
		panic(err)
	}

	Payloads = data
}

func buildAlert(input Input, payload string) string {
	str := input.URL + "?"
	for _, key := range input.Keys {
		str += key + "=" + payload + "&"
	}
	return str
}

func buildPayload(input Input) string {
	str := input.URL + "?"
	for i, key := range input.Keys {
		str += key + "=" + fmt.Sprintf(Canary, i) + "&"
	}
	return str
}
