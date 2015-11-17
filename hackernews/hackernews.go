package hackernews

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// The Hacker News struct
type HackerNews struct {
	config Config
}

// Configuration, stored in ~/.hngorc
type Config struct {
	ApiBaseUrl        string
	ItemsLimit        int
	CacheFilePath     string
	OpenCommand       []string
	ShowCommandOutput bool
	// TODO: Add CacheTimeToLiveSecs int
}

// A Story
type Story struct {
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
	if json.Unmarshal(jsonIds, &ids) != nil {
		return []Story{}, err
	}

	// Limit the number of Story IDs
	ids = ids[:hn.config.ItemsLimit]
	stories := make([]Story, len(ids))

	sem := make(chan bool, len(ids))

	// Fetch all Stories in parallel
	for i, id := range ids {
		go func(i, id int) {
			story, err := hn.GetStory(id)
			if err != nil {
				story = Story{id, "ERROR", "", err.Error(), "no url"}
			}
			stories[i] = story
			sem <- true
		}(i, id)
	}

	// Wait for each go routine to signal
	for i := 0; i < len(ids); i++ {
		<-sem
	}

	// Store the cache and return stories
	return stories, hn.storeCacheFile(stories)
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
		Time  int    `json:"time"`
		// other attributs are ignored
	}

	err = json.Unmarshal(jsonItem, &item)
	if err != nil {
		return Story{}, err
	}

	t := time.Unix(int64(item.Time), 0)

	return Story{
		Id:    item.Id,
		Title: item.Title,
		Url:   item.Url,
		Date:  t.String(),
	}, nil
}

func (hn HackerNews) GetCachedStoryByIndex(index int) (Story, error) {

	stories, err := hn.readCacheFile()
	if err != nil {
		return Story{}, err
	}

	if index < 0 || index >= len(stories) {
		return Story{}, errors.New("Index out of range")
	}

	return stories[index], nil
}
