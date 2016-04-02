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
    //    itemMapping := bleve.NewDocumentMapping()
        
    //     // fieldMapping := bleve.NewTextFieldMapping()
    //     // fieldMapping.Analyzer = "simple"
        
    //     storeFieldOnlyMapping := bleve.NewTextFieldMapping()
    //     storeFieldOnlyMapping.Analyzer = "simple"
    //     storeFieldOnlyMapping.Index = false
    //     storeFieldOnlyMapping.IncludeTermVectors = false
    //     storeFieldOnlyMapping.IncludeInAll = false
        
    //    itemMapping.AddFieldMappingsAt("ExtId", storeFieldOnlyMapping)
    //    itemMapping.AddFieldMappingsAt("ItemId", storeFieldOnlyMapping)
    //    itemMapping.AddFieldMappingsAt("Articul", storeFieldOnlyMapping)
    //    itemMapping.AddFieldMappingsAt("Tags", storeFieldOnlyMapping)
    //    itemMapping.AddFieldMappingsAt("Categories", storeFieldOnlyMapping)
    //    itemMapping.AddFieldMappingsAt("Attributes", storeFieldOnlyMapping)
        
       mapping := bleve.NewIndexMapping()
    //    mapping.AddDocumentMapping("items", itemMapping)
       
       Search, err = bleve.New(dbSearch, mapping) 
       
       return err
    }
    
    return err
}