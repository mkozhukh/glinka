package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Link struct {
	Raw    string
	Global string
	Status LinkStatus
}

type LinksStore struct {
	Records map[string]LinkRecord
}

func NewLinksStore() *LinksStore {
	store := LinksStore{}
	store.Records = make(map[string]LinkRecord)
	return &store
}

type LinkRecord struct {
	URL    string
	Count  int
	Status LinkStatus
	Links  []Link
	Error  string
}

type LinkStatus int

const (
	StatusOK LinkStatus = iota
	StatusError
	StatusBinary
	StatusExternal
	StatusMixedContent
)

func (s LinksStore) save(file string) error {
	data, err := yaml.Marshal(&s)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0744)
}

func (s LinksStore) load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(data), &s)
}
