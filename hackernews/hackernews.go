package hackernews

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HackerNews struct {
	baseUrl string
}

type Story struct {
	I       int
	Id      int
	Title   string
	Date    string
	Content string
	Url     string
}

func NewHackerNews(baseUrl string) (HackerNews, error) {
	return HackerNews{
		baseUrl: baseUrl,
	}, nil
}

func (hn HackerNews) GetTopStories(limit int) ([]Story, error) {

	// Get top story ids
	res, err := http.Get(hn.baseUrl + "topstories.json")
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

	ids = ids[:limit]

	stories := make([]Story, len(ids))

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

	err = storeCacheFile(ids)
	if err != nil {
		return stories, err
	}

	return stories, nil
}

func (hn HackerNews) GetStoryByIndex(index int) (Story, error) {

	ids, err := readCacheFile()
	if err != nil {
		return Story{}, err
	}

	if index < 0 || index >= len(ids) {
		return Story{}, errors.New("Index out of range")
	}

	return hn.GetStory(ids[index])
}

func (hn HackerNews) GetStory(id int) (Story, error) {

	url := hn.baseUrl + "item/" + strconv.Itoa(id) + ".json"

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
		//by          string `json:"by"`
		//descendants int    `json:"descendants"`
		Id int `json:"id"`
		//kids        []int  `json:"kids"`
		//score int    `json:"score"`
		//time  int    `json:"time"`
		Title string `json:"title"`
		//kind  string `json:"type"`
		Url string `json:"url"`
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

func storeCacheFile(ids []int) error {

	// TODO: Move cache file path to settings
	cacheFile := "/tmp/hncache"

	b, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cacheFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readCacheFile() ([]int, error) {

	// TODO: Define cache file in settings
	cacheFile := "/tmp/hncache"

	b, err := ioutil.ReadFile(cacheFile)

	var ids []int

	err = json.Unmarshal(b, &ids)

	return ids, err
}
