package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

func buildUrl(context Context, payload string) string {
	parsed, err := url.Parse(context.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := parsed.Query()
	q.Add(context.Key, payload)
	parsed.RawQuery = q.Encode()
	decoded, err := url.QueryUnescape(q.Encode())
	if err != nil {
		log.Println(err)
	}
	parsed.RawQuery = context.Prefix + decoded
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
