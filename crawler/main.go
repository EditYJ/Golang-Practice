package main

import (
	"demo/crawler/datasave"
	"demo/crawler/engine"
	"demo/crawler/feijisu/parser"
	"demo/crawler/scheduler"
)

func main() {
	//simpleEngine := engine.SimpleEngine{}
	itemChan, err := datasave.ItemSaver()
	if err != nil {
		panic(err)
	}
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}

	concurrentEngine.Run(engine.Request{
		Url:        "http://www.feijisu5.com/acg/0/0/all/1.html",
		ParserFunc: parser.ParseVideoList,
	})
}
