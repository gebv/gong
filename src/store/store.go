package store

import (
    "utils"
    "github.com/blevesearch/bleve"
    // "github.com/golang/glog"
    "github.com/boltdb/bolt"
    "time"
    "encoding/gob"
    
    _ "github.com/blevesearch/bleve/analysis/analyzers/standard_analyzer"
    _ "github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
    _ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
)

var Search bleve.Index
var Store *bolt.DB

type DatabaseModel struct {
    IsRemoved bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

// BeforeCreated перед созданием заполняем поля по умолчанию
func (d *DatabaseModel) BeforeCreated() {
    d.CreatedAt = time.Now()
}

// BeforeUpdated перед обновлением заполняем поля по умолчанию
func (d *DatabaseModel) BeforeUpdated() {
    d.UpdatedAt = time.Now()
}

func (d *DatabaseModel) BeforeDeleted() {
    d.IsRemoved = true
}

func InitStore() error {
    gob.Register(Classifer{})
    gob.Register(Item{})
    
    var err error
    var dbSearch = utils.Cfg.Storage.DbSearch
    var dbStore = utils.Cfg.Storage.DbStore
    
    if Store, err = bolt.Open(dbStore, 0600, nil); err != nil {
        return err
    }
    
    Store.Batch(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(ItemBucket)
        
        return err
    })
     
    if Search, err = bleve.Open(dbSearch); err != nil {        
        notAnalizedField := bleve.NewTextFieldMapping()
        notAnalizedField.Analyzer = "keyword"
        
        defaultMappinfField := bleve.NewTextFieldMapping();
        defaultMappinfField.Analyzer = "standard"
        
        itemMapping := bleve.NewDocumentMapping()
        itemMapping.AddFieldMappingsAt("ItemId", notAnalizedField)
        itemMapping.AddFieldMappingsAt("ExtId", notAnalizedField)
        itemMapping.AddFieldMappingsAt("Articul", notAnalizedField)
        itemMapping.AddFieldMappingsAt("Title", defaultMappinfField)
        
        itemPropsMapping := bleve.NewDocumentMapping()
        itemPropsMapping.DefaultAnalyzer = "keyword"
        itemPropsMapping.AddFieldMappingsAt("Content", defaultMappinfField)
        itemPropsMapping.AddFieldMappingsAt("Config", defaultMappinfField)
        
        itemMapping.AddFieldMappingsAt("Categories", notAnalizedField)
        itemMapping.AddFieldMappingsAt("Attributes", notAnalizedField)
        itemMapping.AddSubDocumentMapping("Props", itemPropsMapping) // Content and Config
        itemMapping.AddFieldMappingsAt("Tags", notAnalizedField)
        
        itemMapping.AddFieldMappingsAt("IsRemoved", bleve.NewBooleanFieldMapping())
        itemMapping.AddFieldMappingsAt("CreatedAt", bleve.NewDateTimeFieldMapping())
        itemMapping.AddFieldMappingsAt("UpdatedAt", bleve.NewDateTimeFieldMapping())
        
        classiferMapping := bleve.NewDocumentMapping()
        classiferMapping.AddFieldMappingsAt("ClassiferId", notAnalizedField)
        classiferMapping.AddFieldMappingsAt("ExtId", notAnalizedField)
        classiferMapping.AddFieldMappingsAt("Title", defaultMappinfField)
        classiferMapping.AddFieldMappingsAt("Path", notAnalizedField)
        classiferMapping.AddFieldMappingsAt("Tags", notAnalizedField)
        classiferMapping.AddFieldMappingsAt("IsRemoved", bleve.NewBooleanFieldMapping())
        classiferMapping.AddFieldMappingsAt("CreatedAt", bleve.NewDateTimeFieldMapping())
        classiferMapping.AddFieldMappingsAt("UpdatedAt", bleve.NewDateTimeFieldMapping())

        mapping := bleve.NewIndexMapping()
        mapping.AddDocumentMapping("Item", itemMapping)
        mapping.AddDocumentMapping("Classifer", classiferMapping)

        Search, err = bleve.New(dbSearch, mapping) 

        return err
    }
    
    return err
}