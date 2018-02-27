package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/purell"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type spiderMaster struct {
	domain     string
	report     *spiderReport
	nest       []spiderWorker
	queue      []*url.URL
	writePoint uint
	readPoint  uint
	counter    uint
	queued     map[string]bool
}

func (m *spiderMaster) run(start string) error {
	m.nest = make([]spiderWorker, 3)
	m.queue = []*url.URL{}
	m.queued = make(map[string]bool)
	m.readPoint = 0
	m.writePoint = 0
	m.counter = 0
	m.report = newSpiderReport()

	tasks := make(chan *url.URL, 3)
	results := make(chan *spiderResult, 3)

	for i := range m.nest {
		go m.nest[i].start(tasks, results)
	}

	bar := pb.StartNew(10)
	m.addToQueue(start, nil)

	for {
		if m.isDone() {
			m.report.toString()
			return nil
		}

		// scheduler new job
		if m.isQueue() && len(tasks) != cap(tasks) {
			tasks <- m.getFromQueue()
		}

		//get processed results
		select {
		case result := <-results:
			if result.links != nil {
				for i := range result.links {
					m.addToQueue(result.links[i], result.parent)
				}
			}
			m.counter++

			bar.Total = int64(m.writePoint)
			bar.Increment()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (m *spiderMaster) relativeToGlobal(link string, parent *url.URL) *url.URL {
	info, err := url.Parse(link)

	if err != nil {
		m.report.addError("[E] Invalid URL, " + link)
	}

	if parent == nil {
		//top level link doesn't have a parent
		return info
	}

	if info.Host == "" {
		//relative link
		if strings.HasPrefix(info.Path, ".") {
			info.Path = path.Base(parent.Path) + "/" + info.Path
		} else {
			info.Host = parent.Host
			info.Scheme = parent.Scheme
		}
	}

	if info.Scheme == "" {
		info.Scheme = parent.Scheme
	}

	if info.Host != parent.Host {
		m.report.addInfo(fmt.Sprintf(
			"[I] External link, from %s to %s",
			parent.String(),
			link))
		return nil
	}

	if info.Scheme != parent.Scheme {
		m.report.addError(fmt.Sprintf(
			"[E] Scheme changed, from %s to %s",
			parent.String(),
			link))
		return nil
	}

	return info
}

func (m *spiderMaster) addToQueue(data string, parent *url.URL) {
	url := m.relativeToGlobal(data, parent)
	if url == nil {
		return
	}

	normalized := purell.MustNormalizeURLString(url.String(), purell.FlagsUnsafeGreedy)

	if m.queued[normalized] {
		return
	}

	m.queue = append(m.queue, url)
	m.writePoint++
	m.queued[normalized] = true
	//log.Printf("%d %s", m.writePoint, normalized)
}

func (m *spiderMaster) getFromQueue() *url.URL {
	if m.readPoint == m.writePoint {
		return nil
	}
	url := m.queue[m.readPoint]
	m.readPoint++
	return url
}

func (m *spiderMaster) isDone() bool {
	return m.writePoint == m.counter
}

func (m *spiderMaster) isQueue() bool {
	return m.readPoint != m.writePoint
}
