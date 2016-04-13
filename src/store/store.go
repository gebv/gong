package store

import (
	"encoding/gob"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"utils"

	_ "github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	_ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	_ "github.com/blevesearch/bleve/analysis/analyzers/standard_analyzer"
)

var StoreBucketName = "files"

var Search bleve.Index
var Store *bolt.DB

func InitStore() error {
	gob.Register(File{})

	var err error
	var dbSearch = utils.Cfg.Storage.DbSearch
	var dbStore = utils.Cfg.Storage.DbStore

	glog.Infof("init store: config=%#v", utils.Cfg)

	if Store, err = bolt.Open(dbStore, 0600, nil); err != nil {
		return err
	}

	Store.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(StoreBucketName))

		return err
	})

	if Search, err = bleve.Open(dbSearch); err != nil {
		// numberField := bleve.NewNumericFieldMapping()

		notAnalizedField := bleve.NewTextFieldMapping()
		notAnalizedField.Analyzer = "keyword"

		defaultMappinfField := bleve.NewTextFieldMapping()
		defaultMappinfField.Analyzer = "standard"

		fileMapping := bleve.NewDocumentMapping()
		fileMapping.AddFieldMappingsAt("Id", notAnalizedField)
		fileMapping.AddFieldMappingsAt("ExtId", notAnalizedField)
		fileMapping.AddFieldMappingsAt("Articul", notAnalizedField)
		fileMapping.AddFieldMappingsAt("Description", defaultMappinfField)

		propsMapping := bleve.NewDocumentMapping()
		propsMapping.DefaultAnalyzer = "keyword"
		propsMapping.AddFieldMappingsAt("Content", defaultMappinfField)
		propsMapping.AddFieldMappingsAt("Config", defaultMappinfField)

		fileMapping.AddFieldMappingsAt("Collections", notAnalizedField)
		// fileMapping.AddFieldMappingsAt("Attributes", notAnalizedField)
		fileMapping.AddSubDocumentMapping("Props", propsMapping) // Content and Config
		fileMapping.AddFieldMappingsAt("Tags", notAnalizedField)

		fileMapping.AddFieldMappingsAt("IsRemoved", bleve.NewBooleanFieldMapping())
		fileMapping.AddFieldMappingsAt("CreatedAt", bleve.NewDateTimeFieldMapping())
		fileMapping.AddFieldMappingsAt("UpdatedAt", bleve.NewDateTimeFieldMapping())

		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping(NameFileType, fileMapping)

		Search, err = bleve.New(dbSearch, mapping)

		return err
	}

	return err
}
