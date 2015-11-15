package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"./hackernews"
)

var config = struct {
	itemsInList int
	hnBaseUrl   string
}{10, "https://hacker-news.firebaseio.com/v0/"}

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

	hn, _ := hackernews.NewHackerNews(config.hnBaseUrl)

	stories, err := hn.GetTopStories(config.itemsInList)
	if err != nil {
		fmt.Println("Could not list top stories")
		fmt.Printf("Error: %v\n", err.Error())
	}

	for _, story := range stories {
		fmt.Printf("%2v: %v\n", story.I, story.Title)
	}
}

func showStoryByIndex(index int) {

	hn, _ := hackernews.NewHackerNews(config.hnBaseUrl)

	story, err := hn.GetStoryByIndex(index)
	if err != nil {
		fmt.Println("Could not get story")
		fmt.Printf("Error: %v\n", err.Error())
	}

	fmt.Printf("%#v\n", story)
}

func getIndexFromArgs() (int, error) {

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return -1, errors.New("Argument not a number: " + err.Error())
	}

	if i < 0 || i >= config.itemsInList {
		return -1, errors.New("Index out of range")
	}

	return i, nil
}
