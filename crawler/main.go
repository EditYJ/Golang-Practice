package main

import (
	"demo/crawler/engine"
	"demo/crawler/feijisu/parser"
)

func main() {
	engine.Run(engine.Request{
		Url:        "http://www.feijisu5.com/acg/0/0/all/1.html",
		ParserFunc: parser.ParseVideoList,
	})
}
