package store

import "time"

const NameFileType = "File"

// IsNew возвращает true в случае если ID = 0
func (i *File) IsNew() bool {
	// return uuid.Equal(uuid.Nil, i.Id)
	return i.Id.IsNil()
}

func (i File) Type() string {
	return NameFileType
}

// BeforeCreated перед созданием заполняем поля по умолчанию
func (d *File) BeforeCreated() {
	d.CreatedAt = time.Now()
}

// BeforeUpdated перед обновлением заполняем поля по умолчанию
func (d *File) BeforeUpdated() {
	d.UpdatedAt = time.Now()
}

func (d *File) BeforeDeleted() {
	d.IsRemoved = true
}

// Helpfull models

func NewUpdateFileDTO() *UpdateFileDTO {
	model := new(UpdateFileDTO)
	model.Props = make(map[string]interface{})
	return model
}

type UpdateFileDTO struct {
	File
}

func NewCreateFileDTO() *CreateFileDTO {
	model := new(CreateFileDTO)
	model.Props = make(map[string]interface{})
	return model
}

type CreateFileDTO struct {
	File
}

// Attributes

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

// SetValues set all elements Values
func (a *Attribute) SetValues(v []interface{}) {

	for _, value := range v {
		a.AddValues(value)
	}
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
