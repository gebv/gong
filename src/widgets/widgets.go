package widgets

import (
	"github.com/shurcooL/github_flavored_markdown"
	"html/template"
	// texttemplate "text/template"
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"store"
)

const (
	// CountWidgetPerPage count widgets per page
	CountWidgetPerPage = 250
)

const (
	// ROUTES
	ROUTES = "routes"

	// LAYOUTS
	LAYOUTS = "layouts"
)

// Helpfull functions

// extractPropertyByName вспомогательная функция получения значения из props
func extractPropertyByName(file *store.File, propKey string) (string, error) {
	// TODO: Вынести в file
	src, isExist := file.Props[propKey]

	if _, isValid := src.(string); !isExist && !isValid {
		return "", fmt.Errorf("not sets or not valid key '%v' from widget ext_id='%v' categories=%v", propKey, file.ExtId, file.Collections)
	}

	return src.(string), nil
}

func NewConfig() Config {
	return Config(map[string]interface{}{})
}

type Config map[string]interface{}

func (c Config) Get(name string) interface{} {
	return c[name]
}

func (c Config) GetAsBool(name string) bool {
	_value, ok := c[name].(bool)
	return ok && _value
}

func (c Config) GetAsString(name string) string {
	_value, ok := c[name].(string)
	if !ok {
		return ""
	}
	return _value
}

func (c *Config) SetNewConfig(v interface{}) {
	if _v, ok := v.(map[string]interface{}); ok {
		*c = _v
	} else {
		*c = NewConfig()
	}
}

func NewContext(c echo.Context) *Context {
	model := new(Context)

	model.Global = make(map[string]interface{})
	model.Self = make(map[string]interface{})

	model.Context = c
	// model.Request = r

	model.Funcs = template.FuncMap{
		"log": func(args ...interface{}) template.HTML {
			glog.Infof("[DEV]: %v", args)
			return template.HTML("")
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"md": func(v interface{}) template.HTML {
			var str string
			switch v.(type) {
			case string, template.HTML:
				str = v.(string)
			default:
				str = fmt.Sprintf("%v", v)
			}
			return template.HTML(github_flavored_markdown.Markdown([]byte(str)))
		},
		"setter": func(m map[string]interface{}, key string, v interface{}) template.HTML {
			m[key] = v
			return ""
		},
		"getter": func(m map[string]interface{}, key string) interface{} {
			return m[key]
		},
		"get_widget": func(group, key string) interface{} {
			item, err := store.NewOrLoadFileOfBucket(group, key)

			if err != nil {
				glog.Warningf("template funcs: get_execute file: group=%v, key=%v, err=%v", group, key, err)
			}

			return item
		},
		"redirect": func(path string) template.HTML {
			model.Context.Redirect(302, path)
			model.Abort()
			return template.HTML("")
		},
		"search": func(group, query string) interface{} {
			filter := store.NewSearchFileter()
			filter.AddCollections(group)
			filter.SetQuery(query)

			items := store.SearchPerPage(filter)

			return items
		},
		"widget": func(ctx *Context, args ...interface{}) template.HTML {

			context := NewContext(ctx.Context)
			context.Global = ctx.Global // копируется ссылка значений

			if len(args) == 0 || len(args) < 2 {
				glog.Warningf("template funcs: min 2 args %v\n", args)

				return ""
			}

			_group, ok1 := args[0].(string)
			_key, ok2 := args[1].(string)

			if !ok1 || !ok2 {
				glog.Warningf("template funcs: execute file: not valid widget group name or key %v\n", args)
				return ""
			}

			buff := bytes.NewBufferString("")

			context.Self["_group"] = _group
			context.Self["_name"] = _key

			if err := context.RenderWidget(buff, _group, _key); err != nil {
				glog.Warningf("template funcs: execute file: group=%v, kye=%v, err=%v", _group, _key, err)
			}

			// TODO: render widget
			return template.HTML(buff.String())
		},
	}

	if err := model.initAppSettings(); err != nil {

		glog.Errorf("context: init app settings err=%v", err)
	}

	return model
}

type Context struct {
	Self   map[string]interface{} `toml:"self"`
	Global map[string]interface{} `toml:"global"`

	Funcs template.FuncMap

	Context echo.Context

	TraceWidgets []string
	depth        int

	isAbort bool

	Routing configRouting `toml:"routing"`
	Theme   configTheme   `toml:"theme"`
	Pages   configPage    `toml:"pages"`

	Store   interface{} // Постоянное хранилище
	Cookie  interface{} //
	Session interface{} //
}

// initAppSettings инициализация настроек приложения
func (c *Context) initAppSettings() error {
	// item, err := store.GetByPath(AppSettingsBucketName, AppSettingsFileName)
	item, err := store.FindFile(AppSettingsBucketName, AppSettingsFileName)

	if err != nil {
		glog.Warningf("app settings: error get file name='%v.%v' err=%v", AppSettingsBucketName, AppSettingsFileName, err)
		return err
	}

	configRaw, err := extractPropertyByName(item, "Config")

	if err != nil {
		glog.Warningf("app settings: error get config name='%v.%v' err=%v", AppSettingsBucketName, AppSettingsFileName, err)
		return err
	}

	compiledConfigBuff := bytes.NewBufferString("")
	err = c.render(compiledConfigBuff, "route_config", configRaw, c)

	if err != nil {
		return err
	}

	if _, err := toml.Decode(compiledConfigBuff.String(), c); err != nil {
		glog.Warningf("app settings: config decode type=%v err=%v", "toml", err)

		return err
	}

	return nil
}

func (c *Context) render(writer *bytes.Buffer, key, raw string, data interface{}) error {
	tpl, err := template.New(key).Funcs(template.FuncMap(c.Funcs)).Parse(raw)

	if err != nil {
		glog.Warningf("render: parse template err=%v", err)

		return err
	}

	tt := template.Must(tpl, err)
	if err := tt.Execute(writer, c); err != nil {
		glog.Warningf("render: compile template err=%v", err)

		return err
	}

	return nil
}

// executeSimpleMode выполнить в простом режиме
func (c *Context) executeSimpleMode(writer io.Writer, urlPath string) error {
	// в качестве имени файла используется path

	var bucketName = c.DefaultRouteBucket()

	if err := c.traceWidgets(bucketName, urlPath); err != nil {
		return err
	}

	c.Context.Response().Header().Set("Content-Type", c.Pages.DefaultContentType)

	return c.RenderWidget(writer, bucketName, urlPath)
}

// Execute получает file с конфигурацией и выполяет действия
func (c *Context) Execute(writer io.Writer, urlPath string) error {

	if c.IsSimpleModeRouting() {
		err := c.executeSimpleMode(writer, urlPath)

		if err == store.ErrNotFound {
			_file := BFI(c.Pages.PageNotFound)

			glog.Infof("execute: error page 404 value=%v, is_valid=%v", _file, _file.IsValid())
			c.Context.Response().WriteHeader(http.StatusNotFound)

			if !_file.IsValid() {
				return err
			}

			return c.RenderWidget(writer, _file.Bucket(), _file.File())
		}

		return err
	}

	// TODO: Роутинг описан в настройках
	// TODO: Content-type должен быть описан в самом файле
	// TODO: В сложном случае, из списка buckets находим первый файл и исполняем его

	// config := newRouteConfig()

	// if _, err := toml.Decode(compiledConfigBuff.String(), config); err != nil {
	//     glog.Warningf("routing: config decode type=%v err=%v", "toml", err)

	//     return err
	// }

	// // Проверяем корректность роутинга
	// if len(config.GoTo()) != 2 {
	//     return fmt.Errorf("routing: not valid route %v", config.To)
	// }

	// if _, err := toml.Decode(compiledConfigBuff.String(), c); err != nil {
	//     glog.Warningf("routing: config decode type=%v err=%v", "toml", err)

	//     return err
	// }

	// c.Global["method"] = c.Context.Request().Method()
	// c.Global["path"] = c.Context.Request().URL().Path()
	// c.Global["GetQueryParam"] = c.Context.Request().URL().QueryParam
	// c.Global["GetFormValue"] = c.Context.Request().FormValue

	// // ACL
	// allowed := false
	// currentMethodName := strings.ToLower(c.Context.Request().Method())

	// for _, methodName := range config.Methods {
	//     if currentMethodName == strings.ToLower(methodName) {
	//         allowed = true
	//     }
	// }

	// if !allowed {
	//     glog.Warningf("routing: not allowed method=%v, allowed methods=%v", c.Context.Request().Method(), config.Methods)

	//     return fmt.Errorf("routing: not allowed method")
	// }

	// for name, value := range config.Headers {
	//     c.Context.Response().Header().Add(name, value)
	// }

	// //
	// // CONTENT
	// //

	// contentRaw, err := extractPropertyByName(_file, "Content")

	// if err != nil {
	//     glog.Warningf("execute file: get CONTENT group=%v name=%v err=%v", ROUTES, urlPath, err)
	//     return err
	// }

	// compiledContentBuff := bytes.NewBufferString("")
	// err = c.render(compiledContentBuff, "routing_content", contentRaw, c)

	// if err != nil {
	//     return err
	// }

	return fmt.Errorf("routing: не реализован режим %v", c.Routing.Mode)

}

func (c *Context) RenderWidget(writer io.Writer, group, name string) error {
	if c.isAbort {
		return fmt.Errorf("execute file: abort")
	}

	if err := c.traceWidgets(group, name); err != nil {

		return err
	}

	widget, err := store.FindFile(group, name)

	if err != nil {
		glog.Warningf("execute file: error get widget group=%v name=%v err=%v", group, name, err)

		return err
	}

	c.Self["_group"] = group
	c.Self["_name"] = name

	//
	// CONFIG
	//

	configRaw, err := extractPropertyByName(widget, "Config")

	if err != nil {
		glog.Warningf("execute file: get CONFIG group=%v name=%v err=%v", group, name, err)
		return err
	}

	compiledConfigBuff := bytes.NewBufferString("")
	err = c.render(compiledConfigBuff, "widget_config", configRaw, c)

	if err != nil {
		return err
	}

	// TODO: Разные типы конфигов

	if _, err := toml.Decode(compiledConfigBuff.String(), c); err != nil {
		glog.Warningf("routing: config decode type=%v err=%v", "toml", err)

		return err
	}

	// Известны текущие настройки виджета

	//
	// CONTENT
	//

	contentRaw, err := extractPropertyByName(widget, "Content")

	if err != nil {
		glog.Warningf("execute file: get CONTENT group=%v name=%v err=%v", group, name, err)
		return err
	}

	compiledContentBuff := bytes.NewBufferString("")
	err = c.render(compiledContentBuff, "widget_content", contentRaw, c)

	if err != nil {
		return err
	}

	// Готовый контент

	// TODO: если markdown

	//     switch p.Self.GetAsString("render") {
	//         case "markdown", "md":
	//             widgetBuff.Write(github_flavored_markdown.Markdown([]byte(contentRaw)))
	// //            widgetBuff.Write(blackfriday.MarkdownCommon([]byte(widget.Raw)))
	//         default:
	//             widgetBuff.WriteString(contentRaw)
	//     }

	contentBuff := bytes.NewBufferString("")
	err = c.render(contentBuff, "widget", compiledContentBuff.String(), c)

	if err != nil {
		glog.Warningf("execute file: compile content group=%v name=%v err=%v", group, name, err)
		return err
	}

	//
	// LAYOUT if exist
	//

	if _layout, isExist := c.Self["layout"]; isExist {
		if layout, isValid := _layout.(string); isValid && len(layout) > 0 {

			c.Self = make(map[string]interface{}) // Для последующих виджетов self свой

			c.Self["Content"] = template.HTML(contentBuff.String())

			return c.RenderWidget(writer, c.Theme.BucketName, layout)
		}
	}

	_, err = writer.Write(contentBuff.Bytes())

	if err != nil {
		glog.Warningf("execute file: write content group=%v name=%v err=%v", group, name, err)
	}

	return err
}
