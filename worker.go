package main

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/purell"
	"golang.org/x/net/html"
)

type spiderWorker struct {
}

func (w *spiderWorker) start(tasks chan string, result chan *LinkRecord) {
	for {
		select {
		case url := <-tasks:
			links, errorText, status := w.get(url)
			res := LinkRecord{
				URL:    url,
				Links:  links,
				Status: status,
				Error:  errorText,
			}
			result <- &res
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (w *spiderWorker) get(urlString string) ([]Link, string, LinkStatus) {
	head, err := http.Head(urlString)
	if err != nil {
		return nil, err.Error(), StatusError
	}

	str := head.Header.Get("Content-Type")
	if !strings.HasPrefix(str, "text/html") && !strings.HasPrefix(str, "text/plain") {
		return nil, "", StatusBinary
	}

	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err.Error(), StatusError
	}

	parent, err := url.Parse(urlString)
	if err != nil {
		return nil, err.Error(), StatusError
	}

	return w.parseHTML(resp, parent), "", StatusOK
}

func (w *spiderWorker) parseHTML(resp *http.Response, parent *url.URL) []Link {
	links := []Link{}
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
							links = append(links, w.parseLink(a.Val, parent))
						}
					}
				}
			}
		}
	}
}

func (w *spiderWorker) parseLink(raw string, parent *url.URL) Link {
	info, err := url.Parse(raw)
	link := Link{
		Raw:    raw,
		Status: StatusOK,
	}

	if err != nil {
		link.Status = StatusError
		return link
	}

	if parent != nil {
		if info.Host == "" {
			//relative link
			if strings.HasPrefix(info.Path, ".") {
				info.Path = path.Base(parent.Path) + "/" + info.Path
			}
			info.Host = parent.Host
			info.Scheme = parent.Scheme
		}

		if info.Scheme == "" {
			info.Scheme = parent.Scheme
		}

		if info.Host != parent.Host {
			link.Status = StatusExternal
		} else if info.Scheme != parent.Scheme {
			link.Status = StatusMixedContent
		}
	}
	// if strings.HasSuffix(info.Path, " ") {
	// 	m.report.addWarning(fmt.Sprintf(
	// 		"Whitespace after link, %s at %s",
	// 		info,
	// 		parent,
	// 	))
	// }

	link.Global = purell.MustNormalizeURLString(info.String(), purell.FlagsUsuallySafeNonGreedy|purell.FlagRemoveDirectoryIndex|purell.FlagRemoveFragment|purell.FlagRemoveDuplicateSlashes|purell.FlagSortQuery)
	return link
}
