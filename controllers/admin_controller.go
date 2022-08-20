package controllers

import (
	"drtech.co/gl2gl/controllers/base"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
)

type AdminController struct {
	base.BaseController
	logger *logrus.Entry
}

func (c *AdminController) Setup(app *iris.Application) {
	c.logger = logrus.WithField("Name", "AdminController")
	app.Post("/admin/login", c.Login)
}

func (c *AdminController) Login(context *context2.Context) {

}
