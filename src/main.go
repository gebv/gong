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