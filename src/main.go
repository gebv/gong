package main

import (
	"flag"
	"github.com/golang/glog"

	"utils"

	"server"
	"store"
)

var fileConfig = flag.String("config", "config.toml", "Файл настроек приложения")

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

/*
//

import "github.com/blevesearch/bleve"
import _ "github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
import _ "github.com/blevesearch/bleve/analysis/analyzers/standard_analyzer"
// import "strconv"

type Item struct {
    Int int
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

    // for i := 0; i < 1000; i++ {
    //     index.Index(strconv.Itoa(i), &Item{i, []string{"/url/text/foo"+strconv.Itoa(i)+"/bar"}, "-url-text-foo1-bar " + strconv.Itoa(i)})
    // }

    // query := bleve.NewQueryStringQuery("/url/text/foo50/bar")
    var min, max float64
    min = 30
    max = 50
    query := bleve.NewNumericRangeQuery(&min, &max).SetField("Int")
    search := bleve.NewSearchRequest(query)
    searchResults, _ := index.Search(search)

    println(searchResults.String())
}
*/
