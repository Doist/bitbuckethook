package main

import (
	"encoding/json"
)

type Payload struct {
	CanonURL   string     `json:"canon_url"`
	Commits    []commit   `json:"commits"`
	Repository repository `json:"repository"`
	User       string     `json:"user"`
}

type commit struct {
	Author       string   `json:"author"`
	Branch       string   `json:"branch"`
	Files        []file   `json:"files"`
	Message      string   `json:"message"`
	Node         string   `json:"node"`
	Parents      []string `json:"parents"`
	RawAuthor    string   `json:"raw_author"`
	RawNode      string   `json:"raw_node"`
	Timestamp    string   `json:"timestamp"`
	UTCTimestamp string   `json:"utc_timestamp"`
}

type file struct {
	File string `json:"file"`
	Type string `json:"type"`
}

type repository struct {
	AbsoluteUrl string `json:"absolute_url"`
	IsFork      bool   `json:"fork"`
	IsPrivate   bool   `json:"is_private"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	SCM         string `json:"scm"`
	Slug        string `json:"slug"`
	Website     string `json:"website"`
}

func parsePayload(body string) (payload Payload, err error) {
	err = json.Unmarshal([]byte(body), &payload)
	return
}
