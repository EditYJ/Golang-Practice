package engine

import (
	"demo/crawler/fetcher"
	"log"
	"time"
)

type SimpleEngine struct {
}

func (se SimpleEngine) Run(seeds ...Request) {
	// 需要维护的任务队列
	var requests []Request

	requests = append(requests, seeds...)

	for len(requests) > 0 {
		// 1. 取出requests中第一个request进行处理
		request := requests[0]
		requests = requests[1:]

		parseResult, err := worker(request)
		if err != nil {
			continue
		}
		requests = append(requests, parseResult.Requests...)
		for _, item := range parseResult.Items {
			log.Printf("爬取结果：%v", item)
		}
		time.Sleep(1500*time.Millisecond)
	}
}

func worker(r Request) (ParseResult, error) {

	//log.Printf("正在爬取 %s ...", r.Url)

	// 得到源网页源码
	fetchRes, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("请求网页 %s 失败，原因：%v", r.Url, err)
		return ParseResult{}, err
	}
	parseResult := r.ParserFunc(fetchRes)
	return parseResult, nil
}
