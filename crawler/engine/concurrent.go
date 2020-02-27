package engine

import (
	"log"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
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
		ce.Scheduler.Submit(r)
	}

	itemCount := 0
	for {
		// 接收结果
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got Item %d: %v", itemCount, item)
			itemCount++
		}

		// 提交任务给调度器
		for _, request := range result.Requests {
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
