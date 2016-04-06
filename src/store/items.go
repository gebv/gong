package store

import (
    "time"
    "github.com/satori/go.uuid"
    "github.com/golang/glog"
)

func NewAttrValue(key, value string) Attribute {
    model := Attribute{}
    model.Key = key
    model.Value = value
    
    return model
}

func NewAttrValues(key string, values ...string) Attribute {
    model := Attribute{}
    model.Key = key
    model.Values = values
    
    return model
}

type Attribute struct {
    Key string
    Value string
    Values []string
}

func NewItem() *Item {
    model := new(Item)
    model.Props = make(map[string]interface{})
    
    return model
}

type Item struct {
    ItemId uuid.UUID
    
    // Для внешних зависимостей
    ExtId string `json:",omitempty"`
    
    // Для пользователя
    Articul string `json:",omitempty"`
    
    Title string `json:",omitempty"`
    
    Categories []string `json:",omitempty"` // classifer path
    Attributes []Attribute `json:",omitempty"`
    
    Props map[string]interface{}        
    Tags []string `json:",omitempty"`
    
    IsRemoved bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (*Item) Type() string {
    return "Item"
}

// BeforeCreated перед созданием заполняем поля по умолчанию
func (d *Item) BeforeCreated() {
    d.CreatedAt = time.Now()
}

// BeforeUpdated перед обновлением заполняем поля по умолчанию
func (d *Item) BeforeUpdated() {
    d.UpdatedAt = time.Now()
}

func (d *Item) BeforeDeleted() {
    d.IsRemoved = true
}

// AddCategory добавить к позиции категорию
func (i *Item) AddCategory(category string) {
    i.Categories = append(i.Categories, category)
}

// AddAttrValue добавить к позиции аттрибут
func (i *Item) AddAttrValue(key, value string) {
    i.Attributes = append(i.Attributes, NewAttrValue(key, value))
}

// AddAttrValues добавить к позиции аттрибут
func (i *Item) AddAttrValues(key string, values ...string) {
    i.Attributes = append(i.Attributes, NewAttrValues(key, values...))
}

func (i *Item) TransformFrom(in interface{}) {
    switch in.(type) {
        case *UpdateItemDTO:
            dto := in.(*UpdateItemDTO)
            
            i.ExtId = dto.ExtId
            i.Title = dto.Title
            i.Tags = dto.Tags
            i.Props = make(map[string]interface{})
            
            for key, value := range dto.Props {
                i.Props[key] = value
            }
            
            i.Attributes = dto.Attributes
        case *CreateItemDTO:
            dto := in.(*CreateItemDTO)
            
            i.ExtId = dto.ExtId
            
            if len(i.ExtId) == 0 {
                i.ExtId = uuid.NewV1().String()
            }
            
            i.Title = dto.Title
            i.Categories = dto.Categories
            i.Tags = dto.Tags
            i.Props = make(map[string]interface{})
            
            for key, value := range dto.Props {
                i.Props[key] = value
            }
            
            i.Attributes = dto.Attributes
        default:
            glog.Warningf("Not supported type %T", in)
    }
}

func NewCreateItemDTO() *CreateItemDTO {
    var model = new(CreateItemDTO)
    model.Props = make(map[string]interface{})
    
    return model
}

type CreateItemDTO struct {
    Item
}

func NewUpdateItemDTO() *UpdateItemDTO {
    var model = new(UpdateItemDTO)
    model.Props = make(map[string]interface{})
    
    return model
}

type UpdateItemDTO struct {
    Item
}