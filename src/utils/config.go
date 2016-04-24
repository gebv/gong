package utils

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

//dico struct
//config.toml
// name = "StorageConfig"
//[[fields]]
//name = "DbSearch"
//type = "string"
//tag = '''toml:"db_search"'''

//[[fields]]
//name = "DbStore"
//type = "string"
//tag = '''toml:"db_store"'''

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewStorageConfig() *StorageConfig {
	model := new(StorageConfig)

	return model
}

type StorageConfig struct {
	DbSearch string `toml:"db_search"`

	DbStore string `toml:"db_store"`
}

// SetDbSearch set DbSearch
func (s *StorageConfig) SetDbSearch(v string) {
	s.DbSearch = v
}

// GetDbSearch get DbSearch
func (s *StorageConfig) GetDbSearch() string {
	return s.DbSearch
}

// SetDbStore set DbStore
func (s *StorageConfig) SetDbStore(v string) {
	s.DbStore = v
}

// GetDbStore get DbStore
func (s *StorageConfig) GetDbStore() string {
	return s.DbStore
}

//<<<AUTOGENERATE.DICO

//dico struct
//config.toml
// name = "ApiConfig"

//[[fields]]
//name = "Bind"
//type = "string"
//tag = '''toml:"bind"'''

//[[fields]]
//name = "CookieStoreKey"
//type = "string"
//tag = '''toml:"cookie_store_key"'''

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewApiConfig() *ApiConfig {
	model := new(ApiConfig)

	return model
}

type ApiConfig struct {
	Bind string `toml:"bind"`

	CookieStoreKey string `toml:"cookie_store_key"`
}

// SetBind set Bind
func (a *ApiConfig) SetBind(v string) {
	a.Bind = v
}

// GetBind get Bind
func (a *ApiConfig) GetBind() string {
	return a.Bind
}

// SetCookieStoreKey set CookieStoreKey
func (a *ApiConfig) SetCookieStoreKey(v string) {
	a.CookieStoreKey = v
}

// GetCookieStoreKey get CookieStoreKey
func (a *ApiConfig) GetCookieStoreKey() string {
	return a.CookieStoreKey
}

//<<<AUTOGENERATE.DICO

//dico struct
//config.toml
// name = "Config"

//[[fields]]
//name = "Storage"
//type = "StorageConfig"
//tag = '''toml:"storage"'''

//[[fields]]
//name = "Api"
//type = "ApiConfig"
//tag = '''toml:"api"'''

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  struct

func NewConfig() *Config {
	model := new(Config)

	return model
}

type Config struct {
	Storage StorageConfig `toml:"storage"`

	Api ApiConfig `toml:"api"`
}

// SetStorage set Storage
func (c *Config) SetStorage(v StorageConfig) {
	c.Storage = v
}

// GetStorage get Storage
func (c *Config) GetStorage() StorageConfig {
	return c.Storage
}

// SetApi set Api
func (c *Config) SetApi(v ApiConfig) {
	c.Api = v
}

// GetApi get Api
func (c *Config) GetApi() ApiConfig {
	return c.Api
}

//<<<AUTOGENERATE.DICO

var Cfg Config

func InitConfig(file string) error {
	_, err := toml.DecodeFile(findConfigFile(file), &Cfg)

	return err
}

func findConfigFile(fileName string) string {
	if len(fileName) == 0 {
		panic("Empty file name")
	}

	if _, err := os.Stat("./config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("./config/" + fileName)
	} else if _, err := os.Stat("../config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("../config/" + fileName)
	} else if _, err := os.Stat(fileName); err == nil {
		fileName, _ = filepath.Abs(fileName)
	} else {
		panic("Not found " + fileName)
	}

	return fileName
}
