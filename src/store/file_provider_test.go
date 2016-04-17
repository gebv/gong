package store

import "testing"

// import "strconv"

func TestCreateAndFindBucket(t *testing.T) {
	initConfigAndDataBase()

	var file *File
	var err error

	extId := "ExtId"

	dto := NewCreateFileDTO()
	dto.ExtId = extId
	dto.Description = "Description #1"
	dto.AddTags("1")

	dto.SetOneProps("ExtId", dto.ExtId)

	file = NewFile()

	if err := file.TransformFrom(dto); err != nil {

		t.Errorf("create bucket: transform err=%v, bucket=%v", file)
		return
	}

	file.AddCollections(CollNameBucket)

	if err := CreateFile(file); err != nil {
		t.Errorf("create bucket: err=%v, bucket=%v", file)
	}

	// Поиск созданного бакета

	// NewOrLoadBucket: Поиск по ext_id
	file, err = NewOrLoadBucket(extId)

	if err != nil {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	if file.Id.String() != "1" {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	// Check files
	if file.ExtId != extId {
		t.Error("find bucket: not expected ExtId")
		return
	}

	if file.GetOneProps("ExtId").(string) != extId {
		t.Error("find bucket: not expected props")
		return
	}

	if !file.IncludeTags("1") || file.IncludeTags("2") {
		t.Error("find bucket: not expected tags")
		return
	}

	// FindFileByExtId: Поиск по ext_id
	file, err = FindFileByExtId(CollNameBucket, extId)

	if err != nil {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	if file.Id.String() != "1" {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	// Check files
	if file.ExtId != extId {
		t.Error("find bucket: not expected ExtId")
		return
	}

	if file.GetOneProps("ExtId").(string) != extId {
		t.Error("find bucket: not expected props")
		return
	}

	if !file.IncludeTags("1") || file.IncludeTags("2") {
		t.Error("find bucket: not expected tags")
		return
	}

	// FindFile: Поиск по ext_id
	file, err = FindFile(CollNameBucket, extId)

	if err != nil {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	if file.Id.String() != "1" {
		t.Errorf("find bucket: by ext_id=%v, err=%v", extId, err)
		return
	}

	// Check files
	if file.ExtId != extId {
		t.Error("find bucket: not expected ExtId")
		return
	}

	if file.GetOneProps("ExtId").(string) != extId {
		t.Error("find bucket: not expected props")
		return
	}

	if !file.IncludeTags("1") || file.IncludeTags("2") {
		t.Error("find bucket: not expected tags")
		return
	}

	// Поиск по id
	file, err = FindFileById("1")

	if err != nil {
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	if file.Id.String() != "1" {
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	// Check files
	if file.ExtId != extId {
		t.Error("find bucket: not expected ExtId")
		return
	}

	if file.GetOneProps("ExtId").(string) != extId {
		t.Error("find bucket: not expected props")
		return
	}

	if !file.IncludeTags("1") || file.IncludeTags("2") {
		t.Error("find bucket: not expected tags")
		return
	}

	// Обновление файла
	newDescription := "new description"
	newPropsValue := "newprops"
	newPropsField := "newpropsfield"
	updatePropsField := "newvaluevaluenew"

	file.SetDescription(newDescription)
	file.SetOneProps("ExtId", updatePropsField)
	file.SetOneProps(newPropsField, newPropsValue)
	// file.SetExtId(newExtId)

	if err := UpdateFile(file); err != nil {
		t.Errorf("update bucket: id=%v, err=%v", file.Id.String(), err)

		return
	}

	// Поиск по id
	file, err = FindFileById("1")

	if err != nil {
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	if file.GetDescription() != newDescription {
		t.Error("find bucket: not expected description")
		return
	}

	if file.GetOneProps("ExtId").(string) != updatePropsField {
		t.Error("find bucket: not expected props")
		return
	}

	if file.GetOneProps(newPropsField).(string) != newPropsValue {
		t.Error("find bucket: not expected props")
		return
	}

	// УДаление bucket 1

	err = DeleteFile(file)

	if err != nil {
		t.Errorf("delete bucket: by id=%v, err=%v", 1, err)
		return
	}

	// Поиск по id
	file, err = FindFileById("1")

	if err != nil {
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	if file.IsRemoved != true {
		t.Error("find bucket: not expected IsRemoved")
		return
	}

	// Поиск по ext_id (в случае поиска по ext_id запись будет не найдена)

	file, err = FindFileByExtId(CollNameBucket, extId)

	if err != ErrNotFound {
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	file, err = NewOrLoadBucket(extId)

	if err != nil {
		// Так как ищем по id а не по ext_id
		t.Errorf("find bucket: by id=%v, err=%v", 1, err)
		return
	}

	clearDatatestFiles()
}

func TestCreatingAndFindBuckets(t *testing.T) {
	initConfigAndDataBase()

	a0 := 20

	ids := LoadBuckets(a0)

	a1 := 30
	a2 := 5
	a3 := ids[3]
	a4 := ids[13]

	LoadFiles(a3, a4, a1, a2)

	filter := NewSearchFileter()
	filter.AddCollections(CollNameBucket)

	result := SearchPerPage(filter)

	if result.Total != a0 {
		t.Errorf("search buckets: not expected count items")
		return
	}

	// кол-во в определенном бакете

	filter = NewSearchFileter()
	filter.AddCollections(CollNameFile)
	filter.AddCollections(a3)

	result = SearchPerPage(filter)

	if result.Total != a1 {
		t.Errorf("search buckets: not expected count items")
		return
	}

	// кол-во в определенном бакете

	filter = NewSearchFileter()
	filter.AddCollections(a3)
	filter.AddCollections(CollNameFile)

	result = SearchPerPage(filter)

	if len(result.Items) != itemsPerPage {
		t.Errorf("search buckets: not expected count items")
		return
	}

	if result.HasNext != true {
		t.Errorf("search buckets: not expected search result")
		return
	}

	if result.NextPage != 1 {
		t.Errorf("search buckets: not expected search result")
		return
	}

	// кол-во в определенном бакете

	filter = NewSearchFileter()
	filter.AddCollections(a4)
	filter.AddCollections(CollNameFile)

	result = SearchPerPage(filter)

	if result.Total != a1/a2 {
		t.Errorf("search buckets: not expected count items")
		return
	}

	if result.HasNext != false {
		t.Errorf("search buckets: not expected search result")
		return
	}

	// Search

	filter = NewSearchFileter()
	filter.AddCollections(a3)
	filter.AddCollections(CollNameFile)
	filter.SetPage(1)

	result = SearchPerPage(filter)

	if len(result.Items)+itemsPerPage != a1 {
		t.Errorf("search buckets: not expected count items")
		return
	}

	if result.HasNext != false {
		t.Errorf("search buckets: not expected search result")
		return
	}

	// Удаление одного bucket и проверка поиска

	file := NewFile()
	file.Id.SetString("3")

	if err := DeleteFile(file); err != nil {
		t.Errorf("remove bucket: err=%v", err)
		return
	}

	filter = NewSearchFileter()
	filter.AddCollections(CollNameBucket)

	result = SearchPerPage(filter)

	if result.Total != a0-1 {
		t.Errorf("search buckets: not expected search result")
		return
	}

	clearDatatestFiles()
}

// func TestCreateFileToRemovedBucket(t *testing.T) {
// 	initConfigAndDataBase()

// 	a0 := 20

// 	ids := LoadBuckets(a0)

// 	a1 := 30
// 	a2 := 5
// 	a3 := ids[3]
// 	a4 := ids[13]

// 	LoadFiles(a3, a4, a1, a2)

// 	file := NewFile()
// 	file.Id.SetString("4")

// 	DeleteFile(file)

// 	// Добавление файла в 4 бакет

// 	file = NewFile()
// 	file.ExtId = "newfileectid"
// 	file.Description = "description"
// 	file.AddCollections(ids[4])

// 	err := CreateFile(file)

// 	if err == nil {
// 		t.Error("create file: not expected result")
// 		return
// 	}

// 	clearDatatestFiles()
// }
