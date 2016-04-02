package utils

import (
    "github.com/BurntSushi/toml"
    "os"
    "path/filepath"
)

type StorageConfig struct {
    DbSearch string `toml:"db_search"`
    DbStore string `toml:"db_store"`
}

type ApiConfig struct {
    Bind string
}

type Config struct {
    Storage StorageConfig
    Api ApiConfig
}

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