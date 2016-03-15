package main

import (
    // "github.com/zenazn/goji"
    // "github.com/zenazn/goji/web"
    "github.com/boltdb/bolt"
    
    "html/template"
    "github.com/gorilla/mux"
    // "regexp"
    "flag"
    "net/http"
    "net/url"
    "log"
    "bytes"
)

var db *bolt.DB

func NewPageData() *PageData {
    model := new(PageData)
    model.V = make(map[string]interface{})
    model.Funcs = template.FuncMap{
        "eq": func(a, b interface{}) bool {
            return a == b
        },
        "widget": func(pageData *PageData, slug string) template.HTML {
            
            return renderWidget("widgets", slug, pageData)
        },
    }
    return model
}
type PageData struct {
    V map[string]interface{}
    Funcs template.FuncMap
    URL url.URL
}

const viewTpl = `<!DOCTYPE html>
<html>
	<head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/uikit/2.25.0/css/uikit.almost-flat.min.css">
	</head>
    <body class="uk-position-relative">
        {{ widget . "/top" }}
        <nav class="uk-navbar uk-navbar">
            <a href="/" class="uk-navbar-brand">{{widget . "/brand" }}</a>
            <ul class="uk-navbar-nav">
                {{ widget . "/navbar" }}
            </ul>
            <div class="uk-navbar-flip">
                <ul class="uk-navbar-nav uk-hidden-small">
                    <li><a href="/">Настройки</a></li>
                </ul>
            </div>
        </nav>
	    {{ .V.Content }}
        <footer>
        {{ widget . "/footer" }}
        </footer>
	</body>
</html>`

const editTpl = `<!DOCTYPE html>
<html>
	<head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/uikit/2.25.0/css/uikit.almost-flat.min.css">
	</head>
    <body class="uk-position-relative">
        {{ widget . "/top" }}
        <nav class="uk-navbar uk-navbar">
            <a href="/" class="uk-navbar-brand">{{widget . "/brand" }}</a>
            <ul class="uk-navbar-nav">
                {{ widget . "/navbar" }}
            </ul>
            <div class="uk-navbar-content">@{{.V.Space}} {{.V.Name}}</div>
            <div class="uk-navbar-flip">
                <ul class="uk-navbar-nav uk-hidden-small">
                    <li><a href="/">Настройки</a></li>
                </ul>
            </div>
        </nav>
        <div class="uk-container">
            <form class="uk-form" action="{{.V.FormAction}}" method="POST">
                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <textarea class="uk-width-1-1" cols="30" rows="5" placeholder="..." name="value">{{.V.Value}}</textarea>
                    </div>
                </div>
                <div class="uk-form-row">
                    <button class="uk-button">Save</button>
                    <a href="{{.V.Name}}" class="uk-button">Cancel</a>
                </div>
            </form>
        </div>
        <footer>
        {{ widget . "/footer" }}
        </footer>
	</body>
</html>`

var dbFile = flag.String("db", "/tmp/gong.db", "Path to the BoltDB file")
var addr = flag.String("bind", ":8080", "Listen addres")

func getSpecialValue(space, key string, buf *bytes.Buffer) {
    db.View(func(tx *bolt.Tx) error {
        log.Printf("get: bucket=%v key=%v value=%v", space, key, string(tx.Bucket([]byte(space)).Get([]byte(key))))
        buf.Write(tx.Bucket([]byte(space)).Get([]byte(key)))
        return nil
    })
    return
}

func getString(bucket, key string) string {
    _src := bytes.NewBuffer([]byte{})
    getSpecialValue(bucket, key, _src)
    return _src.String()
}

func renderWidget(space, key string, data *PageData) template.HTML {
    _html := bytes.NewBuffer([]byte{}) 
    _tpl := getString(space,key)
    
    t := template.Must(template.New("widget").Funcs(data.Funcs).Parse(_tpl))
    t.Execute(_html, data)
    return template.HTML(_html.String())
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
    data := NewPageData()
    data.URL = *r.URL
    data.V["Space"] = mux.Vars(r)["space"]
    data.V["Action"] = mux.Vars(r)["action"]
    data.V["Name"] = "/"+mux.Vars(r)["name"] 
    data.V["FormAction"] = "/@"+ mux.Vars(r)["space"]+ ":"+mux.Vars(r)["action"]+"/"+mux.Vars(r)["name"] 
    data.V["Value"] = getString(mux.Vars(r)["space"], "/"+mux.Vars(r)["name"] )
    t := template.Must(template.New("page").Funcs(data.Funcs).Parse(editTpl))
    t.Execute(w, data)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
    data := NewPageData()
    data.URL = *r.URL
    _bucketName := mux.Vars(r)["space"]
    data.V["Space"] = _bucketName
    _key := "/"+mux.Vars(r)["name"]
    data.V["Name"] = _key 
    
    r.ParseForm()
    log.Printf("save: bucket=%v key=%v value=%v", _bucketName, _key, r.PostForm.Get("value"))
    db.Update(func(tx *bolt.Tx) error{
        tx.CreateBucketIfNotExists([]byte(_bucketName))
        
        return tx.Bucket([]byte(_bucketName)).Put([]byte(_key), []byte(r.PostForm.Get("value")))
    })
    
    http.Redirect(w, r, _key, http.StatusFound)
}

func ViewPageHandler(w http.ResponseWriter, r *http.Request) {
    data := NewPageData()
    data.URL = *r.URL
    _key := "/"+mux.Vars(r)["name"]
    data.V["Name"] = _key
    
    _html := bytes.NewBuffer([]byte{})
    _tpl := getString("pages", _key)
    
    tContent := template.Must(template.New("page_content").Funcs(data.Funcs).Parse(_tpl))
    tContent.Execute(_html, data)
    
    data.V["Content"] = template.HTML(_html.String())
    tPage := template.Must(template.New("page").Funcs(data.Funcs).Parse(viewTpl))
    tPage.Execute(w, data)
}

func setupDB() error {
    var err error
    db, err = bolt.Open(*dbFile, 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("pages")); err != nil {
            log.Fatal(err)
        }
        
        if _, err := tx.CreateBucketIfNotExists([]byte("widgets")); err != nil {
            log.Fatal(err)
        }

		return nil
	})

	if err != nil {
		db.Close()
		return err
	}

	return nil
}

func main()  {
    flag.Parse()
    
    if err := setupDB(); err != nil {
        log.Fatal(err);
        return
    }
    defer db.Close()
    
    r := mux.NewRouter()
    r.Path(`/@{space:(pages|widgets)}:{action:(edit)}/{name:[A-Za-z0-9_.\-/]+}`).HandlerFunc(EditHandler).Methods("GET")
    r.Path(`/@{space:(pages|widgets)}:{action:(edit)}/{name:[A-Za-z0-9_.\-/]+}`).HandlerFunc(SaveHandler).Methods("POST")
    r.Path(`/{name:[A-Za-z0-9_.\-/]+}`).HandlerFunc(ViewPageHandler).Methods("GET")
    
    http.ListenAndServe(*addr, r)
}