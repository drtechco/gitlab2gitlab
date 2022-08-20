package services

import (
	"drtech.co/gl2gl/models/results"
	jsoniter "github.com/json-iterator/go"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/modern-go/reflect2"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type _TimeValEncoder struct {
	Layout string
}

func (e _TimeValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	t := *((*time.Time)(ptr))
	return t.IsZero()
}
func (e _TimeValEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *((*time.Time)(ptr))
	stream.WriteString(t.Format(e.Layout))
}

type NullKeyVal struct {
	Key string
	Val interface{}
}
type XJsonConfig struct {
	Config     jsoniter.Config
	TimeLayout string
}
type _HttpContextService struct {
	logger      *logrus.Entry
	jsonAPI     map[string]jsoniter.API
	jsonAPILock sync.Mutex
}

func (c *_HttpContextService) GetLang(ctx *context2.Context) string {
	return sessions.Get(ctx).GetStringDefault("Lang", "en-us")
}

func (c *_HttpContextService) GetI18nService(context *context2.Context) *I18nService {
	return GetI18nService(c.GetLang(context))
}

func (c *_HttpContextService) GetErrorCodeService(context *context2.Context) *ErrorCodeService {
	return GetErrorCodeService(c.GetLang(context))
}

func (c *_HttpContextService) EMsg(ctx *context2.Context, code int) string {
	return c.GetErrorCodeService(ctx).Msg(code)
}

func (c *_HttpContextService) EJson(ctx *context2.Context, code int) {
	c.Json(ctx, results.ServerResult{
		Code:    code,
		Message: c.EMsg(ctx, code),
		Data:    nil,
	})

}
func (c *_HttpContextService) OkJson(ctx *context2.Context, data interface{}, opts ...NullKeyVal) {
	c.Json(ctx, results.ServerResult{
		Code:    0,
		Message: "",
		Data:    data,
	}, opts...)

}

/*
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJsonRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
	TimeLayout					  string
*/
func (c *_HttpContextService) Json(ctx *context2.Context, data interface{}, opts ...NullKeyVal) {
	json := c.getJsonAPI(opts)
	ctx.ContentType(context2.ContentJSONHeaderValue)
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Debugf("Marshal JSON: %v", err)
		ctx.StatusCode(http.StatusInternalServerError)
	}
	_, err = ctx.Write(jsonData)
	if err != nil {
		c.logger.Debugf("Write JSON: %v", err)
		ctx.StatusCode(http.StatusInternalServerError)
	}
}

func makeDefaultXConfig() XJsonConfig {
	return XJsonConfig{
		Config: jsoniter.Config{
			IndentionStep:                 0,
			MarshalFloatWith6Digits:       false,
			EscapeHTML:                    true,
			SortMapKeys:                   true,
			UseNumber:                     false,
			DisallowUnknownFields:         false,
			TagKey:                        "",
			OnlyTaggedField:               false,
			ValidateJsonRawMessage:        true,
			ObjectFieldMustBeSimpleString: false,
			CaseSensitive:                 false,
		},
		TimeLayout: "2006-01-02 15:04:05",
	}
}

func (c *_HttpContextService) Bool2Str(b bool) string {
	if b {
		return "1"
	} else {
		return "0"
	}
}

func (c *_HttpContextService) getJsonAPI(opts []NullKeyVal) jsoniter.API {
	if len(opts) > 0 {
		xconfig := makeDefaultXConfig()
		var sb strings.Builder
		for _, kv := range opts {
			sb.WriteString(kv.Key)
			if kv.Key == "IndentionStep" {
				xconfig.Config.IndentionStep = kv.Val.(int)
				sb.WriteString(strconv.Itoa(kv.Val.(int)))
			}
			if kv.Key == "MarshalFloatWith6Digits" {
				xconfig.Config.MarshalFloatWith6Digits = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "EscapeHTML" {
				xconfig.Config.EscapeHTML = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "SortMapKeys" {
				xconfig.Config.SortMapKeys = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "UseNumber" {
				xconfig.Config.UseNumber = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "DisallowUnknownFields" {
				xconfig.Config.DisallowUnknownFields = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "TagKey" {
				xconfig.Config.TagKey = kv.Val.(string)
				sb.WriteString(kv.Val.(string))
			}
			if kv.Key == "OnlyTaggedField" {
				xconfig.Config.OnlyTaggedField = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "ValidateJsonRawMessage" {
				xconfig.Config.ValidateJsonRawMessage = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "ObjectFieldMustBeSimpleString" {
				xconfig.Config.ObjectFieldMustBeSimpleString = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "CaseSensitive" {
				xconfig.Config.CaseSensitive = kv.Val.(bool)
				sb.WriteString(c.Bool2Str(kv.Val.(bool)))
			}
			if kv.Key == "TimeLayout" {
				xconfig.TimeLayout = kv.Val.(string)
				sb.WriteString(kv.Val.(string))
			}
		}
		sign := sb.String()
		c.jsonAPILock.Lock()
		defer c.jsonAPILock.Unlock()
		api, has := c.jsonAPI[sign]
		if !has {
			api = c.makeAJsonApi(xconfig)
			c.jsonAPI[sign] = api
		}
		return api
	}
	return c.jsonAPI["default"]
}

func (c *_HttpContextService) makeAJsonApi(xconfig XJsonConfig) jsoniter.API {
	api := xconfig.Config.Froze()
	exs := make(jsoniter.EncoderExtension)
	exs[reflect2.TypeOf(time.Time{})] = _TimeValEncoder{Layout: xconfig.TimeLayout}
	api.RegisterExtension(exs)
	return api
}

var HttpContextService = _HttpContextService{
	logger: logrus.WithField("Name", "HttpContextService"),
}

func SetupHttpContextService() {
	HttpContextService.jsonAPI = make(map[string]jsoniter.API)
	HttpContextService.jsonAPI["default"] = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	exs := make(jsoniter.EncoderExtension)
	exs[reflect2.TypeOf(time.Time{})] = _TimeValEncoder{Layout: "2006-01-02 15:04:05"}
	HttpContextService.jsonAPI["default"].RegisterExtension(exs)
}
