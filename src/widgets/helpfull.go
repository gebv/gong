package widgets

import (
	"strings"
)

var (
	AppSettingsBucketName = "settings"
	AppSettingsFileName   = "app"
)

type appCore struct {
	ThemeName string `toml:"theme"`
}

const (
	// ModeSimple ключ как роут, в качестве ключа используется URL path и ищется во всех buckets пока не встретится, иначе page not_found
	ModeSimple = "simple"

	// ModeRoute маршрутизация осуществляется полностью роутингом
	ModeRoute = "route"
)

type configRouting struct {
	Buckets []string `toml:"buckets"`
	Mode    string   `toml:"mode"`
}

type configTheme struct {
	BucketName string `toml:"name"`
}

type configPage struct {
	DefaultContentType string `toml:"default_content_type"`
	PageNotFound       string `toml:"not_found"`
	PageNotAllowed     string `toml:"not_allowed"`
}

// BFI bucket file identifier
type BFI string

func (f BFI) IsValid() bool {
	return len(strings.Fields(string(f))) == 2
}

func (f BFI) Bucket() string {
	if !f.IsValid() {
		return ""
	}

	return strings.Fields(string(f))[0]
}

func (f BFI) File() string {
	if !f.IsValid() {
		return ""
	}

	return strings.Fields(string(f))[1]
}

func (s *Context) IsSimpleModeRouting() bool {
	return s.Routing.Mode == ModeSimple
}

// DefaultRouteBucket возвращает первый bucket из списка routing.buckets
func (s *Context) DefaultRouteBucket() string {
	if len(s.Routing.Buckets) == 0 {
		return ""
	}

	return s.Routing.Buckets[0]
}

// Abort остановка выполнения цепочки файлов
func (c *Context) Abort() {
	c.isAbort = true
}

// Helpfull functions

func (c *Context) IsPost() bool {
	return strings.ToLower(c.Context.Request().Method()) == "post"
}

func (c *Context) Path() string {
	return c.Context.Request().URL().Path()
}

func (c *Context) GetQueryParam(key string) string {
	return c.Context.Request().URL().QueryParam(key)
}

func (c *Context) GetFormValue(key string) string {
	return c.Context.Request().FormValue(key)
}
