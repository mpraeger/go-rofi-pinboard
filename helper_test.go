package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"reflect"
	"testing"
)

func Test_getConfigFile(t *testing.T) {
	// get current user
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	cfgfile := fmt.Sprintf("%s/.rofi-pinboard.yaml", usr.HomeDir)

	tests := []struct {
		name string
		want string
	}{
		{"cfgfile", cfgfile},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigFile(); got != tt.want {
				t.Errorf("getConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveBookmarks(t *testing.T) {

	type args struct {
		cfg       *Config
		bookmarks *[]PinboardBookmark
	}

	// configure for success
	var cfg Config

	// configure for failure
	var fcfg Config

	bookmarks := make([]PinboardBookmark, 0)

	// setup args
	var targs args
	targs.cfg = &cfg
	targs.bookmarks = &bookmarks

	var fargs args
	fargs.cfg = &fcfg
	fargs.bookmarks = &bookmarks

	cfg.APIToken = "test:0123456789"
	fcfg.APIToken = "test:0123456789"
	tmpfile, err := ioutil.TempFile("", "bookmarks.unittests")
	if err != nil {
		t.Fatalf("Could not open temp file: %s", err)
	}

	// properly clean up
	cfg.Tempfile = tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(cfg.Tempfile)

	// generate a random hash
	bmhash := make([]byte, 16)
	rand.Read(bmhash)

	for i := 0; i < 9; i++ {
		var bm PinboardBookmark
		bm.URL = "https://unit.test/"
		bm.Title = "Unit Test"
		bm.Tags = "unit,test"
		bm.Hash = bmhash
		bookmarks = append(bookmarks, bm)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"saveBookmarks",
			targs,
			false,
		},
		{
			"saveBookmarksFailure",
			fargs,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveBookmarks(tt.args.cfg, tt.args.bookmarks); (err != nil) != tt.wantErr {
				t.Errorf("saveBookmarks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_loadBookmarks(t *testing.T) {
	// vars
	type args struct {
		cfg *Config
	}

	// generate a random hash
	bmhash := make([]byte, 16)
	rand.Read(bmhash)

	// correct config
	var cfg Config

	// failing config
	var fcfg Config

	// wrong temp file
	var wcfg Config

	// setup args
	var targs args
	targs.cfg = &cfg

	// failing args
	var fargs args
	fargs.cfg = &fcfg

	// empty tempfile args
	var wargs args
	wargs.cfg = &wcfg

	bookmarks := make([]PinboardBookmark, 0)
	ebookmarks := make([]PinboardBookmark, 0)

	wcfg.Tempfile = fmt.Sprintf("/var/tmp/%s", bmhash)

	// setup
	cfg.APIToken = "test:0123456789"

	tmpfile, err := ioutil.TempFile("", "bookmarks.unittests")
	if err != nil {
		t.Fatalf("Could not open temp file: %s", err)
	}

	// properly clean up
	cfg.Tempfile = tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(cfg.Tempfile)

	for i := 0; i < 9; i++ {
		var bm PinboardBookmark
		bm.URL = "https://unit.test/"
		bm.Title = "Unit Test"
		bm.Tags = "unit,test"
		bm.Hash = bmhash
		bookmarks = append(bookmarks, bm)
	}

	err = saveBookmarks(&cfg, &bookmarks)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		want    []PinboardBookmark
		wantErr bool
	}{
		{
			"loadBookmarks",
			targs,
			bookmarks,
			false,
		},
		{
			"loadBookmarksFailure",
			fargs,
			ebookmarks,
			true,
		},
		{
			"loadBookmarksEmpty",
			wargs,
			ebookmarks,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadBookmarks(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBookmarks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadBookmarks() = %v, want %v", got, tt.want)
			}
		})
	}
}
