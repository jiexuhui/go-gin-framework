package core

import (
	"fmt"
	"livefun/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config string

func init() {
	v := viper.New()
	v.SetConfigFile("app.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.LF_CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&global.LF_CONFIG); err != nil {
		fmt.Println(err)
	}

}
