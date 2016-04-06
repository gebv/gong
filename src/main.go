package main

import (
    "flag"
    "github.com/golang/glog"
    
    "utils"
    
    "store"
    "server"
)

var fileConfig = flag.String("config", "config.toml", "Файл настроек приложения");

func main() {
    flag.Parse()
    
    if err := utils.InitConfig(*fileConfig); err != nil {
        glog.Fatal(err)
        return
    }
    
    if err := store.InitStore(); err != nil {
        glog.Fatal(err)
        return
    }
    
    server.RunServer()
}

//
/* 
import "github.com/blevesearch/bleve"
import _ "github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
import _ "github.com/blevesearch/bleve/analysis/analyzers/standard_analyzer"

type Item struct {
    Url []string
    Title string
}

func (*Item) Type() string {
    return "Item"
}

func main() {
    dbFileName := "_example.bleve"
    
    index, err := bleve.Open(dbFileName)
    
    if err == bleve.ErrorIndexPathDoesNotExist {
        notAnalyzedIndex := bleve.NewTextFieldMapping()
        notAnalyzedIndex.Analyzer = "keyword"
        
        textIndex := bleve.NewTextFieldMapping()
        textIndex.Analyzer = "standard"
        
        itemMapping := bleve.NewDocumentMapping()
        itemMapping.AddFieldMappingsAt("Url", notAnalyzedIndex)
        itemMapping.AddFieldMappingsAt("Title", textIndex)
        
        mapping := bleve.NewIndexMapping()
        // mapping.DefaultAnalyzer = "keyword"
        mapping.AddDocumentMapping("Item", itemMapping)    
        
        index ,err = bleve.New(dbFileName, mapping)
        
        if err != nil {
            panic(err)
        }
    }
    
    index.Index("1", &Item{[]string{"/url/text/foo1/bar", "url-text-foo2-bar"}, "Text more text"})
    index.Index("2", &Item{[]string{"/url/text/foo2/bar"}, "-url-text-foo1-bar"})
    
    query := bleve.NewTermQuery("bar").SetField("Url")
    search := bleve.NewSearchRequest(query)
    searchResults, _ := index.Search(search)
    
    println(searchResults.String())
}
*/