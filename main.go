package main

import (
	"github.com/gogf/gf/frame/g"
	_ "review-order/boot"
	_ "review-order/router"
)

func main() {
	g.Server().Run()
}
