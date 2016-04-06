package widgets

import (
    "github.com/shurcooL/github_flavored_markdown"
    "html/template"
    texttemplate "text/template"
    "github.com/golang/glog"
    "github.com/BurntSushi/toml"
    "github.com/labstack/echo"
    // "net/http"
    "bytes"
    "strings"
    "io"
	"fmt"
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
          glog.Infof("[DEV]: %v", args);
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
        "redirect": func(path string) template.HTML {
            // model.Context.Redirect(302, path)
            return template.HTML("")
        }, 
        "widget": func(ctx *Context, args ...interface{}) template.HTML {
            
            context := NewContext(ctx.Context)
            context.Global = ctx.Global // копируется ссылка значений
            
            if len(args) == 0 || len(args) < 2 {
                glog.Warningf("template funcs: min 2 args %v\n", args)
                
                return ""
            }
            
            _group, ok1 :=  args[0].(string)
            _key, ok2 :=  args[1].(string)
            
            if !ok1 || !ok2 {
                glog.Warningf("template funcs: widget: not valid widget group name or key %v\n", args)
                return ""
            }
            
            buff := bytes.NewBufferString("")
            
            context.Self["_group"] = _group 
            context.Self["_name"] = _key 
            
            if err := context.RenderWidget(buff, _group, _key); err != nil {
                glog.Warningf("template funcs: widget: group=%v, kye=%v, err=%v", _group, _key, err)
            }
            
            // TODO: render widget
            return template.HTML(buff.String())
        },
    }
    
    return model
}

// TODO: Добавить Request и Response для реализации возможности создавать обработчик запросов в виджете

type routeConfig struct {
    Methods []string `toml:"methods"`
    Collection string `toml:"collection"`
    Key string `toml:"key"`
}

type Context struct {
    Self map[string]interface{} `toml:"self"`
    Global map[string]interface{} `toml:"global"`
    
    Funcs template.FuncMap
    
    Context echo.Context
    
    TraceWidgets []string
    depth int
}

// getStringByPropKey вспомогательная функция получения значения из props
func getStringByPropKey(item *store.Item, propKey string) (string, error) {
    // TODO: Вынести в item
    src, isExist := item.Props[propKey]
    
    if _, isValid := src.(string); !isExist && !isValid {
        return "", fmt.Errorf("not sets or not valid key '%v' from widget ext_id='%v' categories=%v", propKey, item.ExtId, item.Categories)
    }
    
    return src.(string), nil
}

func (c *Context) traceWidgets(group, name string) error {
    c.TraceWidgets = append(c.TraceWidgets, group, name)
    
    if c.depth > CountWidgetPerPage {
        return fmt.Errorf("more widgets per request (for example 'layout' is looped, clear 'layout')")
    }
    
    c.depth++
    
    return nil
}

func (c *Context) render(writer *bytes.Buffer, key, raw string, data interface{}) (error) {
    tpl, err := texttemplate.New(key).Funcs(texttemplate.FuncMap(c.Funcs)).Parse(raw)
        
    if err != nil {
        glog.Warningf("render: parse template err=%v", err)
        
        return err
    }
    
    tt := texttemplate.Must(tpl, err)
    if err := tt.Execute(writer, c); err != nil {
        glog.Warningf("render: compile template err=%v", err)
        
        return err
    }
    
    return nil
}

// Execute получает widget с конфигурацией и выполяет действия
func (c *Context) Execute(writer io.Writer, path string) error {
    
    if err := c.traceWidgets(ROUTES, path); err != nil {
        return err
    }
    
    routeWidget, err := store.GetByPath(ROUTES, path)
    
    if err != nil {
        glog.Warningf("routing: error get widget name=%v err=%v", path, err)
        return fmt.Errorf("routing: not found route config by name '%v'", path)
    }
    
    configRaw, err := getStringByPropKey(routeWidget, "Config")
    
    if err != nil {
        glog.Warningf("routing: error get config name=%v err=%v", path, err)
        return err
    }
    
    compiledConfigBuff := bytes.NewBufferString("")
    err = c.render(compiledConfigBuff, "route_config", configRaw, c)
    
    if err != nil {
        return err
    }
    
    // TODO: Разные типы конфигов
    
    config := &routeConfig{}
    
    if _, err := toml.Decode(compiledConfigBuff.String(), config); err != nil {
        glog.Warningf("routing: config decode type=%v err=%v", "toml", err)
        
        return err
    }
    
    // ACL
    allowed := false
    currentMethodName := strings.ToLower(c.Context.Request().Method())
    
    for _, methodName := range config.Methods {
        if currentMethodName == strings.ToLower(methodName) {
            allowed = true
        }
    }
    
    if !allowed {
        glog.Warningf("routing: not allowed method=%v, allowed methods=%v", c.Context.Request().Method(), config.Methods)
        
        return fmt.Errorf("routing: not allowed method")
    }
    
    c.Global["_route"] = config
    
    return c.RenderWidget(writer, config.Collection, config.Key)
    
}

func (c *Context) RenderWidget(writer io.Writer, group, name string) error {
    if err := c.traceWidgets(group, name); err != nil {
        
        return err
    }
    
    widget, err := store.GetByPath(group, name)
    
    if err != nil {
        glog.Warningf("widget: error get widget group=%v name=%v err=%v", group, name, err)
        
        return err
    }
    
    c.Self["_group"] = group 
    c.Self["_name"] = name 
    
    //
    // CONFIG
    //
    
    configRaw, err := getStringByPropKey(widget, "Config")
    
    if err != nil {
        glog.Warningf("widget: get CONFIG group=%v name=%v err=%v", group, name, err)
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
    
    contentRaw, err := getStringByPropKey(widget, "Content")
    
    if err != nil {
        glog.Warningf("widget: get CONTENT group=%v name=%v err=%v", group, name, err)
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
        glog.Warningf("widget: compile content group=%v name=%v err=%v", group, name, err)
        return err
    }
    
    //
    // LAYOUT is exist
    //
    
    if _layout, isExist := c.Self["layout"]; isExist {
        if layout, isValid := _layout.(string); isValid && len(layout) > 0 {
            
            c.Self = make(map[string]interface{}) // Для последующих виджетов self свой
            
            c.Self["Content"] = template.HTML(contentBuff.String())
            
            return c.RenderWidget(writer, LAYOUTS, layout) 
        }   
    }
    
    _, err = writer.Write(contentBuff.Bytes())
    
    if err != nil {
        glog.Warningf("widget: write content group=%v name=%v err=%v", group, name, err)
    }
    
    return err
}