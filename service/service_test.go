package service

import (
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"testing"
)

func TestService_GetBalanceOfAddr(t *testing.T) {

	cof := &config.DefaultConfig
	service, err := NewService(cof)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	address := "DsjNv5h3jr1cZs5XF8TrPv2vhCbikwnE51B"
	balance, err := service.GetBalanceOfAddr(address)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	fmt.Println("balance : ", balance)
}

func TestService_HandleAddress(t *testing.T) {
	cof := &config.DefaultConfig
	service, err := NewService(cof)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	err = service.HandleAddress(1000)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
}
