package parser

import (
	"bytes"
	"demo/crawler/engine"
	"regexp"
)

const BaseUrl = "http://www.feijisu5.com"

const VideoListRe = `<a class="js-tongjic" href="(.*/[a-z]+/[0-9]+/)" title="(.*)" target="_blank">`

func ParseVideoList(contents []byte) engine.ParseResult {
	compile := regexp.MustCompile(VideoListRe)
	findRes := compile.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, item := range findRes {
		result.Items = append(result.Items, string(item[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:        completeUrl(item[1]),
			ParserFunc: ParseVideoInfo,
		})
	}
	return result
}

func completeUrl(url []byte) string {
	var res string
	if index := bytes.IndexAny(url, "/"); index == 0 {
		res = BaseUrl + string(url)
	} else {
		res = string(url)
	}
	return res
}
