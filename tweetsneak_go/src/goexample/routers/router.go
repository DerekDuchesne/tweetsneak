package routers

import (
	"goexample/controllers"
	"github.com/astaxie/beegae"
)

func init() {
    beegae.Router("/", &controllers.MainController{}, "get:Index")
    beegae.Router("/search", &controllers.MainController{}, "get:Search")
}
