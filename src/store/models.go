package store

import "github.com/satori/go.uuid"
import "github.com/golang/glog"
import "fmt"
import "time"

var (
	CollNameBucket = "bucket"
	CollNameFile   = "file"
)

// type Equaler interface {
//  Equal(Equaler) bool
// }

//dico struct
//config.toml
// name ="File"

// [[transform]]
// type = "UpdateFileDTO"
// custom = '''
// f.SetProps(d.Props)
// '''

// [[transform.map]]
// to ="Articul"
// from = "Articul"
// [[transform.map]]
// to ="Description"
// from = "Description"
// [[transform.map]]
// to ="Tags"
// from = "Tags"

// [[transform]]
// type = "CreateFileDTO"
// custom = '''
// f.SetProps(d.Props)
// '''

// [[transform.map]]
// to ="ExtId"
// from = "ExtId"
// [[transform.map]]
// to ="Articul"
// from = "Articul"
// [[transform.map]]
// to ="Description"
// from = "Description"
// [[transform.map]]
// to ="Collections"
// from = "Collections"
// [[transform.map]]
// to ="Tags"
// from = "Tags"

// [[fields]]
// name = "Id"
// type = "uuid.UUID"

// [[fields]]
// name = "ExtId"
// type = "string"
// comment = "Представление третей стороны"

// [[fields]]
// name = "Articul"
// type = "string"
// comment = "Пользовательское представление"
// tag = '''json:",omitempty"'''

// [[fields]]
// name = "Description"
// type = "string"
// comment = "Описание"

// [[fields]]
// name = "Collections"
// type = "[]string"
// comment = "Принадлежнать к множествам, наборам, коллекциям"

// # [[fields]]
// # name = "Attrs"
// # type = "[]Attribute"
// # comment = "Аттрибуты"
// # tag = '''json:",omitempty"'''

// [[fields]]
// name = "Props"
// type = "map[string]interface{}"

// [[fields]]
// name = "Tags"
// type = "[]string"

// [[fields]]
// name = "IsRemoved"
// type = "bool"

// [[fields]]
// name = "CreatedAt"
// type = "time.Time"

// [[fields]]
// name = "UpdatedAt"
// type = "time.Time"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewFile() *File {
	model := new(File)

	model.Props = make(map[string]interface{})

	return model
}

type File struct {
	Id uuid.UUID

	// Представление третей стороны
	ExtId string

	// Пользовательское представление
	Articul string `json:",omitempty"`

	// Описание
	Description string

	// Принадлежнать к множествам, наборам, коллекциям
	Collections []string

	Props map[string]interface{}

	Tags []string

	IsRemoved bool

	CreatedAt time.Time

	UpdatedAt time.Time
}

// SetId set Id
func (f *File) SetId(v uuid.UUID) {
	f.Id = v
}

// GetId get Id
func (f *File) GetId() uuid.UUID {
	return f.Id
}

// SetExtId set ExtId
func (f *File) SetExtId(v string) {
	f.ExtId = v
}

// GetExtId get ExtId
func (f *File) GetExtId() string {
	return f.ExtId
}

// SetArticul set Articul
func (f *File) SetArticul(v string) {
	f.Articul = v
}

// GetArticul get Articul
func (f *File) GetArticul() string {
	return f.Articul
}

// SetDescription set Description
func (f *File) SetDescription(v string) {
	f.Description = v
}

// GetDescription get Description
func (f *File) GetDescription() string {
	return f.Description
}

// AddCollections add element Collections
func (f *File) AddCollections(v string) {
	if f.IncludeCollections(v) {
		return
	}

	f.Collections = append(f.Collections, v)
}

// RemoveCollections remove element Collections
func (f *File) RemoveCollections(v string) {
	if !f.IncludeCollections(v) {
		return
	}

	_i := f.IndexCollections(v)

	f.Collections = append(f.Collections[:_i], f.Collections[_i+1:]...)
}

// GetCollections get Collections
func (f *File) GetCollections() []string {
	return f.Collections
}

// IndexCollections get index element Collections
func (f *File) IndexCollections(v string) int {
	for _index, _v := range f.Collections {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeCollections has exist value Collections
func (f *File) IncludeCollections(v string) bool {
	return f.IndexCollections(v) > -1
}

// SetProps set all elements Props
func (f *File) SetProps(v map[string]interface{}) {
	f.Props = make(map[string]interface{})

	for key, value := range v {
		f.Props[key] = value
	}
}

// AddProps add element by key
func (f *File) SetOneProps(k string, v interface{}) {
	f.Props[k] = v
}

// RemoveProps remove element by key
func (f *File) RemoveProps(k string) {
	if _, exist := f.Props[k]; exist {
		delete(f.Props, k)
	}
}

// GetProps get Props
func (f *File) GetProps() map[string]interface{} {
	return f.Props
}

// ExistProps has exist key Props
func (f *File) ExistKeyProps(k string) bool {
	_, exist := f.Props[k]

	return exist
}

func (f *File) GetOneProps(k string) interface{} {
	return f.Props[k]
}

// AddTags add element Tags
func (f *File) AddTags(v string) {
	if f.IncludeTags(v) {
		return
	}

	f.Tags = append(f.Tags, v)
}

// RemoveTags remove element Tags
func (f *File) RemoveTags(v string) {
	if !f.IncludeTags(v) {
		return
	}

	_i := f.IndexTags(v)

	f.Tags = append(f.Tags[:_i], f.Tags[_i+1:]...)
}

// GetTags get Tags
func (f *File) GetTags() []string {
	return f.Tags
}

// IndexTags get index element Tags
func (f *File) IndexTags(v string) int {
	for _index, _v := range f.Tags {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeTags has exist value Tags
func (f *File) IncludeTags(v string) bool {
	return f.IndexTags(v) > -1
}

// SetIsRemoved set IsRemoved
func (f *File) SetIsRemoved(v bool) {
	f.IsRemoved = v
}

// GetIsRemoved get IsRemoved
func (f *File) GetIsRemoved() bool {
	return f.IsRemoved
}

// SetCreatedAt set CreatedAt
func (f *File) SetCreatedAt(v time.Time) {
	f.CreatedAt = v
}

// GetCreatedAt get CreatedAt
func (f *File) GetCreatedAt() time.Time {
	return f.CreatedAt
}

// SetUpdatedAt set UpdatedAt
func (f *File) SetUpdatedAt(v time.Time) {
	f.UpdatedAt = v
}

// GetUpdatedAt get UpdatedAt
func (f *File) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

func (f *File) TransformFrom(v interface{}) error {
	switch v.(type) {

	case *UpdateFileDTO:
		d := v.(*UpdateFileDTO)

		f.Articul = d.Articul
		f.Description = d.Description
		f.Tags = d.Tags
		f.SetProps(d.Props)

	case *CreateFileDTO:
		d := v.(*CreateFileDTO)

		f.ExtId = d.ExtId
		f.Articul = d.Articul
		f.Description = d.Description
		f.Collections = d.Collections
		f.Tags = d.Tags
		f.SetProps(d.Props)

	default:
		glog.Errorf("Not supported type %v", v)
		return fmt.Errorf("not_supported")
	}

	return nil
}

//<<<AUTOGENERATE.DICO

//dico struct
//config.toml
// name ="Attribute"

// [[fields]]
// name = "Key"
// type = "interface{}"
// tag = '''json:",omitempty"'''

// [[fields]]
// name = "Value"
// type = "interface{}"
// tag = '''json:",omitempty"'''

// [[fields]]
// name = "Values"
// type = "[]interface{}"
// tag = '''json:",omitempty"'''

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewAttribute() *Attribute {
	model := new(Attribute)

	return model
}

type Attribute struct {
	Key interface{} `json:",omitempty"`

	Value interface{} `json:",omitempty"`

	Values []interface{} `json:",omitempty"`
}

// SetKey set Key
func (a *Attribute) SetKey(v interface{}) {
	a.Key = v
}

// GetKey get Key
func (a *Attribute) GetKey() interface{} {
	return a.Key
}

// SetValue set Value
func (a *Attribute) SetValue(v interface{}) {
	a.Value = v
}

// GetValue get Value
func (a *Attribute) GetValue() interface{} {
	return a.Value
}

// AddValues add element Values
func (a *Attribute) AddValues(v interface{}) {
	if a.IncludeValues(v) {
		return
	}

	a.Values = append(a.Values, v)
}

// RemoveValues remove element Values
func (a *Attribute) RemoveValues(v interface{}) {
	if !a.IncludeValues(v) {
		return
	}

	_i := a.IndexValues(v)

	a.Values = append(a.Values[:_i], a.Values[_i+1:]...)
}

// GetValues get Values
func (a *Attribute) GetValues() []interface{} {
	return a.Values
}

// IndexValues get index element Values
func (a *Attribute) IndexValues(v interface{}) int {
	for _index, _v := range a.Values {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeValues has exist value Values
func (a *Attribute) IncludeValues(v interface{}) bool {
	return a.IndexValues(v) > -1
}

//<<<AUTOGENERATE.DICO
