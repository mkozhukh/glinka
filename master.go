package main

import (
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
	queued     map[string]bool
}

func (m *spiderMaster) run(start string) *LinksStore {
	m.queue = []string{}
	m.queued = make(map[string]bool)
	m.store = NewLinksStore()

	m.readPoint = 0
	m.writePoint = 0
	m.counter = 0

	threads := 3
	tasks := make(chan string, threads)
	results := make(chan *LinkRecord, threads)

	nest := make([]spiderWorker, threads)
	for i := range nest {
		go nest[i].start(tasks, results)
	}

	bar := pb.StartNew(1)
	m.addToQueue(Link{Global: start, Raw: start, Status: StatusOK})

	for {
		if m.isDone() {
			bar.Finish()
			return m.store
		}

		// scheduler new job
		if m.isQueue() && len(tasks) != cap(tasks) {
			tasks <- m.getFromQueue()
		}

		//get processed results
		select {
		case result := <-results:
			m.store.Records[result.URL] = *result
			if result.Status == StatusOK {
				for _, link := range result.Links {
					if link.Status == StatusOK {
						m.addToQueue(link)
					}
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

func (m *spiderMaster) addToQueue(url Link) {
	if m.queued[url.Global] {
		return
	}

	m.queue = append(m.queue, url.Global)
	m.writePoint++
	m.queued[url.Global] = true
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
