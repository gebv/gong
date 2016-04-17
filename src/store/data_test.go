package store

import (
	"github.com/golang/glog"
	"os"
	"strconv"
	"utils"
)

func initConfigAndDataBase() {
	clearDatatestFiles()

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

func clearDatatestFiles() {
	os.Remove(".testdatabase.bolt")
	os.RemoveAll(".testdatabase.bleve")
}

func LoadBuckets(countBuckets int) map[int]string {
	ids := make(map[int]string)

	initConfigAndDataBase()

	dto := NewCreateFileDTO()

	dto.AddCollections(CollNameBucket)

	for i := 12; i < countBuckets+12; i++ {
		iString := strconv.Itoa(i)
		dto.ExtId = "ExtId" + iString

		dto.Description = "Description #" + iString
		dto.AddTags(iString)

		dto.SetOneProps("ExtId", dto.ExtId)

		file := NewFile()

		if err := file.TransformFrom(dto); err != nil {

			continue
		}

		file.AddCollections(CollNameBucket)

		if err := CreateFile(file); err != nil {
			return ids
		}

		ids[i] = file.Id.String()
	}

	return ids
}

// LoadFiles загрузка файлов в 2 бакета. в первый все файлы во второй только кратные separator
func LoadFiles(bucketId1, bucketId2 string, totalFiles, separator int) {
	for i := 212; i < totalFiles+212; i++ {
		newFile := NewCreateFileDTO()
		newFile.ExtId = "file" + strconv.Itoa(i)
		newFile.Description = "description" + strconv.Itoa(i)

		// NewOrLoadFile
		file := NewFile()
		file.TransformFrom(newFile)

		file.AddCollections(CollNameFile)
		file.AddCollections(bucketId1)

		CreateFile(file)

		if i%separator == 0 {
			// NewOrLoadFile
			file = NewFile()
			file.TransformFrom(newFile)
			file.AddCollections(CollNameFile)
			file.AddCollections(bucketId2)

			CreateFile(file)
		}
	}
}
