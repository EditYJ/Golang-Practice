package datasave

import "log"

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	count := 0
	go func() {
		for {
			item := <-out
			log.Printf("Itemserver: got item #%d: %v", count, item)
			count++
		}
	}()
	return out
}
