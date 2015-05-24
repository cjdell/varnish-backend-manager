package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path"
)

type ConfigStore struct {
	entries  []*ConfigEntry
	basePath string
}

type ConfigEntry struct {
	Host    string
	Backend string
}

func NewConfigStore(basePath string) *ConfigStore {
	return &ConfigStore{
		entries:  make([]*ConfigEntry, 0, 0),
		basePath: basePath,
	}
}

func (config *ConfigStore) GetEntries() []*ConfigEntry {
	return config.entries
}

func (config *ConfigStore) SetEntry(entry *ConfigEntry) {
	config.DeleteEntry(entry.Host)

	config.entries = append(config.entries, entry)
}

func (config *ConfigStore) DeleteEntry(host string) {
	toDelete := -1

	for i, entry := range config.entries {
		if entry.Host == host {
			toDelete = i
		}
	}

	if toDelete != -1 {
		config.entries = append(config.entries[:toDelete], config.entries[toDelete+1:]...)
	}
}

func (config *ConfigStore) Load() {
	jsonFile := path.Join(config.basePath, "config.json")

	f, err := os.Open(jsonFile)

	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)

	err = dec.Decode(&config.entries)
	if err != nil {
		log.Fatal(err)
	}
}

func (config *ConfigStore) Save() {
	jsonFile := path.Join(config.basePath, "config.json")

	f, err := os.Create(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)

	err = enc.Encode(config.entries)
	if err != nil {
		log.Fatal(err)
	}

	w.Flush()
}
