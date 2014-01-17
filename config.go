package main

import (
	"encoding/json"
	"io/ioutil"
)

type hookConfig map[string][]string

func parseConfig(filename string) (config hookConfig, err error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return
	}

	err = json.Unmarshal(content, &config)
	return
}
