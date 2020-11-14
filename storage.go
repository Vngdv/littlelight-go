package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	storageFileName = "storage.json"
)

var (
	storage Storage
)

// This should be saved in a Database or something keeping everything in RAM is inefficient
type Storage struct {
	Guilds       map[string]Guild
	ChannelNames map[string]string
}

type Guild struct {
	ChannelCategory string
	CreationChannel string
}

func loadStorage() {

	// read file
	data, err := ioutil.ReadFile(storageFileName)
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data, &storage)
	if err != nil {
		fmt.Println("error:", err)
	}
}
