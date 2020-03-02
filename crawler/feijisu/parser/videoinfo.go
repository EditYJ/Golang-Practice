package parser

import (
	"demo/crawler/engine"
	"demo/crawler/model"
	"regexp"
)

type VideoInfo struct {
	typeRe         *regexp.Regexp
	typeContentRe  *regexp.Regexp
	yearRe         *regexp.Regexp
	regionRe       *regexp.Regexp
	nameRe         *regexp.Regexp
	updateStatueRe *regexp.Regexp
	descriptionRe  *regexp.Regexp
	imgUrlRe       *regexp.Regexp
	baiduDbankRe   *regexp.Regexp
}

var videoInfo = VideoInfo{
	typeRe:        regexp.MustCompile(`<p class="item"><span>类型：</span>([\s\S]*?)</p>`),
	typeContentRe: regexp.MustCompile(`<a.*?>(.*?)</a>`),

	yearRe:         regexp.MustCompile(`<span>年代：</span><a.*?>(.*?)</a></p>`),
	regionRe:       regexp.MustCompile(`<span>地区：</span><a.*?>(.*?)</a></p>`),
	nameRe:         regexp.MustCompile(`<h1>(.*)</h1>`),
	updateStatueRe: regexp.MustCompile(`<p class="tag">(.*)</p>`),
	descriptionRe:  regexp.MustCompile(`<span>简介：</span>([\s\S]*?)</p>`),
	imgUrlRe:       regexp.MustCompile(`<div class="b-detailcover tj-cover" ><img src="(.*?)".*?</div>`),
	baiduDbankRe:   regexp.MustCompile(`<a class="down" href="(.*?)" .*?百度网盘</a>`),
}

var notMainWebVideoInfo = VideoInfo{
	typeRe:        regexp.MustCompile(`<dd>[\s]*?<b>类型：</b>([\s\S]*?)</dd>`),
	typeContentRe: regexp.MustCompile(`<a.*?>(.+?)</a>`),

	yearRe:         regexp.MustCompile(`<b>年代：</b>([0-9]+?)</dd>`),
	regionRe:       regexp.MustCompile(`<b>地区：</b>(.*?)\s*?<b>`),
	nameRe:         regexp.MustCompile(`<dt class="name">(.+?)<`),
	updateStatueRe: regexp.MustCompile(`class="name".+?<span style="font-size:12px;margin-left:10px;">(.+?)</span>`),
	descriptionRe:  regexp.MustCompile(`<div class="des2"><b>剧情：</b>([\s\S]+?)</div>`),
	imgUrlRe:       regexp.MustCompile(`<div class="pic"><img src="(.+?)".*?</div>`),
	baiduDbankRe:   regexp.MustCompile(`<a title='百度云下载'  href='(.+?)' target="_self">`),
}

//videoIdRe      = regexp.MustCompile(`<link rel="alternate".*?href=".+/([0-9]+?)/" />`))

func ParseVideoInfo(contents []byte, extraParam map[string]string) engine.ParseResult {
	var videoRe VideoInfo
	pageType := extraParam["pageType"]
	result := engine.ParseResult{}
	video := model.Video{}

	if pageType == "0" {
		videoRe = videoInfo
	} else if pageType == "1" {
		videoRe = notMainWebVideoInfo
	}

	// 得到电影类型
	typeRes := filterString(contents, videoRe.typeRe)
	typeContentRes := videoRe.typeContentRe.FindAllSubmatch([]byte(typeRes), -1)
	for _, typeContent := range typeContentRes {
		video.Type = append(video.Type, string(typeContent[1]))
	}
	// 得到视频id
	video.Id = extraParam["id"]
	// 得到视频名
	video.Name = filterString(contents, videoRe.nameRe)
	// 得到视频更新状态
	video.UpdateStatus = filterString(contents, videoRe.updateStatueRe)
	// 得到视频年代
	video.Time = filterString(contents, videoRe.yearRe)
	// 得到视频地区
	video.Region = filterString(contents, videoRe.regionRe)
	// 得到视频描述
	video.Description = filterString(contents, videoRe.descriptionRe)
	// 得到视频简介图片
	video.ImgUrl = filterString(contents, videoRe.imgUrlRe)
	// 得到百度网盘爬取页面地址
	video.BaiduDbank = filterString(contents, videoRe.baiduDbankRe)

	result.Items = []interface{}{video}

	if video.BaiduDbank != "" {
		result.Requests = append(result.Requests, engine.Request{
			Url: video.BaiduDbank,
			ParserFunc: func(c []byte) engine.ParseResult {
				return BaiduDbank(c, map[string]string{"id": video.Id})
			},
		})
	}

	return result
}
