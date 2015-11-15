package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"

	"./hackernews"
)

var configFilePath = "" // set in init()

var config = struct {
	ApiBaseUrl  string
	ItemsInList int
}{"https://hacker-news.firebaseio.com/v0/", 10}

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
		fmt.Printf("Config file %q not found. Will try to create one for you.\n", configFilePath)
		err = createConfigFile()
		if err != nil {
			fmt.Printf("Could not create config file %q: %v\n", configFilePath, err.Error())
			os.Exit(1)
		}

		return
	}

	// Unmarshal into config
	err = json.Unmarshal(bs, &config)
	if err != nil {
		fmt.Printf("Error parsing config file %q: %v\n", configFilePath, err.Error())
		os.Exit(2)
	}
}

func main() {

	if len(os.Args) < 2 {
		showTopStories()
		return
	}

	i, err := getIndexFromArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	showStoryByIndex(i)

}

func showTopStories() {

	fmt.Println("Loading Hacker News Top Stories...")

	hn, _ := hackernews.NewHackerNews(config.ApiBaseUrl)

	stories, err := hn.GetTopStories(config.ItemsInList)
	if err != nil {
		fmt.Println("Could not list top stories")
		fmt.Printf("Error: %v\n", err.Error())
	}

	for _, story := range stories {
		fmt.Printf("%2v: %v\n", story.I, story.Title)
	}
}

func showStoryByIndex(index int) {

	hn, _ := hackernews.NewHackerNews(config.ApiBaseUrl)

	story, err := hn.GetStoryByIndex(index)
	if err != nil {
		fmt.Println("Could not get story")
		fmt.Printf("Error: %v\n", err.Error())
	}

	fmt.Printf("%#v\n", story)
}

// Parses the first command line argument as an integer and
// checks if it's inrange (between 0 and ItemsInList)
func getIndexFromArgs() (int, error) {

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return -1, errors.New("Argument not a number: " + err.Error())
	}

	if i < 0 || i >= config.ItemsInList {
		return -1, errors.New("Index out of range")
	}

	return i, nil
}

// Saves config as JSON to the config file
func createConfigFile() error {

	bs, err := json.Marshal(config)
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
