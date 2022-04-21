package main

import (
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"github.com/bobo-boom/dcrnsugar/service"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("load config err %v\n", err)
		return
	}
	s, err := service.NewService(conf)
	if err != nil {
		fmt.Printf("new service errr %v\n", err)
		return
	}
	s.Start()

}
