package server

import (
    "github.com/shurcooL/github_flavored_markdown"
    "html/template"
    texttemplate "text/template"
    "github.com/golang/glog"
    "github.com/BurntSushi/toml"
    // "net/http"
    "bytes"
    "io"
	"fmt"
    "store"
)

const (
    // CountWidgetPerPage count widgets per page
    CountWidgetPerPage = 250
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

func NewContext() *Context {
    model := new(Context)
    model.V = make(map[string]interface{})
    model.Self = NewConfig()
    
    // model.ResponseWriter = w
    // model.Request = r
    
    model.Funcs = template.FuncMap{
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
        "widget": func(ctx *Context, args ...interface{}) template.HTML {
            
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
            ctx.V["_Group"] = _group 
            ctx.V["_Name"] = _key 
            if err := ctx.RenderWidget(buff, _group, _key); err != nil {
                glog.Warningf("template funcs: widget: group=%v, kye=%v, err=%v", _group, _key, err)
            }
            
            // TODO: render widget
            return template.HTML(buff.String())
        },
    }
    
    return model
}

type Context struct {
    V map[string]interface{}
    Self Config
    Funcs template.FuncMap
    // Request *http.Request
    // ResponseWriter http.ResponseWriter 
    
    TraceWidgets []string
    depth int
}

func (p *Context) RenderWidget(writer io.Writer, group, name string) error {
    p.TraceWidgets = append(p.TraceWidgets, name)
    
    if p.depth > CountWidgetPerPage {
        return fmt.Errorf("more widgets per page (for example 'layout' is looped, clear 'layout')")
    }
    p.depth++
    
    
    widget, err := store.GetByPath(group, name)
    
    if err != nil {
        return err
    }
    
    configBuff := bytes.NewBufferString("")
    
    _configRaw, ok := widget.Props["Config"]
    configRaw := ""
    
    if ok {
        configRaw = _configRaw.(string)
    }
    
    ttmpl, err := texttemplate.New("compilation_config").Funcs(texttemplate.FuncMap(p.Funcs)).Parse(configRaw)
        
    if err != nil {
        return err
    }
    
    tt := texttemplate.Must(ttmpl, err)
    if err := tt.Execute(configBuff, p); err != nil {
        glog.Warningf("render widget: toml template compile: %v", err)
    }
    
    if _, err := toml.Decode(configBuff.String(), &p.V); err != nil {
        glog.Warningf("render widget: toml decode: %v", err)
    }
    
    // settings current widget
    p.Self.SetNewConfig(p.V["self"])
    delete(p.V, "self")   
    
    widgetBuff := bytes.NewBufferString("")
    _contentRaw, ok := widget.Props["Config"]
    contentRaw := ""
    
    if ok {
        contentRaw = _contentRaw.(string)
    }
    
    switch p.Self.GetAsString("render") {
        case "markdown", "md":
            widgetBuff.Write(github_flavored_markdown.Markdown([]byte(contentRaw)))
//            widgetBuff.Write(blackfriday.MarkdownCommon([]byte(widget.Raw)))
        default:
            widgetBuff.WriteString(widget.Props["Content"].(string))
    }
    
    if layout := p.Self.GetAsString("layout"); len(layout) > 0 {
        contentBuff := bytes.NewBufferString("")
        
        tmpl, err := template.New("compilation_page_layout").Funcs(p.Funcs).Parse(widgetBuff.String());
        
        if err != nil {
            glog.Warningf("render widget: compile template: %s, \n===\n%s\n===\n", name, []byte(contentRaw))
            return err
        }
        
        template.Must(tmpl, err).Execute(contentBuff, p)
        
        p.V["Content"] = template.HTML(contentBuff.String()) 
        
        // TODO: Вынести константу layoyts в настройки
        
        return p.RenderWidget(writer, "layouts", layout) 
    }
    
    tmpl, err := template.New("compilation_data").Funcs(p.Funcs).Parse(widgetBuff.String())
    
    if err != nil {
        return err
    }
    
    t := template.Must(tmpl, err)
    
    return t.Execute(writer, p) 
}