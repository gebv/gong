package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	// "github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

//dico errors
//config.toml
//[[errors]]
//name="not found"

//[[errors]]
//name="unknown"
//message="unknown error, see details in the log"
//comment="unknown error, see details in the log"

//[[errors]]
//name="not allowed"

//[[errors]]
//name="not supported"

//[[errors]]
//name="not valid_data"

//[[errors]]
//name="rejected"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  errors

var ErrNotFound = fmt.Errorf("not_found")

// ErrUnknown unknown error, see details in the log.
var ErrUnknown = fmt.Errorf("unknown error, see details in the log")

var ErrNotAllowed = fmt.Errorf("not_allowed")

var ErrNotSupported = fmt.Errorf("not_supported")

var ErrNotValidData = fmt.Errorf("not_valid_data")

var ErrRejected = fmt.Errorf("rejected")

//<<<AUTOGENERATE.DICO

var AllowedCollectionNames = map[string]bool{
	CollNameBucket: true,
	CollNameFile:   true,
}

var DisabledFileNames = map[string]bool{
	".metainfo": true,
}

// FindFile поиск файла (не удланный) в bucket по его внешнему идентификатору (ext_id)
func FindFile(bucketId, fileId string) (*File, error) {
	// if !isAllowedCollectionName(bucketId) {
	// 	glog.Errorf("manager: not allowed collection name '%v'", bucketId)
	// 	return nil, ErrNotAllowed
	// }

	if !isAllowedFileName(fileId) {
		glog.Errorf("manager: not allowed file name '%v'", fileId)
		return nil, ErrNotAllowed
	}

	// Search

	filter := NewSearchFileter()
	filter.SetHasEnabled(true)

	if bucketId == CollNameBucket {
		// Ищем bucket

		bucketIdAsUUID := uuid.FromStringOrNil(fileId)

		if !uuid.Equal(bucketIdAsUUID, uuid.Nil) {
			// fileId является bucketId типа UUID
			// Ищем bucket по прямой ссылке

			return FindFileById(bucketIdAsUUID.String())
		}

		// Ищем bucket по fileId как ExtId bucket

		filter.AddCollections(CollNameBucket) // Ищем среди buckets
		filter.SetExtId(fileId)
	} else {
		// ищем file в buckets

		var bucketFile *File
		var err error

		// Если bucketId как UUID
		bucketIdAsUUID := uuid.FromStringOrNil(bucketId)

		if !uuid.Equal(bucketIdAsUUID, uuid.Nil) {
			// bucketId как uuid

			bucketFile, err = FindFileById(bucketIdAsUUID.String())

			if err != nil {
				if err == ErrNotFound {
					return nil, ErrNotFound
				}

				glog.Errorf("manager: find file name=%v, err=%v", fileId, err)
				return nil, ErrUnknown
			}
		} else {
			// fileId не uuid

			bucketFile, err = FindFile(CollNameBucket, bucketId)

			if err != nil {
				glog.Errorf("manager: error find bucket by name=%v, err=%v", bucketId, err)
				return nil, ErrUnknown
			}

			filter.SetExtId(fileId)
		}

		// Если fileId как UUID
		fileIdAsUUID := uuid.FromStringOrNil(fileId)

		if !uuid.Equal(fileIdAsUUID, uuid.Nil) {

			return FindFileById(fileIdAsUUID.String())
		}

		filter.AddCollections(CollNameFile)
		filter.AddCollections(bucketFile.Id.String())
	}

	// index = uniq(bucketId+fileId)
	filter.SetSize(1)

	searchRes := SearchPerPage(filter)

	if len(searchRes.GetItems()) != 1 {

		if len(searchRes.GetItems()) == 0 {
			return nil, ErrNotFound
		}

		glog.Errorf("manager: search file bucket=%v, fileId=%v, err=%v", bucketId, fileId, "expected one record")

		return nil, ErrUnknown
	}

	file := searchRes.GetItems()[0]

	return file, nil
}

// FindFileById поиск файла по его идетификатору
func FindFileById(id string) (*File, error) {
	file := NewFile()

	return file, findOne(id, file)
}

// NewOrLoadBucket создать bucket
//  bucketName string id или ext_id
func NewOrLoadBucket(bucketName string) (*File, error) {
	// если существует файл, загружаем его
	// >> поиск по бакетам, по ext_id или id равным bucketName

	file, err := FindFile(CollNameBucket, bucketName)

	if err == ErrNotFound {
		file = NewFile()
		file.AddCollections(CollNameBucket)

	} else if err != nil {

		glog.Errorf("manager: find bucket=%v, err=%v", bucketName, err)
		return nil, ErrUnknown
	}

	return file, nil
}

// NewOrLoadFileOfBucket создать файл в bucket
//  bucketName string id или ext_id
//  fileName string id или ext_id
func NewOrLoadFileOfBucket(bucketName, fileName string) (*File, error) {
	if !isAllowedFileName(fileName) {
		glog.Errorf("manager: not allowed file name '%v'", fileName)
		return nil, ErrNotAllowed
	}

	// if !isAllowedCollectionName(bucketName) {
	// 	glog.Errorf("manager: not allowed collection name '%v'", bucketName)
	// 	return nil, ErrNotAllowed
	// }

	// Если файл уже такой существуует, загрузить его
	// Загружаем bucketinfo и проверяем доступность его

	file, err := FindFile(bucketName, fileName)

	if err == ErrNotFound {
		file = NewFile()

		bucketFile, err := FindFile(CollNameBucket, bucketName)

		if err != nil {
			return nil, err
		}

		file.AddCollections(CollNameFile)
		file.AddCollections(bucketFile.Id.String())
	} else if err != nil {
		glog.Errorf("manager: find file=%v from bucket=%v, err=%v", fileName, bucketName, err)

		return nil, ErrUnknown
	}

	return file, nil
}

// CreateFile создать файл
func CreateFile(bucketName string, file *File) (err error) {
	if !file.IsNew() {
		glog.Errorf("manager: rejected creating file because file not new file=%v", file.Id.String())
		return ErrRejected
	}

	var bucketFile *File

	if file.IncludeCollections(CollNameFile) {
		bucketFile, err = FindFile(CollNameBucket, bucketName)

		if err != nil {
			glog.Errorf("manager: rejected creating file because bucket=%v not existing or another reason err=%v", bucketName, err)
			return ErrRejected
		}

		file.AddCollections(bucketFile.Id.String())
	} else if file.IncludeCollections(CollNameBucket) && bucketName == CollNameBucket {

		glog.Infof("manager: create new bucket=%v", file.GetExtId())

	} else {
		glog.Errorf("manager: rejected creating file because not valid collections, colls=%v", file.Collections)
		return ErrNotValidData
	}

	fileExisting, err := FindFile(CollNameBucket, file.GetExtId())

	if err != ErrNotFound {
		glog.Errorf("manager: rejected creating file because ext_id existing from bucket=%v, bucket_id=%v, existing_file_id=%v, err=%v", bucketName, bucketFile.Id.String(), fileExisting.Id.String(), err)

		return ErrRejected
	}

	file.SetId(uuid.NewV1())
	file.BeforeCreated()

	return UpsertFile(file)
}

// UpsertFile обновить файл
func UpsertFile(file *File) error {

	if file.IsNew() {
		glog.Errorf("manager: rejected upsert file because file is new")
		return ErrNotValidData
	}

	file.BeforeUpdated()

	if err := updateOne(file.Id.String(), file); err != nil {
		glog.Errorf("upsert: error update err=%v, model=%v", err, file)

		return err
	}

	return nil
}

// Delete помечаем файл как удаленный
func Delete(file *File) error {
	if file.IsNew() {
		glog.Errorf("manager: rejected upsert file because file is new")

		return ErrNotValidData
	}

	file.BeforeDeleted()

	return UpsertFile(file)
}

func findOne(id string, file *File) error {
	return Store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))

		var buff = bytes.NewBuffer([]byte{})
		dec := gob.NewDecoder(buff)

		buff.Write(b.Get([]byte(id)))

		return dec.Decode(file)
	})
}

func updateOne(id string, file *File) error {
	return Store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))

		var buff = bytes.NewBuffer([]byte{})
		enc := gob.NewEncoder(buff)

		if err := enc.Encode(file); err != nil {
			glog.Errorln("encode: ", err)

			return err
		}

		if err := b.Put([]byte(id), buff.Bytes()); err != nil {
			glog.Errorln("update: ", err)

			return err
		}

		if err := Search.Index(id, file); err != nil {
			glog.Errorln("update: ", err)

			return err
		}

		return nil
	})
}

// Helpfull functions

func isAllowedCollectionName(name string) bool {
	return AllowedCollectionNames[name]
}

func isAllowedFileName(name string) bool {
	return !DisabledFileNames[name]
}
