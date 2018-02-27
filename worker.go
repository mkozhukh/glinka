package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type spiderResult struct {
	links  []string
	parent *url.URL
}

type spiderWorker struct {
}

func (w *spiderWorker) start(tasks chan *url.URL, result chan *spiderResult) {
	for {
		select {
		case url := <-tasks:
			res := spiderResult{}
			res.parent = url
			res.links = w.get(url.String(), &res)
			result <- &res
		}
	}
}

func (w *spiderWorker) get(url string, res *spiderResult) []string {
	head, err := http.Head(url)
	if err != nil {
		log.Print(err.Error())
		return nil
	}

	str := head.Header.Get("Content-Type")
	if !strings.HasPrefix(str, "text/html") && !strings.HasPrefix(str, "text/plain") {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err.Error())
		return nil
	}

	return w.parseHTML(resp, res)
}

func (w *spiderWorker) parseHTML(resp *http.Response, res *spiderResult) []string {
	links := []string{}
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						links = append(links, a.Val)
					}
				}
			}
		}
	}
}
