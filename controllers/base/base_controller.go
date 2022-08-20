package base

import (
	"drtech.co/gl2gl/services"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
)

type BaseController struct {
}

func (c *BaseController) GetLang(ctx *context2.Context) string {
	return services.HttpContextService.GetLang(ctx)
}
func (c *BaseController) HasLogin(ctx *context2.Context) bool {
	return c.GetAdminId(ctx) > 0
}

func (c *BaseController) GetAdminId(ctx *context2.Context) int64 {
	//sessions.Get(ctx).
	return sessions.Get(ctx).GetInt64Default("adminId", 0)
	//return  1
}
func (c *BaseController) GetI18nService(context *context2.Context) *services.I18nService {
	return services.GetI18nService(c.GetLang(context))
}

func (c *BaseController) GetErrorCodeService(context *context2.Context) *services.ErrorCodeService {
	return services.GetErrorCodeService(c.GetLang(context))
}

func (c *BaseController) EMsg(ctx *context2.Context, code int) string {
	return c.GetErrorCodeService(ctx).Msg(code)
}

func (c *BaseController) EJson(ctx *context2.Context, code int) {
	services.HttpContextService.EJson(ctx, code)
}
func (c *BaseController) OkJson(ctx *context2.Context, data interface{}, opts ...services.NullKeyVal) {
	services.HttpContextService.OkJson(ctx, data, opts...)
}