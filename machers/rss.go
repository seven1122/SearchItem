package machers

import (
	"SearchItem/search"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type (
	item struct {
		XMLName xml.Name `xml:"item"`
		PubDate string `xml:"pubDate"`
		Title string `xml:"title"`
		Description string `xml:"description"`
		Link string `xml:"link"`
		GUID string `xml:"guid"`
		GeoRssPoint string `xml:"georss:point"`
	}

	image struct {
		XMLName xml.Name `xml:"image"`
		URL string `xml:"url"`
		Title string `xml:"title"`
		Link string `xml:"link"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Title string `xml:"title"`
		Description string `xml:"description"`
		Link string `xml:"link"`
		PubDate string `xml:"pubDate"`
		LastBuildData string `xml:"lastBuildDate"`
		TTL string `xml:"ttl"`
		Language string `xml:"language"`
		ManagingEditor string `xml:"managingEditor"`
		WebMaster string `xml:"webMaster"`
		Image image `xml:"image"`
		Item []item `xml:"item"`
	}

	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel `xml:"channel"`
	}
)

type rssMatcher struct {

}

func init()  {
	var matcher rssMatcher
	search.Register("rss", matcher)

}

func (m rssMatcher) Search(feed *search.Feed, searchItem string) ([]*search.Result, error)  {
	var results []*search.Result
	log.Printf("search Feed Type[%s] for url [%s] site[%s]\n", feed.Type, feed.Uri, feed.Name)
	document, err := m.retrieve(feed)
	if err !=nil{
		return nil, err
	}

	for _, channelItem := range document.Channel.Item{
		matched, err := regexp.MatchString(searchItem, channelItem.Title)
		if err != nil {
			return nil, err
		}
		if matched{
			results = append(results, &search.Result{
				Field:"Title",
				Content: channelItem.Title,
			})
		}

		matched, err = regexp.MatchString(searchItem, channelItem.Description)
		if err != nil {
			return nil, err
		}
		if matched{
			results = append(results, &search.Result{
				Field:"Description",
				Content: channelItem.Description,
			})
		}
	}
	return results, nil
}

func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error)  {
	if feed.Uri == ""{
		return  nil, errors.New("No rss feed uri Provided")
	}

	resp , err := http.Get(feed.Uri)
	if err != nil{
		return nil, err
	}

	defer  resp.Body.Close()

	if  resp.StatusCode != 200{
		return nil, fmt.Errorf("http response error %d\n", resp.StatusCode)
	}

	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err

}