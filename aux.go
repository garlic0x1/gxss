package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var sm sync.Map

func isUniqueOutput(res Result) bool {
	str := res.Type + res.Message

	_, present := sm.Load(str)
	if present {
		return false
	}
	sm.Store(str, true)
	return true
}

func buildUrl(context Context, payload string) string {
	parsed, err := url.Parse(context.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := parsed.Query()
	q.Add(context.Key, context.Prefix+payload)
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func parsePayloads(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	d := yaml.NewDecoder(file)

	var data map[string]map[string]map[string][]string
	err = d.Decode(&data)
	if err != nil {
		panic(err)
	}

	var data2 map[string][]string
	err = d.Decode(&data2)
	if err != nil {
		panic(err)
	}

	Payloads = data
	AttrPayloads = data2
}

func buildPayload(input Input) string {
	str := input.URL + "?"
	for i, key := range input.Keys {
		str += key + "=" + fmt.Sprintf(Canary, i) + "&"
	}
	return str
}
