package hackernews

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HackerNews struct {
	config Config
}

type Config struct {
	ApiBaseUrl    string
	ItemsLimit    int
	CacheFilePath string
	// TODO: Add CacheTimeToLiveSecs int
	// TODO: Add ParallellMode bool
}

type Story struct {
	I       int
	Id      int
	Title   string
	Date    string
	Content string
	Url     string
}

// Create a new HackerNews
func New(config Config) HackerNews {
	return HackerNews{config}
}

// Fetch to stories and store them in cache file
func (hn HackerNews) GetTopStories() ([]Story, error) {

	// Get top story ids
	res, err := http.Get(hn.config.ApiBaseUrl + "topstories.json")
	if err != nil {
		return []Story{}, err
	}

	defer res.Body.Close()
	jsonIds, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Story{}, err
	}

	var ids []int
	err = json.Unmarshal(jsonIds, &ids)
	if err != nil {
		return []Story{}, err
	}

	ids = ids[:hn.config.ItemsLimit]

	stories := make([]Story, len(ids))

	// TODO: Add serial mode when hn.config.ParallellMode is false

	sem := make(chan bool, len(ids))

	for i, id := range ids {
		go func(i, id int) {
			story, err := hn.GetStory(id)
			if err != nil {
				story = Story{i, id, "ERROR", "", err.Error(), "no url"}
			}
			story.I = i
			stories[i] = story
			sem <- true
		}(i, id)
	}

	for i := 0; i < len(ids); i++ {
		<-sem
	}

	err = hn.storeCacheFile(stories)
	if err != nil {
		return stories, err
	}

	return stories, nil
}

func (hn HackerNews) GetStoryByIndex(index int) (Story, error) {

	stories, err := hn.readCacheFile()
	if err != nil {
		// TODO: Return a Story containing error message?
		return Story{}, err
	}

	if index < 0 || index >= len(stories) {
		// TODO: Return a Story containing error message?
		return Story{}, errors.New("Index out of range")
	}

	return stories[index], nil
}

func (hn HackerNews) GetStory(id int) (Story, error) {

	url := hn.config.ApiBaseUrl + "item/" + strconv.Itoa(id) + ".json"

	res, err := http.Get(url)
	if err != nil {
		return Story{}, err
	}

	defer res.Body.Close()
	jsonItem, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Story{}, err
	}

	var item struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Url   string `json:"url"`
		// other attributs are ignored
	}

	err = json.Unmarshal(jsonItem, &item)
	if err != nil {
		return Story{}, err
	}

	story := Story{
		Id:    item.Id,
		Title: item.Title,
		Url:   item.Url,
	}

	return story, nil
}

func (hn HackerNews) storeCacheFile(stories []Story) error {

	bs, err := json.Marshal(stories)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hn.config.CacheFilePath, bs, 0644)
}

func (hn HackerNews) readCacheFile() ([]Story, error) {

	bs, err := ioutil.ReadFile(hn.config.CacheFilePath)
	if err != nil {
		return []Story{}, err
	}

	var stories []Story

	return stories, json.Unmarshal(bs, &stories)
}
