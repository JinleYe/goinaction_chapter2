package search

import (
	"encoding/json"
	"log"
	"os"
)

const dataFile = "data/data.json"

type Feed struct {
	Name string `json:"site"`

	URI  string `json:"link"`
	Type string `json:"type"`
}

func RetrieveFeeds() ([]*Feed, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err

}
