package hackernews

import (
	"encoding/json"
	"io/ioutil"
)

// Store a Story list in cache file
func (hn HackerNews) storeCacheFile(stories []Story) error {

	bs, err := json.Marshal(stories)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hn.config.CacheFilePath, bs, 0644)
}

// Retrieve Story list from cache file
func (hn HackerNews) readCacheFile() ([]Story, error) {

	bs, err := ioutil.ReadFile(hn.config.CacheFilePath)
	if err != nil {
		return []Story{}, err
	}

	var stories []Story

	return stories, json.Unmarshal(bs, &stories)
}
