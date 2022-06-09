package config

import "C"
import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "gopkg.in/yaml.v2"
	"log"
	"os"
)

type Urls struct {
	List []string `yaml: "list"`
}

type App struct {
	Time string `yaml:"time"`
}

type Database struct {
	Driver   string `yaml: "driver"`
	Host     string `yaml: "host"`
	Dbname   string `yaml: "dbname"`
	User     string `yaml: "user"`
	Password string `yaml: "password"`
}
type Configuration struct {
	App      App      `yaml: "app"`
	Database Database `yaml: "database"`
	Urls     Urls     `yaml: "urls"`
}

var ViperConfig Configuration

func init() {

	//读取yaml配置文件, 将yaml配置文件，转换struct类型
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip := viper.New()
	vip.AddConfigPath(path + "/etc") //设置读取的文件路径
	vip.SetConfigName("config")      //设置读取的文件名
	vip.SetConfigType("yaml")        //设置文件的类型
	//尝试进行配置读取
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	err = vip.Unmarshal(&ViperConfig)
	log.Printf("解析到的配置:%v", ViperConfig)
	if err != nil {
		panic(err)
	}
	vip.OnConfigChange(func(e fsnotify.Event) {
		vip.Unmarshal(&ViperConfig)
		log.Println("Info: config change---" + C.App.Time)
	})

}
