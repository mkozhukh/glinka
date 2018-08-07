package main

import (
	"fmt"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

type spiderMaster struct {
	domain string
	store  *LinksStore

	queue      []string
	writePoint uint
	readPoint  uint
	counter    uint

	verbose bool
	quiet   bool
	threads int
}

func (m *spiderMaster) run(start string) *LinksStore {
	m.queue = []string{}
	m.store = NewLinksStore()

	m.readPoint = 0
	m.writePoint = 0
	m.counter = 0

	tasks := make(chan string, m.threads)
	results := make(chan *workerResult, m.threads)

	nest := make([]spiderWorker, m.threads)
	for i := range nest {
		go nest[i].start(tasks, results)
	}

	var bar *pb.ProgressBar
	if !m.verbose && !m.quiet {
		bar = pb.StartNew(1)
	}
	if m.verbose {
		fmt.Printf("Starting %d threads, from %s\n", m.threads, start)
	}
	m.addToQueue(Link{Global: start, Raw: start, Status: StatusOK})

	for {
		if m.isDone() {
			if !m.verbose && !m.quiet {
				bar.Finish()
			}
			return m.store
		}

		// scheduler new job
		if m.isQueue() && len(tasks) != cap(tasks) {
			url := m.getFromQueue()
			if m.verbose {
				fmt.Print("[start] " + url + "\n")
			}
			tasks <- url
		}

		//get processed results
		select {
		case result := <-results:
			if m.verbose {
				fmt.Printf("[end] %s, status: %d, %s\n", result.URL, result.Status, result.Error)
			}
			rec := m.store.Records[result.URL]
			rec.Status = result.Status
			rec.Error = result.Error
			if result.Status == StatusOK && len(result.Links) != 0 {
				rec.Links = make([]string, len(result.Links))
				for i := range result.Links {
					m.addToQueue(result.Links[i])
					rec.Links[i] = result.Links[i].Global
				}
			}

			m.counter++
			if !m.verbose && !m.quiet {
				bar.Total = int64(m.writePoint)
				bar.Increment()
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (m *spiderMaster) addToQueue(url Link) {
	rec, ok := m.store.Records[url.Global]
	if ok {
		rec.Count++
		return
	}

	newrec := LinkRecord{URL: url.Global, Count: 1}
	if url.Status == StatusOK {
		newrec.Status = StatusUnknown
		m.queue = append(m.queue, url.Global)
		m.writePoint++
	} else {
		newrec.Status = url.Status
	}

	if m.verbose {
		fmt.Printf("+ [%d] %s, %s\n", url.Status, url.Global, url.Raw)
	}
	m.store.Records[url.Global] = &newrec
}

func (m *spiderMaster) getFromQueue() string {
	if m.readPoint == m.writePoint {
		return ""
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
