package datasave

import (
	"context"
	"demo/crawler/model"
	"github.com/olivere/elastic/v7"
	"log"
)

func ItemSaver() (chan interface{}, error) {
	//client, err := elastic.NewClient(elastic.SetSniff(false))
	//if err != nil {
	//	return nil, err
	//}
	out := make(chan interface{})
	count := 0
	go func() {
		for {
			item := <-out
			if val, ok := item.(model.Video); ok {
				log.Printf("Itemserver: got item #%d: %+v", count, val)
				count++
			}
			// 存储
			//err := save(item, client)
			//if err != nil {
			//	log.Printf("Item server error saving item #%v: %v", item, err)
			//}
		}
	}()
	return out, nil
}

func save(item interface{}, client *elastic.Client) (err error) {
	// 如果是Video 存入Video表中
	if val, ok := item.(model.Video); ok {
		_, err := client.Index().Index("video_info").Id(val.Id).BodyJson(val).Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
