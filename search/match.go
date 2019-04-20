package search

import "log"

type Result struct {
	Field string
	Content string
}

type Matcher interface {
	Search(feed *Feed, searchItem string) ([]*Result, error)
}

func Match(matcher Matcher, feed *Feed, searchItem string, results chan<- *Result)  {
	searchResults, err := matcher.Search(feed, searchItem)
	if err != nil{
		log.Println(err)
		return
	}
	for _, result := range searchResults{
		results <- result
	}

}

func Display(results chan *Result)  {
	for result := range results{
		log.Println(result.Field,":", result.Content)
	}
}
