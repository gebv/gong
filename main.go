package main

import (
	"github.com/boltdb/bolt"
    "github.com/golang/glog"
    "flag"
	"github.com/gorilla/mux"
    "github.com/gorilla/context"
    "github.com/gorilla/schema"
    "net/http"
	"html/template"
    texttemplate "text/template"
    "encoding/gob"
    "github.com/BurntSushi/toml"
    "bytes"
    "io"
    "errors"
    
//    "github.com/russross/blackfriday"
    "github.com/shurcooL/github_flavored_markdown"
	"fmt"
)

const (
    // CountWidgetPerPage count widgets per page
    CountWidgetPerPage = 250
)

var db *bolt.DB
var dbFile = flag.String("db", "/tmp/gong.db", "Path to the BoltDB file")
var addr = flag.String("bind", ":8080", "Listen addres")

func widgetName(name string) string {
    if len(name) == 0 {
        return "pages/__index"
    }
    return name
}

func setupStorage(file string) (*bolt.DB, error) {
    db, err := bolt.Open(file, 0600, nil)
    
    if err != nil {
        return nil, err
    }
    
    err = db.Update(func(tx *bolt.Tx) error {    
        if _, err := tx.CreateBucketIfNotExists([]byte("widgets")); err != nil {
            return err
        }
        
		return nil
	})

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main()  {
    flag.Parse()
    gob.Register(WidgetData{})
    
    var err error 
    if db, err = setupStorage(*dbFile); err != nil {
        glog.Fatal("setup db: ", err)
    }
    defer db.Close()
    
    r := mux.NewRouter()
    r.Path(`/@{space:(widgets)}:{action:(edit)}/{name:[A-Za-z0-9_.\-/]*}`).HandlerFunc(EditHandler).Methods("GET")
    r.Path(`/@{space:(widgets)}:{action:(edit)}/{name:[A-Za-z0-9_.\-/]*}`).HandlerFunc(OnSaveHandler).Methods("POST")
    r.Path(`/{name:[A-Za-z0-9_.\-/]*}`).HandlerFunc(ViewHandler).Methods("GET")
    
    http.ListenAndServe(*addr, r)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r);
    widgetName := widgetName(mux.Vars(r)["name"])
    ctx.V["Name"] = widgetName
    ctx.V["ShowHelpLink"] = len(r.URL.Query().Get("editable")) > 0
//    spaceName := mux.Vars(r)["space"]
    
    ctx.RenderWidget(w, widgetName)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r);
    spaceName := mux.Vars(r)["space"]
    widgetName := widgetName(mux.Vars(r)["name"])
    ctx.V["WidgetName"] = widgetName 
    ctx.V["FormAction"] = "/@"+ spaceName+ ":"+mux.Vars(r)["action"]+"/"+widgetName
    ctx.V["Widget"], ctx.V["WidgetError"] = GetWidgetData(widgetName)
    
    ctx.V["CurrentWidgetViewLink"] = "/"+widgetName
    
    t := template.Must(template.New("page_edit").Funcs(ctx.Funcs).Parse(editTpl))
    
    if err := ctx.RenderWidget(bytes.NewBufferString(""), widgetName); err != nil {
        ctx.V["Error"] = err
    }
    
    if err := t.Execute(w, ctx); err != nil {
        glog.Warningln("")
    }
}

func OnSaveHandler(w http.ResponseWriter, r *http.Request) {
//    ctx := NewContext(w, r);
    widgetName := widgetName(mux.Vars(r)["name"])
//    spaceName := mux.Vars(r)["space"]
    
    err := r.ParseForm()

    if err != nil {
        http.Error(w, "error parse form", http.StatusBadRequest)
        return
    }

    decoder := schema.NewDecoder()
    dto := &WidgetData{}
    err = decoder.Decode(dto, r.PostForm)

    if err != nil {
        http.Error(w, "error decode form data", http.StatusBadRequest)
        return
    }
    
    dto.Name = widgetName
    
    if err := SaveWidgetData(dto); err != nil {
        http.Error(w, "error save data", http.StatusBadRequest)
        return
    }
    
    // TODO: Задать верный URL редиректа
    http.Redirect(w, r, r.URL.String(), http.StatusFound)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
    model := new(Context)
    model.V = make(map[string]interface{})
    model.Self = NewConfig()
    
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
            
            if len(args) == 0 {
                return ""
            }
            
            _widgetName, ok :=  args[0].(string)
            
            if !ok {
                glog.Warningf("template funcs: widget: not valid widget name %v\n", args)
                return ""
            }
            
            buff := bytes.NewBufferString("")
            ctx.V["Name"] = _widgetName 
            ctx.RenderWidget(buff, _widgetName)
            
            // TODO: render widget
            return template.HTML(buff.String())
        },
    }
    model.ResponseWriter = w
    model.Request = r
    
    return model
}

// WidgetData настройки виджета
//  Config - toml настройки
//  Raw - html|css|js|etc
type WidgetData struct {
    Name string `schema:"-"`
    Config string `schema:"config"`
    Raw string `schema:"raw"`
}

func (w *WidgetData) gobEncode() ([]byte, error) {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err := enc.Encode(w)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func gobDecode(data []byte) (*WidgetData, error) {
    var w *WidgetData
    buf := bytes.NewBuffer(data)
    dec := gob.NewDecoder(buf)
    err := dec.Decode(&w)
    if err != nil {
        return nil, err
    }
    return w, nil
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

type Context struct {
    V map[string]interface{}
    Self Config
    Funcs template.FuncMap
    Request *http.Request
    ResponseWriter http.ResponseWriter 
    
    TraceWidgets []string
    depth int
}

// IsEditable Context mode editable
func (p *Context) IsEditable() bool {
    isEditable, ok := context.GetOk(p.Request, "flags.is_editable")
    if !ok {
        return false
    }
    return isEditable.(bool)
}

func (p *Context) RenderWidget(writer io.Writer, name string) error {
    p.TraceWidgets = append(p.TraceWidgets, name)
    
    if p.depth > CountWidgetPerPage {
        return errors.New("more widgets per page (for example 'layout' is looped, clear 'layout')")
    }
    p.depth++
    
    widget, err := GetWidgetData(name)
    
    if err != nil {
        return err
    }
    
    configBuff := bytes.NewBufferString("")
    
    ttmpl, err := texttemplate.New("compilation_config").Funcs(texttemplate.FuncMap(p.Funcs)).Parse(widget.Config)
    
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
    
    switch p.Self.GetAsString("render") {
        case "markdown", "md":
            widgetBuff.Write(github_flavored_markdown.Markdown([]byte(widget.Raw)))
//            widgetBuff.Write(blackfriday.MarkdownCommon([]byte(widget.Raw)))
        default:
            widgetBuff.WriteString(widget.Raw)
    }
    
    if layout := p.Self.GetAsString("layout"); len(layout) > 0 {
        contentBuff := bytes.NewBufferString("")
        
        tmpl, err := template.New("compilation_page_layout").Funcs(p.Funcs).Parse(widgetBuff.String());
        
        if err != nil {
            glog.Warningf("render widget: compile template: %s, \n===\n%s\n===\n", name, []byte(widget.Raw))
            return err
        }
        
        template.Must(tmpl, err).Execute(contentBuff, p)
        
        p.V["Content"] = template.HTML(contentBuff.String()) 
        
        return p.RenderWidget(writer, layout) 
    }
    
    tmpl, err := template.New("compilation_data").Funcs(p.Funcs).Parse(widgetBuff.String())
    
    if err != nil {
        return err
    }
    
    t := template.Must(tmpl, err)
    return t.Execute(writer, p) 
}

// GetWidgetData get widget by name
func GetWidgetData(name string) (*WidgetData, error) {
    var _bytes []byte
    err := db.View(func(tx *bolt.Tx) error {
        _bytes = tx.Bucket([]byte("widgets")).Get([]byte(name))
        return nil
    })
    
    if err != nil {
        return &WidgetData{Config: defaultConfig, Raw: defaultRaw}, nil
    }
    
    widget, err := gobDecode(_bytes);
    
    if err != nil {
        return &WidgetData{Config: defaultConfig, Raw: defaultRaw}, nil
    }
    
    return widget, nil
}

// SaveWidgetData save widget
func SaveWidgetData(dto *WidgetData) error {
    return db.Update(func(tx *bolt.Tx) error{
        _bytes, err := dto.gobEncode()
        if err != nil {
            return err
        }
        return tx.Bucket([]byte("widgets")).Put([]byte(dto.Name), _bytes)
    })
}

const defaultConfig = `
title = "page title" # if used layout

[self]
render = "" # or markdown
# if the render=markdown in the widget data should not have dynamic parameters

# layout = "" # current widget will be in the .V.Content variable
link_edit = "/@widgets:edit/{{.V.Name}}" # example dynamic parameter
link_title = "edit"
`
const defaultRaw = `{{if .V.ShowHelpLink}}<a href="{{.Self.link_edit}}" title="{{.V.Name}}">edit</a>{{end}}`

const editTpl = `<!DOCTYPE html>
<title>{{.V.Title}}</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/uikit/2.25.0/css/uikit.almost-flat.min.css">
<body class="uk-position-relative">
    <nav class="uk-navbar uk-navbar">
        <a href="/" class="uk-navbar-brand">/</a>
        <ul class="uk-navbar-nav">
        </ul>
        <div class="uk-navbar-content uk-h3"><a href="{{.V.CurrentWidgetViewLink}}">{{.V.WidgetName}}</a></div>
        <div class="uk-navbar-flip">
            <ul class="uk-navbar-nav">
            </ul>
        </div>
    </nav>
    <div class="uk-container">
        {{ if .V.Error}}
        <p class="uk-text-warning">{{.V.Error}}</p>
        {{ end }}
        <form class="uk-form uk-form-stacked" action="{{.V.FormAction}}" method="POST" >
            <div class="uk-form-row">
                <label class="uk-form-label"><a href="/helps/howto#widget-config" class="uk-text-muted uk-icon-question-circle uk-float-right"></a>Config (toml):</label>
                <div class="uk-form-controls">
                    <div id="editor_config" class="uk-width-1-1" style="height: 200px;">{{.V.Widget.Config}}</div>
                    <textarea id="src_config" name="config" style="display: none;">{{.V.Widget.Config}}</textarea>
                    <p class="uk-text-small"></p>
                </div>
            </div>
            <div class="uk-form-row">
                <label class="uk-form-label"><a href="/helps/howto#widget-data" class="uk-text-muted uk-icon-question-circle uk-float-right"></a>Data (html or css or javascript or text or etc):</label>
                <div class="uk-form-controls">
                    <div id="editor_html" class="uk-width-1-1" style="height: 300px;">{{.V.Widget.Raw}}</div>
                    <textarea id="src_raw" name="raw" style="display: none;">{{.V.Widget.Raw}}</textarea>
                    <p class="uk-text-small"></p>
                </div>
            </div>
            <div class="uk-form-row">
                <a href="{{.V.CurrentWidgetViewLink}}" class="uk-float-right uk-button" target="_blank">View widget</a>
                <button class="uk-button">Save</button>
                <a href="{{.V.FormAction}}" class="uk-button">Cancel</a>
            </div>
        </form>
    </div>
    <hr>
    <div class="uk-container uk-margin-top">
        <p><a href="/helps/howto#linked-widgets" class="uk-text-muted uk-icon-question-circle uk-float-right"></a>Linked widgets:</p>
        <ol>
        {{range $widgetName := .TraceWidgets}}
            <li><a href="/@widgets:edit/{{$widgetName}}">{{$widgetName}}</a></li>
        {{ end }}
        <ol>
    </div>
</body>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/ace.js"></script>
<script>
    var editor_config = ace.edit("editor_config");
    editor_config.getSession().setMode("ace/mode/toml");
    var editor_html = ace.edit("editor_html");
    editor_html.getSession().setMode("ace/mode/html");
    
    editor_config.getSession().on('change', function(){
        document.getElementById("src_config").value = editor_config.getSession().getValue()
    });
    
    editor_html.getSession().on('change', function(){
        document.getElementById("src_raw").value = editor_html.getSession().getValue()
    });
</script>
`