package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type spiderResult struct {
	links  []string
	parent *url.URL
	error  error
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
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (w *spiderWorker) get(url string, res *spiderResult) []string {
	head, err := http.Head(url)
	if err != nil {
		res.error = err
		return nil
	}

	str := head.Header.Get("Content-Type")
	if !strings.HasPrefix(str, "text/html") && !strings.HasPrefix(str, "text/plain") {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		res.error = err
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
						if !strings.HasPrefix(a.Val, "mailto:") {
							links = append(links, a.Val)
						}
					}
				}
			}
		}
	}
}
