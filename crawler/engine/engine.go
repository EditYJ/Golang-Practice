package engine

import (
	"demo/crawler/fetcher"
	"log"
	"time"
)

func Run(seeds ...Request) {
	// 需要维护的任务队列
	var requests []Request

	requests = append(requests, seeds...)

	for len(requests) > 0 {
		// 1. 取出requests中第一个request进行处理
		request := requests[0]
		requests = requests[1:]

		log.Printf("正在爬取 %s ...", request.Url)

		// 得到源网页源码
		fetchRes, err := fetcher.Fetch(request.Url)
		if err != nil {
			log.Printf("请求网页 %s 失败，原因：%v", request.Url, err)
			continue
		}
		parseResult := request.ParserFunc(fetchRes)
		requests = append(requests, parseResult.Requests...)
		for _, item := range parseResult.Items {
			log.Printf("爬取结果：%v", item)
		}
		time.Sleep(time.Second)
	}
}
