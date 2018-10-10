package main

import (
	_ "github.com/arxanchain/jove/services/poe-demo/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
