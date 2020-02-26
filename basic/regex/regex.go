package main

import (
	"fmt"
	"regexp"
)

const text = `
my email is y158392613@qq.com
my email is yujie@qq.com
my email is jinqianqian@qq.com
my email is jinqianqian@qq.com.cn
`

func main() {
	re := regexp.MustCompile(`([a-z0-9A-Z]+)@([a-z0-9A-Z]+)(\.[a-z0-9A-Z.]+)`)
	findString := re.FindAllStringSubmatch(text, -1)
	fmt.Println(findString)
}
