package main

import (
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*config)
}
