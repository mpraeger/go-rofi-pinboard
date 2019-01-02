package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	APIToken string `yaml:"api_token"`
	Tempfile string `yaml:"temp_file"`
}

// PinboardBookmark struct used for storing bookmarks on disk
type PinboardBookmark struct {
	Title string `json:"desc"`
	URL   string `json:"url"`
	Tags  string `json:"tags"`
	Hash  []byte `json:"hash"`
}

// load the config
func (c *Config) loadConfig(file string) (*Config, error) {
	var cfgFile string
	var err error

	if len(file) <= 0 {
		cfgFile = getConfigFile()
	} else {
		cfgFile = file
	}

	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return c, fmt.Errorf("yamlFile.Get err: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return c, fmt.Errorf("cannot unmarshal data: %v", err)
	}

	return c, err
}

// save the config
func (c *Config) saveConfig(file string) error {
	var err error
	var cfgFile string
	if len(file) <= 0 {
		cfgFile = getConfigFile()
	} else {
		cfgFile = file
	}

	// marshal
	yamlMarsh, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// write yaml
	err = ioutil.WriteFile(cfgFile, yamlMarsh, 0640)
	if err != nil {
		return err
	}

	fmt.Printf("Writing config file %s.\n", cfgFile)

	return err
}
