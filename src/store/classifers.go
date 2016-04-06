package store

import (
    "github.com/golang/glog"
    "github.com/satori/go.uuid"
    "time"
)

func NewClassifer() *Classifer {
    model := new(Classifer)
    model.Props = make(map[string]interface{})
    
    return model
}

type Classifer struct {
    ClassiferId uuid.UUID
    ExtId string 
    
    Title string
    
    Path string
    
    Props map[string]interface{}
    Tags []string
    
    IsRemoved bool
    CreatedAt time.Time
    UpdatedAt time.Time
}


func (*Classifer) Type() string {
    return "Classifer"
}


// BeforeCreated перед созданием заполняем поля по умолчанию
func (d *Classifer) BeforeCreated() {
    d.CreatedAt = time.Now()
}

// BeforeUpdated перед обновлением заполняем поля по умолчанию
func (d *Classifer) BeforeUpdated() {
    d.UpdatedAt = time.Now()
}

func (d *Classifer) BeforeDeleted() {
    d.IsRemoved = true
}


func (c *Classifer) TransformFrom(in interface{}) {
    switch in.(type) {
        case *UpdateClassiferDTO:
            dto := in.(*UpdateClassiferDTO)
            
            // c.ExtId = dto.ExtId
            c.Title = dto.Title
            // c.Path = dto.Path
            c.Tags = dto.Tags
            c.Props = make(map[string]interface{})
            
            for key, value := range dto.Props {
                c.Props[key] = value
            }
            
        case *CreateClassiferDTO:
            dto := in.(*CreateClassiferDTO)
            
            c.ExtId = dto.ExtId
            
            if len(c.ExtId) == 0 {
                c.ExtId = uuid.NewV4().String()
            }
            
            c.Title = dto.Title
            c.Path = dto.Path
            c.Tags = dto.Tags
            c.Props = make(map[string]interface{})
            
            for key, value := range dto.Props {
                c.Props[key] = value
            }
        default:
            glog.Warningf("Не поддерживается тип %T", in)
    }
}

func NewCreateClassiferDTO() *CreateClassiferDTO {
    model := new(CreateClassiferDTO)
    model.Props = make(map[string]interface{})
    
    return model
}

type CreateClassiferDTO struct {
    ExtId string
    
    Title string
    
    Path string
    
    Props map[string]interface{}
    Tags []string
}

func NewUpdateClassiferDTO() *UpdateClassiferDTO {
    model := new(UpdateClassiferDTO)
    model.Props = make(map[string]interface{})
    
    return model
}

type UpdateClassiferDTO CreateClassiferDTO