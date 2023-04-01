package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const WIKI_API_URL = "https://en.wikipedia.org/w/api.php"

type result struct {
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
}

type response struct {
	Result  result `json:"result"`
	Message string `json:"message"`
	Status  Status `json:"status"`
}

type Status string

const (
	Error   Status = "Error"
	Success Status = "Success"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		response := response{}

		if err := r.ParseForm(); err != nil {
			response.Status = Error
			response.Message = "error getting query parameters"
			rjson, _ := json.Marshal(response)
			fmt.Fprint(w, string(rjson))
			return
		}
		name := r.FormValue("name")
		if name == "" {
			response.Status = Error
			response.Message = "missing name for query"
			rjson, _ := json.Marshal(response)
			fmt.Fprint(w, string(rjson))
			return
		}
		revision, err := fetchLatestRevision(name)
		if err != nil {
			response.Status = Error
			response.Message = err.Error()
			rjson, _ := json.Marshal(response)
			fmt.Fprint(w, string(rjson))
			return
		}
		result := result{Title: revision.Title, ShortDescription: revision.ShortDescription}
		response.Result = result
		response.Message = "result found"
		rjson, _ := json.Marshal(response)
		fmt.Fprint(w, string(rjson))
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func fetchLatestRevision(titles string) (*result, error) {
	// API auto capitalizes the first part of a titles, but no other portions of the titles on the call
	// The experience of my result would be nicer if I provided multiple results of variants of the titles sent
	titles = strings.Title(titles) //Deprecated, does not support full unicode word boundries
	titles = url.QueryEscape(titles)

	reqquery := fmt.Sprintf(
		"action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content",
		titles,
	)
	requrl := fmt.Sprintf("%s?%s", WIKI_API_URL, reqquery)
	resp, err := http.Get(requrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var wikiResult WikiResponse
	if err := json.Unmarshal(body, &wikiResult); err != nil {
		return nil, err
	}

	if len(wikiResult.Query.Pages) <= 0 {
		return nil, fmt.Errorf("no result found")
	}
	page := wikiResult.Query.Pages[0]
	if len(page.Revisions) <= 0 {
		return nil, fmt.Errorf("no result found")
	}
	revision := page.Revisions[0]
	revision.NormalizeHeader()

	return &result{Title: page.Title, ShortDescription: revision.Header.ShortDescription}, nil
}
