package store

import (
    "github.com/satori/go.uuid"
    "github.com/blevesearch/bleve"
    "github.com/boltdb/bolt"
    "github.com/golang/glog"
    "encoding/gob"
    "bytes"
    "fmt"
)

var ClassiferBucket = []byte("classifers")
var ItemBucket = []byte("items")

const (
    CLASSIFERS = "__classifers"
    PAGES = "__pages"
)

var RegisteredCategory = map[string]bool{
    CLASSIFERS: true,
    PAGES: true,
}

func GetByPath(group, key string) (*Item, error) {
    conjuncts := []bleve.Query{}
    
    if len(group) == 0 {
        group = CLASSIFERS
    }
    
    
    conjuncts = append(conjuncts, bleve.NewTermQuery(group).SetField("Categories"))
    conjuncts = append(conjuncts, bleve.NewTermQuery(key).SetField("ExtId"))
    conjuncts = append(conjuncts, bleve.NewBoolFieldQuery(false).SetField("IsRemoved"))
    
    query := bleve.NewConjunctionQuery(conjuncts)
    search := bleve.NewSearchRequest(query)
    searchResults, err := Search.Search(search)
   
    if err != nil {
        return nil, err
    }
    
    var item *Item
    
    if searchResults.Total > 0 {
        item, err = GetItem(searchResults.Hits[0].ID)
        
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("not found")
    }
    
    return item, nil
}

func GetItem(modelId string) (*Item, error) {
    var model = NewItem()
    
    return model, Store.View(func (tx *bolt.Tx) error {
        b := tx.Bucket(ItemBucket)
        
        var buff = bytes.NewBuffer([]byte{})
        dec := gob.NewDecoder(buff) 
        
        buff.Write(b.Get(uuid.FromStringOrNil(modelId).Bytes()))
        
        return dec.Decode(model)
    })
}

func CreateClassifer(dto *CreateItemDTO) (*Item, error) {
    
    return CreateItem(dto, CLASSIFERS)
}

func CreateItem(dto *CreateItemDTO, groupname string) (*Item, error) {
    var modelId = uuid.NewV1()
    var model = NewItem()
    
    // Найти классифер и добавить в него ExtId
    
    if groupname != CLASSIFERS {
        classifer, err  := GetItem(groupname)
        
        if err != nil {
            return model, err
        }

        // Указываем к чему принадлежит позиция
        dto.Categories = append(dto.Categories, classifer.ExtId)
        dto.Categories = append(dto.Categories, classifer.ItemId.String())   
    } else {
        dto.Categories = append(dto.Categories, CLASSIFERS)
    }
    
    return model, Store.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(ItemBucket)
        
        var buff = bytes.NewBuffer([]byte{})
        enc := gob.NewEncoder(buff) 
        
        model.ItemId = modelId
        
        model.TransformFrom(dto)
        model.BeforeCreated()
        model.BeforeUpdated()
        
        if err := enc.Encode(*model); err != nil {
            glog.Errorln("encode: ", err)
            
            return err
        }
        
        if err := b.Put(model.ItemId.Bytes(), buff.Bytes()); err != nil {
            glog.Errorln("save store: ", err)
            
            return err
        }
        
        if err := Search.Index(model.ItemId.String(), model); err != nil {
            glog.Errorln("save search: ", err)
            
            return err
        }    
         
        return nil
    })
}

func UpdateItem(modelId string, dto *UpdateItemDTO, makeDelete bool) (*Item, error) {
    var model = NewItem();
    
    return model, Store.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(ItemBucket)
        
        var buff = bytes.NewBuffer([]byte{})
        enc := gob.NewEncoder(buff)
	    dec := gob.NewDecoder(buff)
        
        buff.Write(b.Get(uuid.FromStringOrNil(modelId).Bytes())) 
        
        if err := dec.Decode(model); err != nil {
            glog.Errorln("decode: ", err)
            
            return err
        }
        
        if makeDelete {
            model.BeforeDeleted()
        } else {
            model.TransformFrom(dto)
        }
        
        model.BeforeUpdated()
        
        if err := enc.Encode(model); err != nil {
            glog.Errorln("encode: ", err)
            
            return err
        }
        
        if err := b.Put(model.ItemId.Bytes(), buff.Bytes()); err != nil {
            glog.Errorln("save store: ", err)
            
            return err
        }
        
        if err := Search.Index(model.ItemId.String(), model); err != nil {
            glog.Errorln("save search: ", err)
            
            return err
        }    
         
        return nil
    })
}

func DeleteItem(modelId string) (*Item, error) {
    return UpdateItem(modelId, NewUpdateItemDTO(), true)
}

func ItemSearchQuery(queryString, classifer_id string) ([]*Item, error) {
    conjuncts := []bleve.Query{}
    
    if len(classifer_id) == 0 {
        classifer_id = CLASSIFERS
    }
    
    // bleve.NewTermQuery
    conjuncts = append(conjuncts, bleve.NewTermQuery(classifer_id).SetField("Categories"))
    conjuncts = append(conjuncts, bleve.NewBoolFieldQuery(false).SetField("IsRemoved"))
    
    if len(queryString) != 0 {
        conjuncts = append(conjuncts, bleve.NewQueryStringQuery(queryString))    
    }
    
    // conjuncts[0] = bleve.NewPrefixQuery("title").SetField("Title")
    // conjuncts[1] = bleve.NewBoolFieldQuery(false).SetField("IsEnabled")
    
    query := bleve.NewConjunctionQuery(conjuncts)
    search := bleve.NewSearchRequest(query)
    searchResults, err := Search.Search(search)
    
    var items []*Item
    
    if err != nil {
        return items, err
    }
        
    for _, hit := range searchResults.Hits {
        model, err := GetItem(hit.ID)
        
        items = append(items, model)
        
        if err != nil {
            glog.Warningf("search: item_id=%v, err=%v", hit.ID, err)
        }
    }
    
    return items, nil
}