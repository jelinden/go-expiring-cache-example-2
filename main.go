package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client http.Client

func init() {
	client = http.Client{
		Timeout: time.Duration(4 * time.Second),
	}
}

func main() {
	weatherToCache()
	go weatherSchedule()

	// Caching is handled with background scheduling.
	http.HandleFunc("/weather/helsinki", func(w http.ResponseWriter, r *http.Request) {
		cacheItem := GetItemFromCache(weatherURLs[0])
		if cacheItem != nil {
			w.Write(cacheItem.Value)
		} else {
			json.NewEncoder(w).Encode("not found")
		}
	})

	// If item is not found in cache, it's added to cache.
	// When the item is not in cache, the request is slower.
	http.HandleFunc("/weather/tikkurila", func(w http.ResponseWriter, r *http.Request) {
		cacheItem := GetItemFromCache(weatherURLs[1])
		if cacheItem != nil {
			w.Write(cacheItem.Value)
		} else {
			if weather := fetchWeather(weatherURLs[1]); weather != nil {
				j, err := json.Marshal(&weather.Query.Results.Channel.Item)
				if err != nil {
					fmt.Println(err)
				} else {
					AddItemToCache(weatherURLs[1], j, time.Minute)
				}
				w.Write(j)
			} else {
				json.NewEncoder(w).Encode("not found")
			}
		}
	})
	http.ListenAndServe(":8080", nil)
}

func weatherSchedule() {
	for range time.Tick(time.Minute * 1) {
		weatherToCache()
	}
}

func weatherToCache() {
	fmt.Println("fetching weather")
	if w := fetchWeather(weatherURLs[0]); w != nil {
		j, err := json.Marshal(&w.Query.Results.Channel.Item)
		if err != nil {
			fmt.Println(err)
		} else {
			AddItemToCache(weatherURLs[0], j, time.Minute)
		}
	}
}

func fetchWeather(url string) *Weather {
	var w Weather
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err = json.Unmarshal(buf, &w); err != nil {
		fmt.Println(err)
	}
	return &w
}
