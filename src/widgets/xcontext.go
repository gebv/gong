package widgets

import "github.com/labstack/echo"
import "github.com/golang/glog"
import "html/template"
import "store"
import "fmt"
import "bytes"
import "io"
import "github.com/BurntSushi/toml"
import "github.com/shurcooL/github_flavored_markdown"
import "strings"
import "strconv"
import "net/url"

// import "net/http"

// tplCache кеш для функции executeTpl
// var tplCache = make(map[string]*template.Template)

type FileInfo struct {
	Id          string
	ExtId       string
	Bucket      string
	Collections []string
	Description string
}

func NewC(context echo.Context) *C {
	c := new(C)
	c.Context = context
	c.Self = make(map[string]interface{})
	c.Global = make(map[string]interface{})
	c.J = make(map[string]interface{})

	c.Settings = NewAppSettings()
	c.Global = c.Settings.Global // подхватить значения global из настройки приложения

	c.funcs = template.FuncMap{
		"clear": func(interface{}) template.HTML {
			return ""
		},
		"md": func(v interface{}) template.HTML {
			var str string
			switch v.(type) {
			case string:
				str = v.(string)
			case template.HTML:
				str = string(v.(template.HTML))
			default:
				str = fmt.Sprintf("%v", v)
			}
			return template.HTML(github_flavored_markdown.Markdown([]byte(str)))
		},
		"render": func(args ...interface{}) template.HTML {
			_c := NewC(c.Context)

			_c.Settings = c.Settings // ссылка
			_c.Global = c.Global     // ссылка
			_c.J = c.J               // ссылка
			// _c.Self["templates"] = c.GetOneSelf("templates")

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

			if err := _c.ExecuteFile(buff, _group, _key, true); err != nil {
				glog.Warningf("template funcs: execute file: group=%v, kye=%v, err=%v", _group, _key, err)
			}

			c.traceFiles = append(c.traceFiles, _c.GettraceFiles()...)

			return template.HTML(buff.String())
		},
	}

	return c
}

//dico struct
//config.toml
// name = "C"
// disableConstructor = true
// extends = ["echo.Context"]

//[[fields]]
// name = "Self"
// type = "map[string]interface{}"
// tag = 'toml:"self"'

//[[fields]]
// comment = "переменная для значений render type = json"
// name = "J"
// type = "map[string]interface{}"
// tag = 'toml:"-"'

//[[fields]]
// name = "Global"
// type = "map[string]interface{}"
// tag = 'toml:"global"'

//[[fields]]
// name = "traceFiles"
// type = "[]FileInfo"
// tag = 'toml:"-"'

//[[fields]]
// name = "depth"
// type = "int"
// tag = 'toml:"-"'

//[[fields]]
// name = "isAbort"
// type = "bool"
// tag = 'toml:"-"'

//[[fields]]
// name = "funcs"
// type = "template.FuncMap"
// tag = 'toml:"-"'

//[[fields]]
// name = "Settings"
// type = "*AppSettings"
// tag = 'toml:"-"'

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

type C struct {
	echo.Context

	Self map[string]interface{} `toml:"self"`

	// переменная для значений render type = json
	J map[string]interface{} `toml:"-"`

	Global map[string]interface{} `toml:"global"`

	traceFiles []FileInfo `toml:"-"`

	depth int `toml:"-"`

	isAbort bool `toml:"-"`

	funcs template.FuncMap `toml:"-"`

	Settings *AppSettings `toml:"-"`
}

// SetSelf set all elements Self
func (c *C) SetSelf(v map[string]interface{}) *C {
	c.Self = make(map[string]interface{})

	for key, value := range v {
		c.Self[key] = value
	}

	return c
}

// AddSelf add element by key
func (c *C) SetOneSelf(k string, v interface{}) *C {
	c.Self[k] = v

	return c
}

// RemoveSelf remove element by key
func (c *C) RemoveSelf(k string) {
	if _, exist := c.Self[k]; exist {
		delete(c.Self, k)
	}
}

// GetSelf get Self
func (c *C) GetSelf() map[string]interface{} {
	return c.Self
}

// ExistSelf has exist key Self
func (c *C) ExistKeySelf(k string) bool {
	_, exist := c.Self[k]

	return exist
}

func (c *C) GetOneSelf(k string) interface{} {
	return c.Self[k]
}

func (c *C) GetOneSelfString(k string) string {
	v, exist := c.Self[k]
	if !exist {
		return ""
	}

	vv, valid := v.(string)

	if !valid {
		return ""
	}

	return vv
}

func (c *C) GetOneSelfArr(k string) []interface{} {
	v, exist := c.Self[k]

	if !exist {
		return []interface{}{}
	}

	vv, valid := v.([]interface{})

	if !valid {
		return []interface{}{}
	}

	return vv
}

func (c *C) GetOneSelfInt(k string) int {
	v, exist := c.Self[k]
	if !exist {
		return 0
	}

	vv, valid := v.(int)

	if !valid {
		return 0
	}

	return vv
}

func (c *C) GetOneSelfBool(k string) bool {
	v, exist := c.Self[k]
	if !exist {
		return false
	}

	vv, valid := v.(bool)

	if !valid {
		return false
	}

	return vv
}

// SetJ set all elements J
func (c *C) SetJ(v map[string]interface{}) *C {
	c.J = make(map[string]interface{})

	for key, value := range v {
		c.J[key] = value
	}

	return c
}

// AddJ add element by key
func (c *C) SetOneJ(k string, v interface{}) *C {
	c.J[k] = v

	return c
}

// RemoveJ remove element by key
func (c *C) RemoveJ(k string) {
	if _, exist := c.J[k]; exist {
		delete(c.J, k)
	}
}

// GetJ get J
func (c *C) GetJ() map[string]interface{} {
	return c.J
}

// ExistJ has exist key J
func (c *C) ExistKeyJ(k string) bool {
	_, exist := c.J[k]

	return exist
}

func (c *C) GetOneJ(k string) interface{} {
	return c.J[k]
}

func (c *C) GetOneJString(k string) string {
	v, exist := c.J[k]
	if !exist {
		return ""
	}

	vv, valid := v.(string)

	if !valid {
		return ""
	}

	return vv
}

func (c *C) GetOneJArr(k string) []interface{} {
	v, exist := c.J[k]

	if !exist {
		return []interface{}{}
	}

	vv, valid := v.([]interface{})

	if !valid {
		return []interface{}{}
	}

	return vv
}

func (c *C) GetOneJInt(k string) int {
	v, exist := c.J[k]
	if !exist {
		return 0
	}

	vv, valid := v.(int)

	if !valid {
		return 0
	}

	return vv
}

func (c *C) GetOneJBool(k string) bool {
	v, exist := c.J[k]
	if !exist {
		return false
	}

	vv, valid := v.(bool)

	if !valid {
		return false
	}

	return vv
}

// SetGlobal set all elements Global
func (c *C) SetGlobal(v map[string]interface{}) *C {
	c.Global = make(map[string]interface{})

	for key, value := range v {
		c.Global[key] = value
	}

	return c
}

// AddGlobal add element by key
func (c *C) SetOneGlobal(k string, v interface{}) *C {
	c.Global[k] = v

	return c
}

// RemoveGlobal remove element by key
func (c *C) RemoveGlobal(k string) {
	if _, exist := c.Global[k]; exist {
		delete(c.Global, k)
	}
}

// GetGlobal get Global
func (c *C) GetGlobal() map[string]interface{} {
	return c.Global
}

// ExistGlobal has exist key Global
func (c *C) ExistKeyGlobal(k string) bool {
	_, exist := c.Global[k]

	return exist
}

func (c *C) GetOneGlobal(k string) interface{} {
	return c.Global[k]
}

func (c *C) GetOneGlobalString(k string) string {
	v, exist := c.Global[k]
	if !exist {
		return ""
	}

	vv, valid := v.(string)

	if !valid {
		return ""
	}

	return vv
}

func (c *C) GetOneGlobalArr(k string) []interface{} {
	v, exist := c.Global[k]

	if !exist {
		return []interface{}{}
	}

	vv, valid := v.([]interface{})

	if !valid {
		return []interface{}{}
	}

	return vv
}

func (c *C) GetOneGlobalInt(k string) int {
	v, exist := c.Global[k]
	if !exist {
		return 0
	}

	vv, valid := v.(int)

	if !valid {
		return 0
	}

	return vv
}

func (c *C) GetOneGlobalBool(k string) bool {
	v, exist := c.Global[k]
	if !exist {
		return false
	}

	vv, valid := v.(bool)

	if !valid {
		return false
	}

	return vv
}

// GettraceFiles get traceFiles
func (c *C) GettraceFiles() []FileInfo {
	return c.traceFiles
}

// Setdepth set depth
func (c *C) Setdepth(v int) {
	c.depth = v
}

// Getdepth get depth
func (c *C) Getdepth() int {
	return c.depth
}

// SetisAbort set isAbort
func (c *C) SetisAbort(v bool) {
	c.isAbort = v
}

// GetisAbort get isAbort
func (c *C) GetisAbort() bool {
	return c.isAbort
}

// Setfuncs set funcs
func (c *C) Setfuncs(v template.FuncMap) {
	c.funcs = v
}

// Getfuncs get funcs
func (c *C) Getfuncs() template.FuncMap {
	return c.funcs
}

// SetSettings set Settings
func (c *C) SetSettings(v *AppSettings) {
	c.Settings = v
}

// GetSettings get Settings
func (c *C) GetSettings() *AppSettings {
	return c.Settings
}

//<<<AUTOGENERATE.DICO

func (c *C) trace(bucketName string, file *store.File) error {
	c.traceFiles = append(c.traceFiles, FileInfo{file.Id.String(), file.ExtId, bucketName, file.Collections, file.Description})

	return nil
}

func (c *C) InitSettings() error {
	return c.initSettings()
}

func (c *C) initSettings() error {
	file, err := store.NewOrLoadFile("settings", "app")

	if err != nil {
		glog.Errorf("context: error init settings %v", err)

		return err
	}

	if file.IsNew() {

		return fmt.Errorf("context: not found settings")
	}

	configRaw := file.GetOnePropsString("Config")

	if len(configRaw) == 0 {
		return fmt.Errorf("context: empty config file=%v", file.Id.String())
	}

	configBuff := bytes.NewBufferString("")

	if err := c.executeTpl(configBuff, "settings.app.config."+strconv.Itoa(file.UpdatedAt.Nanosecond()), configRaw, c); err != nil {

		return err
	}

	if err := c.decodeToml("settings.app.config."+file.UpdatedAt.String(), configBuff.String(), c.Settings); err != nil {

		return err
	}

	glog.Infof("app settings = %#v", c.Settings)

	return nil
}

func (c *C) executeTpl(writer *bytes.Buffer, key, raw string, d interface{}) error {
	tpl, err := template.New(key).Funcs(template.FuncMap(c.funcs)).Parse(raw)

	if err != nil {
		glog.Warningf("context: parse template=%v err=%v", key, err)

		return err
	}

	tt := template.Must(tpl, err)

	if err := tt.Execute(writer, d); err != nil {
		glog.Warningf("context: compile template=%v err=%v", key, err)

		return err
	}

	return nil
}

func (c *C) decodeToml(key, raw string, d interface{}) error {
	if _, err := toml.Decode(raw, d); err != nil {
		glog.Warningf("context: config decode type=%v err=%v", "toml", err)

		return err
	}

	return nil
}

func (c *C) ExecuteFile(writer io.Writer, bucketName, fileName string, contentExecute bool) error {
	if c.isAbort {

		return fmt.Errorf("context: bucket=%v, file=%v stoped, abort", bucketName, fileName)
	}

	file, err := store.FindFile(bucketName, fileName)

	if err != nil {
		return err
	}

	if err := c.trace(bucketName, file); err != nil {

		return err
	}

	c.Self["bucketName"] = bucketName
	c.Self["fileName"] = fileName
	c.Self["file"] = file

	// Config

	fileConfigRaw := file.GetOnePropsString("Config")
	fileConfigBuff := bytes.NewBufferString("")

	fileCacheKey := strconv.Itoa(file.UpdatedAt.Nanosecond())

	if err := c.executeTpl(fileConfigBuff, bucketName+":"+fileName+".Config."+fileCacheKey, fileConfigRaw, c); err != nil {

		return err
	}

	if err := c.decodeToml(bucketName+":"+fileName+".Config."+fileCacheKey, fileConfigBuff.String(), c); err != nil {

		return err
	}

	fileContentRawBuff := bytes.NewBufferString("")
	fileContentRawBuff.WriteString(file.GetOnePropsString("Content"))

	// Загрузка шаблонов, если есть

	if _tpls := c.GetOneSelfArr("templates"); len(_tpls) > 0 {

		for _, _tplName := range _tpls {

			c.Self = make(map[string]interface{}) // Для последующих виджетов self свой

			// Надо не исполнять контент что бы не терялся {{template ... }}{{end}}

			c.ExecuteFile(fileContentRawBuff, c.Settings.Template.Name, _tplName.(string), false)
		}
	}

	// Собираем контен (в контенте могут быть переменные)

	if !contentExecute {
		_, err := writer.Write(fileContentRawBuff.Bytes())

		return err
	}

	fileContentBuff := bytes.NewBufferString("")

	err = c.executeTpl(fileContentBuff, bucketName+":"+fileName+".Content."+fileCacheKey, fileContentRawBuff.String(), c)

	if err != nil {
		return err
	}

	// Собираем файл

	outputBuff := bytes.NewBufferString("")

	err = c.executeTpl(outputBuff, bucketName+":"+fileName+".Output."+fileCacheKey, fileContentBuff.String(), c)

	if err != nil {
		glog.Warningf("context: compile content group=%v name=%v err=%v", bucketName, fileName, err)
		return err
	}

	// Оборачиваем в layout

	glog.Infof("\t\t >> %v, %v, %v << ", bucketName, fileName, c.GetOneSelfString("render"))

	if layoutName := c.GetOneSelfString("layout"); len(layoutName) > 0 {
		c.Self = make(map[string]interface{}) // Для последующих виджетов self свой

		c.Self["Content"] = template.HTML(outputBuff.String())

		return c.ExecuteFile(writer, c.Settings.Template.Name, layoutName, true)
	}

	// Записываем результат

	writer.Write(outputBuff.Bytes())

	if err != nil {
		glog.Warningf("context: write content group=%v name=%v err=%v", bucketName, fileName, err)
	}

	return err
}

// Abort()

// IsPost() bool
// IsGet() bool
// Param(key) string
// FormValue(key) string
// QueryValue(key) string
// Bind(interface) error
// Settings(key) *store.File

// SetSession(key, interface{}) error
// GetSession(key) interface{}
// DelSession(key) intefface{}

// SetCookie(key, interface{}) error
// GetCookie(key) interface{}
// DelCookie(key) intefface{}

func (c *C) IsPost() bool {
	return strings.ToLower(c.Request().Method()) == "post"
}

func (c *C) IsGet() bool {
	return strings.ToLower(c.Request().Method()) == "get"
}

func (c *C) IsPut() bool {
	return strings.ToLower(c.Request().Method()) == "put"
}

func (c *C) IsDelete() bool {
	return strings.ToLower(c.Request().Method()) == "delete"
}

func (c *C) Abort() bool {
	c.isAbort = true

	return c.isAbort
}

func (c *C) EmptyFunc() {
	return
}

// Run
func (c *C) Run(writer io.Writer) error {
	if err := c.initSettings(); err != nil {
		glog.Warningf("context: error init settings")

		return err
	}

	if c.Settings.Routing.Mode == "simple" {
		return c.ExecuteFile(writer, "pages", c.Request().URL().Path(), true)
	}

	r := NewRouter()

	for _, route := range c.Settings.Routing.Routs {
		h := NewHandlerFromString(route.Handler)

		if h.IsEmpty() {
			glog.Warningf("context: not valid handler, route=%v", route)
			continue

		}

		glog.V(2).Infof("\tAdd route: %#v", route)

		_h := r.Handle(route.Path, h)

		// Default GET
		if len(route.Methods) > 0 {

			_h.Methods(route.Methods...)
		}
	}

	var match RouteMatch

	_u, _ := url.Parse(c.Request().URI())
	glog.V(2).Infof("\tFind routing from %v", _u.String())

	if r.Match(&Request{_u, c.Request().Method()}, &match) {
		glog.V(2).Infof("\t Match route: %#v from url=%v vars=%v", match, _u.String(), match.Vars)

		var keys []string
		var values []string

		for key, value := range match.Vars {
			keys = append(keys, key)
			values = append(values, value)
		}

		c.SetParamNames(keys...)
		c.SetParamValues(values...)

		return c.ExecuteFile(writer, match.Handler.Bucket, match.Handler.File, true)
	}

	h := NewHandlerFromString(c.Settings.Routing.NotFound)

	if !h.IsEmpty() {
		return c.ExecuteFile(writer, h.Bucket, h.File, true)
	}

	return nil
}
