package engine

import (
	"log"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan    chan interface{}
}

type Scheduler interface {
	Notifier
	Submit(Request)
	GetWorkChan() chan Request
	Run()
}

type Notifier interface {
	WorkerReady(chan Request)
}

func (ce *ConcurrentEngine) Run(seeds ...Request) {
	// 初始化
	out := make(chan ParseResult)
	ce.Scheduler.Run()
	for i := 0; i < ce.WorkerCount; i++ {
		createWorker(ce.Scheduler.GetWorkChan(), out, ce.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			log.Printf("%s，此Url已爬取过，自动略过...", r.Url)
			continue
		}
		ce.Scheduler.Submit(r)
	}

	for {
		// 接收结果
		result := <-out
		for _, item := range result.Items {
			go func(it interface{}) {
				ce.ItemChan <- it
			}(item)
			//if _, ok := item.(model.Video); ok {
			//	itemCount++
			//	log.Printf("Got Item #%d : %v", itemCount, item.(model.Video).Name)
			//}
		}

		// 提交任务给调度器
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				log.Printf("%s，此Url已爬取过，自动略过...", request.Url)
				continue
			}
			ce.Scheduler.Submit(request)
		}
	}
}

func createWorker(in chan Request, out chan ParseResult, n Notifier) {
	go func() {
		for {
			n.WorkerReady(in)
			request := <-in
			parseResult, err := worker(request)
			if err != nil {
				continue
			}
			out <- parseResult
		}
	}()
}

// url去重
var urls = make(map[string]bool)

func isDuplicate(url string) bool {
	if urls[url] {
		return true
	}
	urls[url] = true
	return false
}
