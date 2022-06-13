package main

import (
	"github.com/robfig/cron"
	. "grap-data/config"
	"grap-data/handler"
	"log"
)

func main() {

	c := cron.New()
	spec := ViperConfig.App.Time
	c.AddFunc(spec, func() {
		log.Println("执行业务代码...")
		handler.Deal()
	})
	c.Start()
	select {}

}
