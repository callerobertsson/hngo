package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"./hackernews"
)

// Path to the config file, set in init() to expanded "~/.hngorc"
var configFilePath = ""

// Configuration with default values
// Will be overridden when the real config is read
// Or saved as default if no config exist
var config = hackernews.Config{
	ApiBaseUrl:        "https://hacker-news.firebaseio.com/v0/",
	ItemsLimit:        10,
	CacheFilePath:     "/tmp/hngocache",
	OpenCommand:       []string{"echo"},
	ShowCommandOutput: true,
}

// Initialization
func init() {

	// Set the config file path
	maybePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Could not get path to config file: %v\n", err.Error())
		os.Exit(1)
	}
	configFilePath = maybePath

	// Read config file
	bs, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		// On error try to create a new config file with default values
		fmt.Printf("Config file %q not found.\n", configFilePath)
		err = createConfigFile()
		if err != nil {
			fmt.Printf("Could not create config file %q: %v\n", configFilePath, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Created config file %q with default values\n", configFilePath)

		// No need to continue, the default config is already there
		return
	}

	// Unmarshal into config
	err = json.Unmarshal(bs, &config)
	if err != nil {
		fmt.Printf("Error parsing config file %q: %v\n", configFilePath, err.Error())
		os.Exit(2)
	}
}

// Primus Motor
func main() {

	// Default if no args
	if len(os.Args) < 2 {
		showTopStories()
		return
	}

	// Try to interpret the first arg as an index
	i, err := getIndexFromArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	openStoryByIndex(i)
}

// Fetch a list of top stories from Hacker News
func showTopStories() {

	fmt.Println("Loading Hacker News Top Stories...")

	hn := hackernews.New(config)

	stories, err := hn.GetTopStories()
	if err != nil {
		fmt.Println("Could not list top stories")
		fmt.Printf("Error: %v\n", err.Error())
		return
	}

	for i, story := range stories {
		fmt.Printf("%2v: %v\n", i, story.Title)
	}
}

// Use OpenCommand to open the Story with index number
func openStoryByIndex(index int) {

	hn := hackernews.New(config)

	story, err := hn.GetCachedStoryByIndex(index)
	if err != nil {
		fmt.Println("Could not get story")
		fmt.Printf("Error: %v\n", err.Error())
		return
	}

	args := makeCommandArgs(config.OpenCommand, story.Url)
	showHide := "hiding"
	if config.ShowCommandOutput {
		showHide = "showing"
	}

	fmt.Printf("Story: %v\n  Date: %v\n  Url: %v\n", story.Title, story.Date, story.Url)
	fmt.Printf("Command: %v (%v output)\n\n", strings.Join(args, " "), showHide)

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Problems opening command")
		return
	}

	if cmd.Start() != nil {
		fmt.Printf("Error: could not execute command %q\n", strings.Join(args, " "), story.Url)
	}

	if config.ShowCommandOutput {
		bs, _ := ioutil.ReadAll(out)
		fmt.Println(string(bs))
	}
}

// Parses the first command line argument as an integer and
// checks if it's inrange (between 0 and ItemsInList)
func getIndexFromArgs() (int, error) {

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return -1, errors.New("Argument not a number: " + err.Error())
	}

	if i < 0 || i >= config.ItemsLimit {
		return -1, errors.New("Index out of range")
	}

	return i, nil
}

// Saves config as JSON to the config file
func createConfigFile() error {

	bs, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Println("Affen!")
		return err
	}

	return ioutil.WriteFile(configFilePath, bs, 0644)
}

// Used to figure out the path to the configuration file
func getConfigFilePath() (string, error) {

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir + "/.hngorc", nil
}

func makeCommandArgs(args []string, url string) []string {
	// TODO: Replace place holder with url
	return append(args, url)
}
