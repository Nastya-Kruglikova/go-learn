package main

import (
    "github.com/spf13/viper"
    "gopkg.in/fsnotify.v1"
    "fmt"
)

func main (){
    viper.SetConfigName("test")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/opt/maytech/etc/")
    viper.ReadInConfig()
    viper.SetDefault("logFile", "/opt/maytech/var/log/test.log")
    fmt.Println(viper.Get("logFile"))
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    })
}