package main

import (
	"govmpools/tools/conf"
	"log"
)

func main() {
	var cfg conf.Config
	err := cfg.Init("/opt/GoVMPools/conf/config.ini")
	if err != nil {
		log.Fatalln(err)
	}

}
