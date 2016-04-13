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
import "os"
import "strconv"

type Item struct {
	Int   int
	Url   []string
	Ints  []int
	Title string
}

func (*Item) Type() string {
	return "Item"
}

func main() {
	dbFileName := "_example.bleve"
	os.RemoveAll(dbFileName)

	index, err := bleve.Open(dbFileName)

	if err == bleve.ErrorIndexPathDoesNotExist {
		notAnalyzedIndex := bleve.NewTextFieldMapping()
		notAnalyzedIndex.Analyzer = "keyword"

		textIndex := bleve.NewTextFieldMapping()
		textIndex.Analyzer = "standard"

		itemMapping := bleve.NewDocumentMapping()
		itemMapping.AddFieldMappingsAt("Int", bleve.NewNumericFieldMapping())
		itemMapping.AddFieldMappingsAt("Ints", notAnalyzedIndex)
		itemMapping.AddFieldMappingsAt("Url", notAnalyzedIndex)
		itemMapping.AddFieldMappingsAt("Title", textIndex)

		mapping := bleve.NewIndexMapping()
		// mapping.DefaultAnalyzer = "keyword"
		mapping.AddDocumentMapping("Item", itemMapping)

		index, err = bleve.New(dbFileName, mapping)

		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 100; i++ {
		index.Index(strconv.Itoa(i), &Item{i, []string{"/url/text/foo" + strconv.Itoa(i) + "/bar"}, []int{i % 2}, "-url-text-foo1-bar " + strconv.Itoa(i)})
	}

	// query := bleve.NewQueryStringQuery("/url/text/foo50/bar")
	// var min, max float64
	// min = 3
	// max = 3
	// isPoint := true
	// query := bleve.NewNumericRangeQuery(&min, &max).SetField("Int")
	query := bleve.NewMatchQuery("1").SetField("Ints")
	search := bleve.NewSearchRequest(query)
	searchResults, _ := index.Search(search)

	println(searchResults.String())
}
*/
