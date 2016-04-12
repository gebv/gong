package store

import "time"
import "github.com/satori/go.uuid"

const NameFileType = "File"

func (i *File) IsNew() bool {
	return uuid.Equal(uuid.Nil, i.Id)
}

func (i *File) Type() string {
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
