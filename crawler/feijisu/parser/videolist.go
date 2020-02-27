package parser

import (
	"bytes"
	"demo/crawler/engine"
	"regexp"
)

const BaseUrl = "http://www.feijisu5.com"

var (
	VideoListRe = regexp.MustCompile(`<a class="js-tongjic" href="(.*/[a-z]+/[0-9]+/)" title="(.*)" target="_blank">`)
	NextListRe  = regexp.MustCompile(`class="pages[\s\S]+href="(.*?)" class="a1">下一页</a> `)
)

func ParseVideoList(contents []byte) engine.ParseResult {
	VideoListRes := VideoListRe.FindAllSubmatch(contents, -1)
	NextListRes := filterString(contents, NextListRe)

	result := engine.ParseResult{}
	for _, item := range VideoListRes {
		result.Items = append(result.Items, string(item[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:        completeUrl(item[1]),
			ParserFunc: ParseVideoInfo,
		})
	}

	result.Requests = append(result.Requests, engine.Request{
		Url:        completeUrl([]byte(NextListRes)),
		ParserFunc: ParseVideoList,
	})

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

// 得到正则过滤的值
func filterString(content []byte, re *regexp.Regexp) string {
	res := re.FindSubmatch(content)
	if len(res) >= 2 {
		return string(res[1])
	}
	return ""
}
