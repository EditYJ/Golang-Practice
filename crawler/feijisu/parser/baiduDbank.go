package parser

import (
	"demo/crawler/engine"
	"demo/crawler/model"
	"demo/crawler/util"
	"regexp"
)

var sharePasswordRe = regexp.MustCompile(`<input id="biao1" class="share-password".*value="(.*?)"/>`)
var shareUrlRe = regexp.MustCompile(`id="paniframe" href="(.*?)"`)

func BaiduDbank(contents []byte, extraParam map[string]string) engine.ParseResult {
	result := engine.ParseResult{}
	dbank := model.Dbank{}

	// 得到视频id
	dbank.VideoId = extraParam["id"]
	// 得到百度网盘分享密码
	dbank.SharePwd = util.FilterString(contents, sharePasswordRe)
	// 得到百度网盘分享地址
	dbank.ShareUrl = util.FilterString(contents, shareUrlRe)

	result.Items = []interface{}{dbank}
	return result
}

