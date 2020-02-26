package parser

import (
	"demo/crawler/engine"
	"demo/crawler/model"
	"demo/crawler/util"
	"regexp"
)

var typeRe = regexp.MustCompile(`<p class="item"><span>类型：</span>([\s\S]*?)</p>`)
var typeContentRe = regexp.MustCompile(`<a.*?>(.*?)</a>`)

var yearRe = regexp.MustCompile(`<span>年代：</span><a.*?>(.*?)</a></p>`)
var regionRe = regexp.MustCompile(`<span>地区：</span><a.*?>(.*?)</a></p>`)
var nameRe = regexp.MustCompile(`<h1>(.*)</h1>`)
var updateStatueRe = regexp.MustCompile(`<p class="tag">(.*)</p>`)
var descriptionRe = regexp.MustCompile(`<span>简介：</span>([\s\S]*?)</p>`)
var imgUrlRe = regexp.MustCompile(`<div class="b-detailcover tj-cover" ><img src="(.*?)".*?</div>`)
var baiduDbankRe = regexp.MustCompile(`<a class="down" href="(.*?)" .*?百度网盘</a>`)
var videoIdRe = regexp.MustCompile(`<link rel="alternate".*?href=".+/([0-9]+?)/" />`)

func ParseVideoInfo(contents []byte) engine.ParseResult {
	result := engine.ParseResult{}
	video := model.Video{}

	// 得到电影类型
	typeRes := util.FilterString(contents, typeRe)
	typeContentRes := typeContentRe.FindAllSubmatch([]byte(typeRes), -1)
	for _, typeContent := range typeContentRes {
		video.Type = append(video.Type, string(typeContent[1]))
	}
	// 得到视频id
	video.Id = util.FilterString(contents, videoIdRe)
	// 得到视频名
	video.Name = util.FilterString(contents, nameRe)
	// 得到视频更新状态
	video.UpdateStatus = util.FilterString(contents, updateStatueRe)
	// 得到视频年代
	video.Time = util.FilterString(contents, yearRe)
	// 得到视频地区
	video.Region = util.FilterString(contents, regionRe)
	// 得到视频描述
	video.Description = util.FilterString(contents, descriptionRe)
	// 得到视频简介图片
	video.ImgUrl = util.FilterString(contents, imgUrlRe)
	// 得到百度网盘爬取页面地址
	video.BaiduDbank = util.FilterString(contents, baiduDbankRe)

	result.Items = []interface{}{video}
	result.Requests = append(result.Requests, engine.Request{
		Url: video.BaiduDbank,
		ParserFunc: func(c []byte) engine.ParseResult {
			return BaiduDbank(c, map[string]string{"id": video.Id})
		},
	})
	return result
}
