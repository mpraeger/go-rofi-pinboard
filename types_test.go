package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestConfig_saveConfig(t *testing.T) {
	// vars
	type fields struct {
		APIToken string
		Tempfile string
	}
	type args struct {
		file string
	}
	// setup

	var targs args
	targs.file = fmt.Sprintf("%s/go-rofi-pinboard.yaml", os.TempDir())

	tempfile := fmt.Sprintf("%s/go-rofi-pinboard.unittests", os.TempDir())
	defer os.Remove(tempfile)

	// tests
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"saveConfig",
			fields{
				APIToken: "unittest:0123456789",
				Tempfile: tempfile,
			},
			targs,
			false,
		},
		{
			"saveConfigFail",
			fields{
				APIToken: "unittest:0123456789",
				Tempfile: tempfile,
			},
			args{
				file: "/this/does/not/exist",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIToken: tt.fields.APIToken,
				Tempfile: tt.fields.Tempfile,
			}
			if err := c.saveConfig(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Config.saveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_loadConfig(t *testing.T) {
	// vars
	type fields struct {
		APIToken string
		Tempfile string
	}
	type args struct {
		file string
	}

	// setup
	tempfile := fmt.Sprintf("%s/go-rofi-pinboard.unittests", os.TempDir())
	cfgFile := fmt.Sprintf("%s/go-rofi-pinboard.yaml", os.TempDir())
	defer os.Remove(tempfile)
	defer os.Remove(cfgFile)

	// args
	var targs args
	targs.file = cfgFile

	// config
	var cfg Config
	cfg.APIToken = "unittests.0123456789"
	cfg.Tempfile = tempfile
	cfg.saveConfig(cfgFile)

	var ecfg Config
	ecfg.APIToken = "unittests.0123456789"

	// tests
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Config
		wantErr bool
	}{
		{
			"loadConfig",
			fields{
				APIToken: cfg.APIToken,
				Tempfile: cfg.Tempfile,
			},
			targs,
			&cfg,
			false,
		},
		{
			"loadConfigTempfileNil",
			fields{
				APIToken: cfg.APIToken,
				Tempfile: "",
			},
			args{
				file: "/this/does/not/exist",
			},
			&ecfg,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIToken: tt.fields.APIToken,
				Tempfile: tt.fields.Tempfile,
			}
			got, err := c.loadConfig(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
