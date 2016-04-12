package store

import (
	"fmt"
	"github.com/golang/glog"
	"os"
	"strconv"
	"testing"
	"time"
	"utils"
)

func initConfigAndDataBase() {
	os.Remove(".testdatabase.bolt")
	os.RemoveAll(".testdatabase.bleve")

	if err := utils.InitConfig("../../config/config.toml.travis"); err != nil {
		glog.Fatal(err)
		return
	}

	utils.Cfg.Storage.DbSearch = ".testdatabase.bleve"
	utils.Cfg.Storage.DbStore = ".testdatabase.bolt"

	if err := InitStore(); err != nil {
		glog.Fatal(err)
		return
	}
}

func LoadFixtures(t *testing.T) map[string]string {
	ids := make(map[string]string)

	initConfigAndDataBase()

	dto := NewCreateFileDTO()

	dto.AddCollections(CollNameBucket)

	for i := 0; i < 10; i++ {
		iString := strconv.Itoa(i)
		dto.ExtId = "ExtId" + iString

		dto.Description = "Description #" + iString
		dto.AddTags(iString)

		dto.AddCollections(iString)
		dto.SetOneProps("ExtId", dto.ExtId)

		file := NewFile()
		if err := file.TransformFrom(dto); err != nil {
			fmt.Printf("load fixtures:", err)
			continue
		}

		if err := CreateFile(CollNameBucket, file); err != nil {
			t.Errorf("create file, err=%v", err)
			return ids
		}

		ids[iString] = file.Id.String()
	}

	return ids
}

func TestSearchSimpleProcedures(t *testing.T) {
	t.Log("load fixtures...")
	ids := LoadFixtures(t)

	filter := NewSearchFileter()
	filter.AddCollections(CollNameBucket)
	res := SearchPerPage(filter)

	if res.Total != 10 {
		t.Error("not expected search result, res=%v", res)
		return
	}

	// Создать пару записей в ExtId2 bucket

	for i := 0; i < 30; i++ {
		newFile := NewCreateFileDTO()
		newFile.ExtId = "file" + strconv.Itoa(i)
		newFile.Description = "description" + strconv.Itoa(i)

		file := NewFile()
		file.TransformFrom(newFile)
		file.AddCollections(CollNameFile)

		CreateFile("ExtId2", file)
	}

	filter = NewSearchFileter()
	filter.AddCollections(ids["2"])
	res = SearchPerPage(filter)

	if res.Total != 30 {
		t.Error("not expected search result, res=%v", res)
		return
	}

	if res.Items[0].ExtId != "file29" {
		t.Error("not expected search order result, res=%v", res.Items[0])
		return
	}
}

func TestCheckSimpleProcedures(t *testing.T) {
	// 1. Load fixtures
	// 2. Поиск по ext_id, uuid
	// 3. Update file
	// 4. Check Update file
	// 5. Delete file
	// 6. Chekc delete file

	// 7. Создание файла в bucket

	t.Log("load fixtures...")
	ids := LoadFixtures(t)

	// Поиск по ExtId

	file, err := FindFile(CollNameBucket, "ExtId4")

	if err != nil {
		t.Errorf("find file by ext_id =%v, err=%v", "ExtId4", err)
		return
	}

	if file.Description != "Description #4" {
		t.Errorf("not valid field 'description'")
		return
	}

	if len(file.Tags) != 5 {
		t.Errorf("not valid length field 'tags'")
		return
	}

	if !file.IncludeCollections("0") {
		t.Errorf("not valid field 'collections'")
		return
	}

	if !file.IncludeCollections("3") {
		t.Errorf("not valid field 'collections'")
		return
	}

	if file.Props["ExtId"].(string) != file.ExtId {
		t.Errorf("not valid field 'props[id]'=%v, file_id=%v", file.Props, file.ExtId)
		return
	}

	// Поиска по UUID

	fileId := ids["3"]

	file, err = FindFile(CollNameBucket, fileId)

	if err != nil {
		t.Errorf("not found file=%v", fileId)
		return
	}

	if file.ExtId != "ExtId3" {
		t.Errorf("not valid field 'ExtId'")
		return
	}

	// Обновление
	file, err = NewOrLoadBucket("ExtId7")

	if err != nil {
		t.Errorf("not found file=%v", "ExtId7")
		return
	}

	if file.Props["ExtId"].(string) != file.ExtId {
		t.Errorf("not valid field 'props[id]'=%v, file_id=%v", file.Props, file.ExtId)
		return
	}

	file.SetOneProps("NewField", "new value")
	file.SetDescription("New description")

	if err := UpsertFile(file); err != nil {
		t.Errorf("save file err=%v", err)
		return
	}

	// Проверка обновления
	file, err = NewOrLoadBucket("ExtId7")

	if err != nil {
		t.Errorf("not found file=%v", "ExtId7")
		return
	}

	if file.Props["NewField"].(string) != "new value" {
		t.Errorf("not valid field 'props[NewFields]'=%v, file_id=%v", file.Props, file.ExtId)
		return
	}

	if file.Description != "New description" {
		t.Errorf("not valid field 'Description'")
		return
	}

	if err := Delete(file); err != nil {
		t.Errorf("delete file, err=%v", err)
		return
	}

	time.Sleep(time.Second * 1)

	// Проверка удаления
	file, err = NewOrLoadBucket("ExtId7")

	if err != nil {
		t.Errorf("not found file=%v", "ExtId7")
		return
	}

	if !file.IsNew() {
		t.Errorf("not correct delete file=%v", "ExtId7")
		return
	}

	// 7. Созданрие файла

	newFile := NewCreateFileDTO()
	newFile.ExtId = "new file"
	newFile.Description = "new description"

	file = NewFile()
	file.TransformFrom(newFile)
	file.AddCollections(CollNameFile)

	if err := CreateFile("ExtId7", file); err != ErrRejected {
		t.Errorf("create new file bucket=%v file=%v err=%v", "ExtId7", newFile.ExtId, err)
		return
	}

	if err := CreateFile("ExtId8", file); err != nil {
		t.Errorf("create new file bucket=%v file=%v err=%v", "ExtId7", newFile.ExtId, err)
		return
	}

	// Проверка созданного файла в bucket ExtId8
	file, err = NewOrLoadFileOfBucket("ExtId8", "new file")

	if err != nil {
		t.Errorf("find file, err=%v", err)
		return
	}

	if file.Description != "new description" {
		t.Errorf("not valid field 'Description'")
		return
	}
}
