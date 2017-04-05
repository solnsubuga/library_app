package api

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/snsubuga/library/models"
)

//Find a book from the API
func Find(id string) (models.ClassifyBookResponse, error) {
	var c models.ClassifyBookResponse
	body, err := ClassifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&owi=" + url.QueryEscape(id))
	if err != nil {
		return models.ClassifyBookResponse{}, err
	}
	err = xml.Unmarshal(body, &c)
	return c, err
}

//Search finds a book
func Search(query string) ([]models.SearchResult, error) {
	body, err := ClassifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query))
	if err != nil {
		return []models.SearchResult{}, err
	}

	var c models.ClassifySearchResponse
	err = xml.Unmarshal(body, &c)
	return c.Results, err
}

//ClassifyAPI api for books
func ClassifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(url); err != nil {
		log.Fatalf("%s", err.Error())
		return []byte{}, err
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return []byte{}, err
	}
	return body, nil
}
