package store

import "github.com/blevesearch/bleve"
import "github.com/golang/glog"

var itemsPerPage = 2

// SearchPerPage поиск записей (по умолчанию записи отсортированы по мере их создания)
func SearchPerPage(filter *SearchFileter) *SearchResult {
	result := NewSearchResult()

	conjuncts := []bleve.Query{}

	for _, collection := range filter.GetCollections() {
		conjuncts = append(conjuncts, bleve.NewTermQuery(collection).SetField("Collections"))
	}

	if len(filter.GetExtId()) != 0 {
		conjuncts = append(conjuncts, bleve.NewTermQuery(filter.GetExtId()).SetField("ExtId"))
	}

	if filter.HasEnabled {
		conjuncts = append(conjuncts, bleve.NewBoolFieldQuery(false).SetField("IsRemoved"))
	}

	if len(filter.GetQuery()) > 0 {
		conjuncts = append(conjuncts, bleve.NewQueryStringQuery(filter.GetQuery()))
	}

	query := bleve.NewConjunctionQuery(conjuncts)
	search := bleve.NewSearchRequest(query)
	search.Size = itemsPerPage
	search.From = itemsPerPage * filter.Page

	searchResults, err := Search.Search(search)

	if err != nil {
		return result
	}

	for _, hit := range searchResults.Hits {
		file := NewFile()

		if err := findOne(hit.ID, file); err != nil {
			glog.Warningf("search: find by item=%v, err=%v", hit.ID, err)
			continue
		}

		result.AddItems(file)
	}

	result.SetTotal(int(searchResults.Total))

	if searchResults.Total > uint64(search.From+len(result.Items)) {
		result.NextPage = filter.Page + 1
		result.HasNext = true
	}

	return result
}

//dico struct
//config.toml
//name ="SearchFileter"
//[[fields]]
//name = "Query"
//type = "string"

//[[fields]]
//name = "Page"
//type = "int"

//[[fields]]
//name = "ExtId"
//type = "string"

//[[fields]]
//name = "HasEnabled"
//type = "bool"

//[[fields]]
//name = "Size"
//type = "int"

//[[fields]]
//name = "Collections"
//type = "[]string"

//[[fields]]
//name = "Params"
//type = "map[interface{}]interface{}"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewSearchFileter() *SearchFileter {
	model := new(SearchFileter)

	model.Params = make(map[interface{}]interface{})

	return model
}

type SearchFileter struct {
	Query string

	Page int

	ExtId string

	HasEnabled bool

	Size int

	Collections []string

	Params map[interface{}]interface{}
}

// SetQuery set Query
func (s *SearchFileter) SetQuery(v string) {
	s.Query = v
}

// GetQuery get Query
func (s *SearchFileter) GetQuery() string {
	return s.Query
}

// SetPage set Page
func (s *SearchFileter) SetPage(v int) {
	s.Page = v
}

// GetPage get Page
func (s *SearchFileter) GetPage() int {
	return s.Page
}

// SetExtId set ExtId
func (s *SearchFileter) SetExtId(v string) {
	s.ExtId = v
}

// GetExtId get ExtId
func (s *SearchFileter) GetExtId() string {
	return s.ExtId
}

// SetHasEnabled set HasEnabled
func (s *SearchFileter) SetHasEnabled(v bool) {
	s.HasEnabled = v
}

// GetHasEnabled get HasEnabled
func (s *SearchFileter) GetHasEnabled() bool {
	return s.HasEnabled
}

// SetSize set Size
func (s *SearchFileter) SetSize(v int) {
	s.Size = v
}

// GetSize get Size
func (s *SearchFileter) GetSize() int {
	return s.Size
}

// AddCollections add element Collections
func (s *SearchFileter) AddCollections(v string) {
	if s.IncludeCollections(v) {
		return
	}

	s.Collections = append(s.Collections, v)
}

// RemoveCollections remove element Collections
func (s *SearchFileter) RemoveCollections(v string) {
	if !s.IncludeCollections(v) {
		return
	}

	_i := s.IndexCollections(v)

	s.Collections = append(s.Collections[:_i], s.Collections[_i+1:]...)
}

// GetCollections get Collections
func (s *SearchFileter) GetCollections() []string {
	return s.Collections
}

// IndexCollections get index element Collections
func (s *SearchFileter) IndexCollections(v string) int {
	for _index, _v := range s.Collections {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeCollections has exist value Collections
func (s *SearchFileter) IncludeCollections(v string) bool {
	return s.IndexCollections(v) > -1
}

// SetParams set all elements Params
func (s *SearchFileter) SetParams(v map[interface{}]interface{}) {
	s.Params = make(map[interface{}]interface{})

	for key, value := range v {
		s.Params[key] = value
	}
}

// AddParams add element by key
func (s *SearchFileter) SetOneParams(k interface{}, v interface{}) {
	s.Params[k] = v
}

// RemoveParams remove element by key
func (s *SearchFileter) RemoveParams(k interface{}) {
	if _, exist := s.Params[k]; exist {
		delete(s.Params, k)
	}
}

// GetParams get Params
func (s *SearchFileter) GetParams() map[interface{}]interface{} {
	return s.Params
}

// ExistParams has exist key Params
func (s *SearchFileter) ExistKeyParams(k interface{}) bool {
	_, exist := s.Params[k]

	return exist
}

func (s *SearchFileter) GetOneParams(k interface{}) interface{} {
	return s.Params[k]
}

//<<<AUTOGENERATE.DICO

//dico struct
//config.toml
//name ="SearchResult"
//[[fields]]
//name="Items"
//type="[]*File"

//[[fields]]
//name="Total"
//type="int"

//[[fields]]
//name="HasNext"
//type="bool"

//[[fields]]
//name="NextPage"
//type="int"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewSearchResult() *SearchResult {
	model := new(SearchResult)

	return model
}

type SearchResult struct {
	Items []*File

	Total int

	HasNext bool

	NextPage int
}

// AddItems add element Items
func (s *SearchResult) AddItems(v *File) {
	if s.IncludeItems(v) {
		return
	}

	s.Items = append(s.Items, v)
}

// RemoveItems remove element Items
func (s *SearchResult) RemoveItems(v *File) {
	if !s.IncludeItems(v) {
		return
	}

	_i := s.IndexItems(v)

	s.Items = append(s.Items[:_i], s.Items[_i+1:]...)
}

// GetItems get Items
func (s *SearchResult) GetItems() []*File {
	return s.Items
}

// IndexItems get index element Items
func (s *SearchResult) IndexItems(v *File) int {
	for _index, _v := range s.Items {
		if _v == v {
			return _index
		}
	}
	return -1
}

// IncludeItems has exist value Items
func (s *SearchResult) IncludeItems(v *File) bool {
	return s.IndexItems(v) > -1
}

// SetTotal set Total
func (s *SearchResult) SetTotal(v int) {
	s.Total = v
}

// GetTotal get Total
func (s *SearchResult) GetTotal() int {
	return s.Total
}

// SetHasNext set HasNext
func (s *SearchResult) SetHasNext(v bool) {
	s.HasNext = v
}

// GetHasNext get HasNext
func (s *SearchResult) GetHasNext() bool {
	return s.HasNext
}

// SetNextPage set NextPage
func (s *SearchResult) SetNextPage(v int) {
	s.NextPage = v
}

// GetNextPage get NextPage
func (s *SearchResult) GetNextPage() int {
	return s.NextPage
}

//<<<AUTOGENERATE.DICO
