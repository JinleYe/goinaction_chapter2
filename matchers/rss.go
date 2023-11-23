package matchers

import (
	"encoding/xml"
	"errors"
	"goinaction/sample1/search"
	"io"
	"log"
	"net/http"
	"regexp"
)

type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`

		Link        string `xml:"link"`
		GUID        string `xml:"guid"`
		GeoRssPoint string `xml:"georss:point"`
	}

	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}

	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

type rssMatcher struct{}

func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	var results []*search.Result
	log.Printf("Search Feed Type[%s] Site[%s] For URI[%s]\n\n", feed.Type, feed.Name, feed.URI)
	document, err := m.retrieve(feed)
	if err != nil {
		return nil, err
	}

	if document == nil {
		log.Printf("Failed to retrieve document from Feed Type[%s] Site[%s] For URI[%s]", feed.Type, feed.Name, feed.URI)
		return results, err
	}

	for _, channelItem := range document.Channel.Item {
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, &search.Result{
				Field:   "Title",
				Content: channelItem.Title,
			})
		}

		matchedDesc, errDesc := regexp.MatchString(searchTerm, channelItem.Description)
		if errDesc != nil {
			return nil, errDesc
		}
		if matchedDesc {
			// use & to get the address of this new value, which is stored in the slice
			results = append(results, &search.Result{
				Field:   "Description",
				Content: channelItem.Description,
			})
		}

	}
	return results, nil
}

func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	if feed.URI == "" {
		return nil, errors.New("no rss feed URI provided")
	}

	response, err := http.Get(feed.URI)
	if err != nil {
		log.Fatalf("Can't get the response from URI: %s", feed.URI)
		//var document rssDocument

		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error occurs while reading the response body, error: %s", err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		log.Printf("URI[%s] responded with status code %s", feed.URI, response.Status)
		return nil, err
	}

	var document rssDocument
	err = xml.NewDecoder(response.Body).Decode(&document)

	return &document, err

}

func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}
