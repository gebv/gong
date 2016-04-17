package store

import (
	"bytes"
	"encoding/gob"
	// "github.com/blevesearch/bleve"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	// "github.com/satori/go.uuid"
)

var allowedCollectionNames = map[string]bool{
	CollNameBucket: true,
	CollNameFile:   true,
}

var disabledFileNames = map[string]bool{
	".metainfo": true,
}

func isAllowedCollectionName(name string) bool {
	return allowedCollectionNames[name]
}

func isAllowedFileName(name string) bool {
	return !disabledFileNames[name]
}

type FileProvider interface {
	FindFile(string, string) (*File, error)
	FindFileById(string) (*File, error)
	FindFileByExtId(string, string) (*File, error)

	NewOrLoadBucket(string) (*File, error)
	NewOrLoadFile(string, string) (*File, error)

	CreateFile(*File) error
	UpdateFile(*File) error
	DeleteFile(*File) error

	// SearchPerPage(filter *SearchFileter) *SearchResult

	// findOne(*File) error
	// updateOne(*File) error
}

func FindFile(bucketName, fileName string) (file *File, err error) {
	var bucket *File

	if bucketName != CollNameBucket {
		file, err = NewOrLoadFile(bucketName, fileName)

		if err != nil {

			if err == ErrNotFound {

				return nil, ErrNotFound
			}

			return nil, fmt.Errorf("FindFile: find bucket by id=%v:%v, err=%v", bucketName, fileName, err)
		}
	}

	if bucketName == CollNameBucket {
		bucket, err = NewOrLoadBucket(fileName)

		if err != nil {

			if err == ErrNotFound {

				return nil, ErrNotFound
			}

			return nil, fmt.Errorf("FindFile: find bucket by id=%v err=%v", fileName, err)
		}

		file = bucket
	}

	// if file.IsRemoved {
	//     return nil, fmt.Errorf("FindFile: file is removed id=%v:%v", bucketName, fileName)
	// }

	return
}

func FindFileById(fileId string) (*File, error) {
	file := NewFile()
	file.Id.SetString(fileId)

	if file.Id.IsNil() {

		return nil, ErrNotFound
	}

	return file, findOne(file)
}

// FindFileByExtId поиск по идентификатору bucket и по ext_id файла
func FindFileByExtId(bucketId, fileExtId string) (file *File, err error) {
	filter := NewSearchFileter()
	filter.SetHasEnabled(true)
	filter.AddCollections(bucketId)
	filter.SetExtId(fileExtId)

	filter.SetSize(1)

	result := SearchPerPage(filter)

	if len(result.Items) == 0 {

		return nil, ErrNotFound
	}

	return result.GetItems()[0], nil
}

// NewOrLoadBucket ищет всевозможными способами бакет, в противном случае выдает пустой бакет
func NewOrLoadBucket(bucketName string) (file *File, err error) {

	if file, err = FindFileById(bucketName); err == nil {
		return
	}

	if file, err = FindFileByExtId(CollNameBucket, bucketName); err == nil {
		return
	}

	file = NewFile()
	file.AddCollections(CollNameBucket)

	return file, nil
}

// IsExistBucket проверяет, существует ли bucket
func IsExistBucket(bucketName string) bool {
	file, _ := NewOrLoadBucket(bucketName)

	return !file.IsNew()
}

// NewOrLoadFile ищет всевозможными способами файл, в противном случае выдает пустой file (принадлежащий bucket если таковой существует)
func NewOrLoadFile(bucketName, fileName string) (file *File, err error) {
	var bucket *File

	if bucketName == CollNameBucket {

		return NewOrLoadBucket(fileName)
	}

	if bucket, err = NewOrLoadBucket(bucketName); err != nil {
		if err == ErrNotFound {

			return nil, ErrNotFound
		}

		glog.Errorf("NewOrLoadFile: id=%v:%v, err=%v", bucketName, fileName, err)

		return nil, err
	}

	if file, err = FindFileById(fileName); err == nil {

		if !file.IncludeCollections(bucket.Id.String()) {

			return nil, fmt.Errorf("file existing by id=%v, but collection is not include bucket id=%v", fileName, bucket.Id.String())
		}

		return
	}

	if file, err = FindFileByExtId(bucket.Id.String(), fileName); err == nil {

		return
	}

	file = NewFile()
	file.AddCollections(CollNameFile)

	if !bucket.IsNew() {
		file.AddCollections(bucket.Id.String())
	}

	return file, nil
}

func CreateFile(file *File) error {
	if !file.IsNew() {
		return ErrNotAllowed
	}

	return updateOne(file)
}

func UpdateFile(file *File) error {
	if file.IsNew() {
		return ErrNotAllowed
	}

	// TODO: проверять по id наличие файла?

	return updateOne(file)
}

func DeleteFile(file *File) error {
	if file.IsNew() {
		return ErrNotAllowed
	}

	file.BeforeDeleted()

	return updateOne(file)
}

func findOne(file *File) error {

	return Store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))

		var buff = bytes.NewBuffer([]byte{})
		dec := gob.NewDecoder(buff)

		buff.Write(b.Get(file.Id.Bytes()))

		return dec.Decode(file)
	})
}

func updateOne(file *File) error {
	return Store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(StoreBucketName))

		var buff = bytes.NewBuffer([]byte{})
		enc := gob.NewEncoder(buff)

		if file.IsNew() {
			_id, _ := b.NextSequence()
			file.Id.SetUint64(_id)

			file.BeforeCreated()
		}

		file.BeforeUpdated()

		if file.IsNew() {
			glog.V(2).Infof("files: create new file=%v", file)
		} else {
			glog.V(2).Infof("files: update file=%v", file)
		}

		if err := enc.Encode(file); err != nil {
			glog.Errorf("update: encode err=%v", err)

			return err
		}

		if err := b.Put(file.Id.Bytes(), buff.Bytes()); err != nil {
			glog.Errorf("update: err=%v", err)

			return err
		}

		if err := Search.Index(file.Id.String(), file); err != nil {
			glog.Errorf("update: err=%v", err)

			return err
		}

		return nil
	})
}
