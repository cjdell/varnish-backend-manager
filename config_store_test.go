package main

import (
	"testing"
)

func TestCanInitWithZeroEntries(t *testing.T) {
	config := NewConfigStore("./")

	entries := config.GetEntries()

	if len(entries) != 0 {
		t.Error(
			"For", entries,
			"expected", 0,
			"got", len(entries),
		)
	}
}

func TestCanWriteAnEntryAndSaveAndRead(t *testing.T) {
	config := NewConfigStore("./")

	config.SetEntry(&ConfigEntry{
		Host:    "example.com",
		Backend: "192.168.1.1",
	})

	config.Save()

	config = NewConfigStore("./")

	config.Load()

	entries := config.GetEntries()

	if len(entries) != 1 {
		t.Error(
			"For", entries,
			"expected", 1,
			"got", len(entries),
		)
	}
}

func TestCanWriteEntriesAndDeleteOne(t *testing.T) {
	config := NewConfigStore("./")

	config.SetEntry(&ConfigEntry{
		Host:    "example.com",
		Backend: "192.168.1.1",
	})
	config.SetEntry(&ConfigEntry{
		Host:    "a.example.com",
		Backend: "192.168.1.1",
	})
	config.SetEntry(&ConfigEntry{
		Host:    "b.example.com",
		Backend: "192.168.1.1",
	})

	config.DeleteEntry("a.example.com")

	entries := config.GetEntries()

	if len(entries) != 2 {
		t.Error(
			"For", entries,
			"expected", 2,
			"got", len(entries),
		)
	}
}

func TestCanWriteEntriesAndUpdateOne(t *testing.T) {
	config := NewConfigStore("./")

	config.SetEntry(&ConfigEntry{
		Host:    "example.com",
		Backend: "192.168.1.1",
	})
	config.SetEntry(&ConfigEntry{
		Host:    "a.example.com",
		Backend: "192.168.1.1",
	})
	config.SetEntry(&ConfigEntry{
		Host:    "b.example.com",
		Backend: "192.168.1.1",
	})

	config.SetEntry(&ConfigEntry{
		Host:    "a.example.com",
		Backend: "192.168.1.2",
	})

	entries := config.GetEntries()

	if len(entries) != 3 {
		t.Error(
			"For", entries,
			"expected", 3,
			"got", len(entries),
		)
	}
}
