package service

import (
	"context"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"sync"
	"testing"
)

func TestService_GetBalanceOfAddr(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	cof := &config.DefaultConfig
	service, err := NewService(cof, ctx)
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
	ctx, _ := context.WithCancel(context.Background())

	service, err := NewService(cof, ctx)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	err = service.HandleAddress(1000)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
}

func TestService_GetAddressAsync(t *testing.T) {
	cof := &config.DefaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := NewService(cof, ctx)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	go func() {
		err2 := service.GetAddressAsync()
		if err2 != nil {
			fmt.Println(err2)
		}
	}()
	//var ad *AddressAndId
	//for {
	//	select {
	//	case ad = <-service.AddressCh:
	//		fmt.Println(*ad)
	//	}
	//}

}
func TestService_GetBalanceOfAddrAsync(t *testing.T) {
	cof := &config.DefaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, err := NewService(cof, ctx)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}

	go func() {
		err2 := service.GetBalanceOfAddrAsync()
		if err2 != nil {
			fmt.Println(err2)
		}
	}()
	//for {
	//	addid := &AddressAndId{
	//		id:      1,
	//		address: "DsjNv5h3jr1cZs5XF8TrPv2vhCbikwnE51B",
	//	}
	//	service.AddressCh <- addid
	//}
}
func TestService_CommitBalanceInfo(t *testing.T) {
	cof := &config.DefaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, err := NewService(cof, ctx)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}

	go func() {
		err2 := service.CommitBalanceInfo()
		if err2 != nil {
			fmt.Println(err2)
		}
	}()

	//for {
	//	balance := &dbtypes.BalanceInfo{
	//		Index:   1,
	//		Balance: 12,
	//	}
	//	service.BalanceInfoCh <- balance
	//}
}
func TestBalanceWorkers(t *testing.T) {

	cof := &config.DefaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, err := NewService(cof, ctx)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	ids := make([]*AddressAndId, 0)
	for i := 0; i < 10; i++ {
		a := &AddressAndId{
			id:      int64(i),
			address: "DsjNv5h3jr1cZs5XF8TrPv2vhCbikwnE51B",
		}
		ids = append(ids, a)

	}
	var wg sync.WaitGroup
	wg.Add(5)
	err, infos := BalanceWorkers(service, ids, 5, &wg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(infos)
	select {}
}
