package routers

import (
	"github.com/arxanchain/baas-poe-demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/upload", &controllers.UploadController{})
	beego.Router("/v1/up", &controllers.UploadController{}, "POST:Check")
}
