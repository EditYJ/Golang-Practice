package datasave

import (
	"context"
	"demo/crawler/model"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestSave(t *testing.T) {
	videoTestData := model.Video{
		Id:           "528",
		Name:         "黑白来看守所第一季 UpdateStatus:更新至13集",
		UpdateStatus: "更新至13集",
		Type:         []string{"搞笑", "热血"},
		Time:         "2016",
		Region:       "日本",
		Description:  "黑白来看守所动画主要讲述了在日\n本最大的监狱中，4个碎碎念囚犯不断逼疯狱卒的黑白来爆笑喜剧。到底是什么样子的精彩故事让人爆笑不止呢，狱卒又是怎样被这4名囚犯给逼疯了呢。",
		ImgUrl:       "http://ae01.alicdn.com/kf/Hd85499c9ec454d\na29f685c1581d018a5R.jpg",
		BaiduDbank:   "http://m.feijisu5.com/down/528.html",
	}
	id, err := save(videoTestData)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	result, err := client.Get().Index("video_info").Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//t.Logf("%v", string(result.Source))
	var video model.Video
	err = json.Unmarshal(result.Source, &video)
	if err != nil {
		panic(err)
	}

	t.Logf("Got %+v", video)
}
