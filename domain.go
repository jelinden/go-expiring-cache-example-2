package main

type Weather struct {
	Query struct {
		Results struct {
			Channel struct {
				Item struct {
					Condition struct {
						Temp string `json:"temp"`
						Text string `json:"text"`
					} `json:"condition"`
					Title string
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}

var weatherURLs = []string{
	"https://query.yahooapis.com/v1/public/yql?q=select%20item%20from%20weather.forecast%20where%20woeid=565346%20and%20u=\"c\"&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys",
	"https://query.yahooapis.com/v1/public/yql?q=select%20item%20from%20weather.forecast%20where%20woeid=573961%20and%20u=\"c\"&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys",
}
