package widgets

import "html/template"
import "github.com/labstack/echo"

//dico struct
//config.toml
// name = "Context"
// disableConstructor = true

//[[fields]]
// name = "Self"
// type = "map[string]interface{}"
// tag = 'toml:"self"'

//[[fields]]
// name = "Global"
// type = "map[string]interface{}"
// tag = 'toml:"global"'

//[[fields]]
// name = "Funcs"
// type = "template.FuncMap"
// tag = 'toml:"-"'

//[[fields]]
// name = "Context"
// type = "echo.Context"
// tag = 'toml:"-"'

//[[fields]]
// name = "TraceWidgets"
// type = "[]WidgetInfo"
// tag = 'toml:"-"'

//[[fields]]
// name = "depth"
// type = "int"

//[[fields]]
// name = "isAbort"
// type = "bool"

//[[fields]]
// name = "Routing"
// type = "configRouting"
// tag = 'toml:"routing"'

//[[fields]]
// name = "Theme"
// type = "configTheme"
// tag = 'toml:"theme"'

//[[fields]]
// name = "Pages"
// type = "configPage"
// tag = 'toml:"pages"'

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

type Context struct {
	Self map[string]interface{} `toml:"self"`

	Global map[string]interface{} `toml:"global"`

	Funcs template.FuncMap `toml:"-"`

	Context echo.Context `toml:"-"`

	TraceWidgets []WidgetInfo `toml:"-"`

	depth int

	isAbort bool

	Routing configRouting `toml:"routing"`

	Theme configTheme `toml:"theme"`

	Pages configPage `toml:"pages"`
}

// SetSelf set all elements Self
func (c *Context) SetSelf(v map[string]interface{}) {
	c.Self = make(map[string]interface{})

	for key, value := range v {
		c.Self[key] = value
	}
}

// AddSelf add element by key
func (c *Context) SetOneSelf(k string, v interface{}) {
	c.Self[k] = v
}

// RemoveSelf remove element by key
func (c *Context) RemoveSelf(k string) {
	if _, exist := c.Self[k]; exist {
		delete(c.Self, k)
	}
}

// GetSelf get Self
func (c *Context) GetSelf() map[string]interface{} {
	return c.Self
}

// ExistSelf has exist key Self
func (c *Context) ExistKeySelf(k string) bool {
	_, exist := c.Self[k]

	return exist
}

func (c *Context) GetOneSelf(k string) interface{} {
	return c.Self[k]
}

// SetGlobal set all elements Global
func (c *Context) SetGlobal(v map[string]interface{}) {
	c.Global = make(map[string]interface{})

	for key, value := range v {
		c.Global[key] = value
	}
}

// AddGlobal add element by key
func (c *Context) SetOneGlobal(k string, v interface{}) {
	c.Global[k] = v
}

// RemoveGlobal remove element by key
func (c *Context) RemoveGlobal(k string) {
	if _, exist := c.Global[k]; exist {
		delete(c.Global, k)
	}
}

// GetGlobal get Global
func (c *Context) GetGlobal() map[string]interface{} {
	return c.Global
}

// ExistGlobal has exist key Global
func (c *Context) ExistKeyGlobal(k string) bool {
	_, exist := c.Global[k]

	return exist
}

func (c *Context) GetOneGlobal(k string) interface{} {
	return c.Global[k]
}

// SetFuncs set Funcs
func (c *Context) SetFuncs(v template.FuncMap) {
	c.Funcs = v
}

// GetFuncs get Funcs
func (c *Context) GetFuncs() template.FuncMap {
	return c.Funcs
}

// SetContext set Context
func (c *Context) SetContext(v echo.Context) {
	c.Context = v
}

// GetContext get Context
func (c *Context) GetContext() echo.Context {
	return c.Context
}

// GetTraceWidgets get TraceWidgets
func (c *Context) GetTraceWidgets() []WidgetInfo {
	return c.TraceWidgets
}

// Setdepth set depth
func (c *Context) Setdepth(v int) {
	c.depth = v
}

// Getdepth get depth
func (c *Context) Getdepth() int {
	return c.depth
}

// SetisAbort set isAbort
func (c *Context) SetisAbort(v bool) {
	c.isAbort = v
}

// GetisAbort get isAbort
func (c *Context) GetisAbort() bool {
	return c.isAbort
}

// SetRouting set Routing
func (c *Context) SetRouting(v configRouting) {
	c.Routing = v
}

// GetRouting get Routing
func (c *Context) GetRouting() configRouting {
	return c.Routing
}

// SetTheme set Theme
func (c *Context) SetTheme(v configTheme) {
	c.Theme = v
}

// GetTheme get Theme
func (c *Context) GetTheme() configTheme {
	return c.Theme
}

// SetPages set Pages
func (c *Context) SetPages(v configPage) {
	c.Pages = v
}

// GetPages get Pages
func (c *Context) GetPages() configPage {
	return c.Pages
}

//<<<AUTOGENERATE.DICO
