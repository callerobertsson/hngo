// Package hackernews implements the HackerNews type used to fetch stories from HackerNews
package hackernews

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// HackerNews type
type HackerNews struct {
	config Config
}

// Config contains the settings for HackerNews type, stored in ~/.hngorc
type Config struct {
	APIBaseURL        string
	ItemsLimit        int
	CacheFilePath     string
	OpenCommand       []string
	ShowCommandOutput bool
}

// Story holds information about a HackerNews story
type Story struct {
	ID      int
	Title   string
	Date    string
	Content string
	URL     string
}

// New creates a new HackerNews
func New(config Config) HackerNews {
	return HackerNews{config}
}

// GetTopStories fetches top stories and store them in cache file
func (hn HackerNews) GetTopStories() ([]Story, error) {

	// Get top story ids
	res, err := http.Get(hn.config.APIBaseURL + "topstories.json")
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

// GetStory fetches a story with ID id from HackerNews
func (hn HackerNews) GetStory(id int) (Story, error) {

	url := hn.config.APIBaseURL + "item/" + strconv.Itoa(id) + ".json"

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
		ID    int    `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
		Time  int    `json:"time"`
		// other attributs are ignored
	}

	err = json.Unmarshal(jsonItem, &item)
	if err != nil {
		return Story{}, err
	}

	t := time.Unix(int64(item.Time), 0)

	return Story{
		ID:    item.ID,
		Title: item.Title,
		URL:   item.URL,
		Date:  t.String(),
	}, nil
}

// GetCachedStoryByIndex fetches a story from the cache
func (hn HackerNews) GetCachedStoryByIndex(index int) (Story, error) {

	stories, err := hn.readCacheFile()
	if err != nil {
		return Story{}, err
	}

	if index < 0 || index >= len(stories) {
		return Story{}, errors.New("index out of range")
	}

	return stories[index], nil
}
