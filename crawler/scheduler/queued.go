package scheduler

import (
	"demo/crawler/engine"
)

type QueuedScheduler struct {
	requestChan chan engine.Request
	// 每个worker有自己的chan，然后我们从中去选择放入request
	workChan chan chan engine.Request
}

func (qs *QueuedScheduler) GetWorkChan() chan engine.Request {
	return make(chan engine.Request)
}

func (qs *QueuedScheduler) Submit(r engine.Request) {
	qs.requestChan <- r
}

// 外界告诉我们有worker已经准备好了，可以从外界接收request
func (qs *QueuedScheduler) WorkerReady(w chan engine.Request) {
	qs.workChan <- w
}

func (qs *QueuedScheduler) Run() {
	qs.requestChan = make(chan engine.Request)
	qs.workChan = make(chan chan engine.Request)
	go func() {
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			// 如果requestQ和workerQ都有东西在排队，那么就可以分配工作了
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case r := <-qs.requestChan:
				// request排队
				requestQ = append(requestQ, r)
			case w := <-qs.workChan:
				// worker排队
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				// 分配工作
				// 去除队列中的位置
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
