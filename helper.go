package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
)

func getConfigFile() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cfgFile := fmt.Sprintf("%s/.rofi-pinboard.yaml", usr.HomeDir)
	return cfgFile
}

// Save the bookmarks
func saveBookmarks(cfg *Config, bookmarks *[]PinboardBookmark) error {
	jsonBookmarks, err := json.Marshal(bookmarks)
	if err != nil {
		return fmt.Errorf("could not json.Marshal bookmarks %s", err)
	}

	// write jsonBookmarks to file
	if len(cfg.Tempfile) <= 0 {
		return errors.New("tempfile is not configured")
	}

	// as they are personal bookmarks it is my preference to set
	// the permissions of the file to 0640
	err = ioutil.WriteFile(cfg.Tempfile, jsonBookmarks, 0640)
	return err
}

// load bookmarks and json.Unmarshal() them
func loadBookmarks(cfg *Config) ([]PinboardBookmark, error) {
	jsonBookmarks := make([]PinboardBookmark, 0)

	// if tempfile is not set exit early
	if len(cfg.Tempfile) <= 0 {
		return jsonBookmarks, errors.New("tempfile not set")
	}

	// read file
	pbk, err := ioutil.ReadFile(cfg.Tempfile)
	if err != nil {
		return jsonBookmarks, err
	}

	// unmarshal
	err = json.Unmarshal(pbk, &jsonBookmarks)
	if err != nil {
		return jsonBookmarks, err
	}

	// return
	return jsonBookmarks, err
}
