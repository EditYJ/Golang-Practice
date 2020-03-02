package parser

import (
	"bytes"
	"demo/crawler/engine"
	"regexp"
	"strings"
)

const BaseUrl = "http://www.feijisu5.com"

var (
	VideoListRe = regexp.MustCompile(`<a class="js-tongjic" href="(.*/[a-z]+/[0-9]+/)" title="(.*)" target="_blank">`)
	NextListRe  = regexp.MustCompile(`class="pages[\s\S]+href="(.*?)" class="a1">下一页</a> `)

	IdRe = regexp.MustCompile(`.+/([0-9]+?)/`)
)

func ParseVideoList(contents []byte) engine.ParseResult {
	VideoListRes := VideoListRe.FindAllSubmatch(contents, -1)
	NextListRes := filterString(contents, NextListRe)

	result := engine.ParseResult{}
	for _, item := range VideoListRes {
		// 是否是主站的视频
		isMainWeb := strings.HasPrefix(completeUrl(item[1]), BaseUrl)
		id := filterString(item[1], IdRe)
		var pageType = "0"
		if isMainWeb {
			pageType = "0"
		} else {
			pageType = "1"
		}

		result.Items = append(result.Items, string(item[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url: completeUrl(item[1]),
			ParserFunc: func(c []byte) engine.ParseResult {
				return ParseVideoInfo(c, map[string]string{"id": id, "pageType": pageType})
			},
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
