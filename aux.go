package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
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

func parsePayloads(payloadsfile string, tagmapfile string) {
	payloads, err := ioutil.ReadFile(payloadsfile)
	tagmap, err := ioutil.ReadFile(tagmapfile)

	var data map[string]map[string]Handler
	err = yaml.Unmarshal(tagmap, &data)
	if err != nil {
		panic(err)
	}

	var data2 map[string][]string
	err = yaml.Unmarshal(payloads, &data2)
	if err != nil {
		panic(err)
	}

	TagMap = data
	Payloads = data2
}

func buildPayload(input Input) string {
	str := input.URL + "?"
	for i, key := range input.Keys {
		str += key + "=" + fmt.Sprintf(Canary, i) + "&"
	}
	return str
}
