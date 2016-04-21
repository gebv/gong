package widgets

// ------------------------------------------------------------
// AppSettings
// ------------------------------------------------------------

//dico struct
//config.toml
// name = "AppSettings"
// disableConstructor = false

//[[fields]]
// name = "Template"
// type = "TemplateSetting"
// tag = 'toml:"template"'

//[[fields]]
// name = "Routing"
// type = "RoutingSetting"
// tag = 'toml:"routing"'

//[[fields]]
// name = "Global"
// type = "map[string]interface{}"
// tag = 'toml:"global"'

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewAppSettings() *AppSettings {
	model := new(AppSettings)

	model.Global = make(map[string]interface{})

	return model
}

type AppSettings struct {
	Template TemplateSetting `toml:"template"`

	Routing RoutingSetting `toml:"routing"`

	Global map[string]interface{} `toml:"global"`
}

// SetTemplate set Template
func (a *AppSettings) SetTemplate(v TemplateSetting) {
	a.Template = v
}

// GetTemplate get Template
func (a *AppSettings) GetTemplate() TemplateSetting {
	return a.Template
}

// SetRouting set Routing
func (a *AppSettings) SetRouting(v RoutingSetting) {
	a.Routing = v
}

// GetRouting get Routing
func (a *AppSettings) GetRouting() RoutingSetting {
	return a.Routing
}

// SetGlobal set all elements Global
func (a *AppSettings) SetGlobal(v map[string]interface{}) *AppSettings {
	a.Global = make(map[string]interface{})

	for key, value := range v {
		a.Global[key] = value
	}

	return a
}

// AddGlobal add element by key
func (a *AppSettings) SetOneGlobal(k string, v interface{}) *AppSettings {
	a.Global[k] = v

	return a
}

// RemoveGlobal remove element by key
func (a *AppSettings) RemoveGlobal(k string) {
	if _, exist := a.Global[k]; exist {
		delete(a.Global, k)
	}
}

// GetGlobal get Global
func (a *AppSettings) GetGlobal() map[string]interface{} {
	return a.Global
}

// ExistGlobal has exist key Global
func (a *AppSettings) ExistKeyGlobal(k string) bool {
	_, exist := a.Global[k]

	return exist
}

func (a *AppSettings) GetOneGlobal(k string) interface{} {
	return a.Global[k]
}

func (a *AppSettings) GetOneGlobalString(k string) string {
	v, exist := a.Global[k]
	if !exist {
		return ""
	}

	vv, valid := v.(string)

	if !valid {
		return ""
	}

	return vv
}

func (a *AppSettings) GetOneGlobalArr(k string) []interface{} {
	v, exist := a.Global[k]

	if !exist {
		return []interface{}{}
	}

	vv, valid := v.([]interface{})

	if !valid {
		return []interface{}{}
	}

	return vv
}

func (a *AppSettings) GetOneGlobalInt(k string) int {
	v, exist := a.Global[k]
	if !exist {
		return 0
	}

	vv, valid := v.(int)

	if !valid {
		return 0
	}

	return vv
}

func (a *AppSettings) GetOneGlobalBool(k string) bool {
	v, exist := a.Global[k]
	if !exist {
		return false
	}

	vv, valid := v.(bool)

	if !valid {
		return false
	}

	return vv
}

//<<<AUTOGENERATE.DICO

// ------------------------------------------------------------
// TemplateSetting
// ------------------------------------------------------------

//dico struct
//config.toml
// name = "TemplateSetting"
// disableConstructor = true
//[[fields]]
// name = "Name"
// type = "string"
// tag = 'toml:"name"'
//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

type TemplateSetting struct {
	Name string `toml:"name"`
}

// SetName set Name
func (t *TemplateSetting) SetName(v string) {
	t.Name = v
}

// GetName get Name
func (t *TemplateSetting) GetName() string {
	return t.Name
}

//<<<AUTOGENERATE.DICO

// ------------------------------------------------------------
// RoutingSetting
// ------------------------------------------------------------

//dico struct
//config.toml
// name = "RoutingSetting"
// disableConstructor = true

//[[fields]]
// name = "Mode"
// type = "string"
// tag = 'toml:"mode"'

//[[fields]]
// comment = "bucket по умолчанию"
// name = "BucketName"
// type = "string"
// tag = 'toml:"bucket"'

//[[fields]]
// name = "NotFound"
// type = "string"
// tag = 'toml:"not_found"'

//[[fields]]
// comment = "json|html| или ни чего по умолчанию"
// name = "DefaultRenderMode"
// type = "string"
// tag = 'toml:"render"'

//[[fields]]
// name = "NotAllowed"
// type = "string"
// tag = 'toml:"not_allowed"'

//[[fields]]
// name = "Routs"
// type = "[]AppRoute"
// tag = 'toml:"routs"'
//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

type RoutingSetting struct {
	Mode string `toml:"mode"`

	// bucket по умолчанию
	BucketName string `toml:"bucket"`

	NotFound string `toml:"not_found"`

	// json|html| или ни чего по умолчанию
	DefaultRenderMode string `toml:"render"`

	NotAllowed string `toml:"not_allowed"`

	Routs []AppRoute `toml:"routs"`
}

// SetMode set Mode
func (r *RoutingSetting) SetMode(v string) {
	r.Mode = v
}

// GetMode get Mode
func (r *RoutingSetting) GetMode() string {
	return r.Mode
}

// SetBucketName set BucketName
func (r *RoutingSetting) SetBucketName(v string) {
	r.BucketName = v
}

// GetBucketName get BucketName
func (r *RoutingSetting) GetBucketName() string {
	return r.BucketName
}

// SetNotFound set NotFound
func (r *RoutingSetting) SetNotFound(v string) {
	r.NotFound = v
}

// GetNotFound get NotFound
func (r *RoutingSetting) GetNotFound() string {
	return r.NotFound
}

// SetDefaultRenderMode set DefaultRenderMode
func (r *RoutingSetting) SetDefaultRenderMode(v string) {
	r.DefaultRenderMode = v
}

// GetDefaultRenderMode get DefaultRenderMode
func (r *RoutingSetting) GetDefaultRenderMode() string {
	return r.DefaultRenderMode
}

// SetNotAllowed set NotAllowed
func (r *RoutingSetting) SetNotAllowed(v string) {
	r.NotAllowed = v
}

// GetNotAllowed get NotAllowed
func (r *RoutingSetting) GetNotAllowed() string {
	return r.NotAllowed
}

// GetRouts get Routs
func (r *RoutingSetting) GetRouts() []AppRoute {
	return r.Routs
}

//<<<AUTOGENERATE.DICO

// ------------------------------------------------------------
// AppRoute
// ------------------------------------------------------------

//dico struct
//config.toml
// name = "AppRoute"
// disableConstructor = true

//[[fields]]
// name = "Path"
// type = "string"
// tag = 'toml:"path"'

//[[fields]]
// name = "Methods"
// type = "[]string"
// tag = 'toml:"methods"'

//[[fields]]
// comment = "bucketname filename"
// name = "Handler"
// type = "string"
// tag = 'toml:"handler"'
//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

type AppRoute struct {
	Path string `toml:"path"`

	Methods []string `toml:"methods"`

	// bucketname filename
	Handler string `toml:"handler"`
}

// SetPath set Path
func (a *AppRoute) SetPath(v string) {
	a.Path = v
}

// GetPath get Path
func (a *AppRoute) GetPath() string {
	return a.Path
}

// SetMethods set all elements Methods
func (a *AppRoute) SetMethods(v []string) {

	for _, value := range v {
		a.AddMethods(value)
	}
}

// AddMethods add element Methods
func (a *AppRoute) AddMethods(v string) {
	if a.IncludeMethods(v) {
		return
	}

	a.Methods = append(a.Methods, v)
}

// RemoveMethods remove element Methods
func (a *AppRoute) RemoveMethods(v string) {
	if !a.IncludeMethods(v) {
		return
	}

	_i := a.IndexMethods(v)

	a.Methods = append(a.Methods[:_i], a.Methods[_i+1:]...)
}

// GetMethods get Methods
func (a *AppRoute) GetMethods() []string {
	return a.Methods
}

// IndexMethods get index element Methods
func (a *AppRoute) IndexMethods(v string) int {
	for _index, _v := range a.Methods {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeMethods has exist value Methods
func (a *AppRoute) IncludeMethods(v string) bool {
	return a.IndexMethods(v) > -1
}

// SetHandler set Handler
func (a *AppRoute) SetHandler(v string) {
	a.Handler = v
}

// GetHandler get Handler
func (a *AppRoute) GetHandler() string {
	return a.Handler
}

//<<<AUTOGENERATE.DICO
