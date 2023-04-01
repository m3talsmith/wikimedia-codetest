package main

import (
	"strings"
)

type WikiRevision struct {
	Content string `json:"content"`
	Header  WikiRevisionHeader
}

type WikiPage struct {
	Title     string         `json:"title"`
	Revisions []WikiRevision `json:"revisions"`
}

type WikiQuery struct {
	Pages []WikiPage `json:"pages"`
}

type WikiResponse struct {
	Query WikiQuery `json:"query"`
}

type WikiRevisionHeader struct {
	ShortDescription string
}

func (revision *WikiRevision) NormalizeHeader() {
	if headerContent, _, ok := strings.Cut(revision.Content, "\n\n"); ok {
		revision.Header = WikiRevisionHeader{}
		parts := strings.SplitAfter(headerContent, `}}\n`)
		sdContent := parts[0]
		sdContent = strings.Trim(sdContent, `\n`)
		sd := strings.Split(sdContent, "|")[1]
		sd, _, _ = strings.Cut(sd, `}}`)
		revision.Header.ShortDescription = sd
	}
}
