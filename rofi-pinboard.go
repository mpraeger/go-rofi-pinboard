package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/umahmood/pinboard"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	version = "0.1.3"
)

// Setup rofi-pinboard
func setupRofiPinboard(c *cli.Context) error {
	var cfg Config
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Pinboard API Token: ")
	apiToken, _ := reader.ReadString('\n')
	cfg.APIToken = strings.TrimSuffix(apiToken, "\n")
	// fmt.Print("\033[H\033[2J")

	// set Tempfile
	tmpdir := os.TempDir()
	cfg.Tempfile = fmt.Sprintf("%s/bookmarks.pinboard", tmpdir)

	// save config
	err := cfg.saveConfig("")

	if err != nil {
		log.Fatalf("Could not save config: %s", err)
	}

	return nil
}

func getURL(c *cli.Context) error {
	var id int
	var err error
	var pbk []PinboardBookmark

	id, err = strconv.Atoi(c.Args().First())
	if err != nil {
		log.Fatalf("Could not convert string to integer: %s", err)
	}

	// load config
	var cfg Config
	cfg.loadConfig("")

	pbk, err = loadBookmarks(&cfg)
	if err != nil {
		log.Fatalf("Could not load bookmarks: %s", err)
	}

	fmt.Println(pbk[id].URL)

	return nil
}

func getAllBookmarks(c *cli.Context) error {
	// load config
	var cfg Config
	cfg.loadConfig("")

	pin := pinboard.New()

	_, err := pin.Auth(cfg.APIToken)
	if err != nil {
		log.Fatalf("Could not authenticate to pinboard API: %s", err)
	}

	// check if we're authenticated
	if pin.IsAuthed() != true {
		log.Fatalf("Could not authenticate to pinboard API.")
	}

	bmarks, err := pin.Bookmarks(nil, 0, 0, time.Time{}, time.Time{}, false)

	if err != nil {
		log.Fatalf("Could not get bookmarks: %s", err)
	}

	pbk := make([]PinboardBookmark, 0)

	// iterate bookmarks
	for id, b := range bmarks {
		var bm PinboardBookmark
		bm.URL = b.URL
		bm.Title = b.Title
		bm.Tags = strings.Join(b.Tags, ",")
		bm.Hash = b.Hash
		pbk = append(pbk, bm)
		fmt.Printf("%d %s %s\n", id, bm.URL, bm.Tags)
	}

	err = saveBookmarks(&cfg, &pbk)
	if err != nil {
		log.Fatalf("Could not save bookmarks: %s", err)
	}

	return nil
}

func main() {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print only the version",
	}
	app := cli.NewApp()
	app.Name = "go-rofi-pinboard"
	app.Version = version
	app.Description = "Search and open bookmarks via the Pinboard API for display with rofi"
	app.Usage = "Search and open bookmarks via the Pinboard API"
	app.Commands = []cli.Command{
		{
			Name:   "all",
			Usage:  "get all bookmarks from Pinboard API",
			Action: getAllBookmarks,
		},
		{
			Name:   "url",
			Usage:  "get url for a bookmark",
			Action: getURL,
		},
		{
			Name:   "setup",
			Usage:  "setup rofi-pinboard",
			Action: setupRofiPinboard,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
