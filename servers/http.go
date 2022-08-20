package servers

import (
	"drtech.co/gl2gl/controllers"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/sessions"
	"io/ioutil"
	"time"
)

func HttpRun(address string) (error, *iris.Application) {
	app := iris.New()
	app.SetRoutesNoLog(false)
	app.Logger().SetLevel("debug")

	sess := sessions.New(sessions.Config{
		Cookie:       "_session_id",
		AllowReclaim: true,
		Expires:      24 * time.Hour, // <=0 意味永久的存活
	})
	//sess.UseDatabase(db)
	controllers.InitApiController(app)
	app.UseRouter(sess.Handler())
	ac := accesslog.New(ioutil.Discard)
	ac.IP = true
	app.UseRouter(ac.Handler)
	go func() {
		err := app.Listen(address)
		if err != nil {
			return
		}
	}()
	return nil, app
}
