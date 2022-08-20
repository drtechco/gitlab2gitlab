package controllers

import (
	"github.com/kataras/iris/v12"
)

type WebController interface {
	Setup(r *iris.Application)
}

var apiControllers = [...]WebController{}

func InitApiController(r *iris.Application) {
	for _, controller := range apiControllers {
		controller.Setup(r)
	}
}
