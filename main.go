package main

import (
	"context"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"github.com/bobo-boom/dcrnsugar/service"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("load config err %v\n", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	s, err := service.NewService(conf, ctx)
	if err != nil {
		fmt.Printf("new service errr %v\n", err)
		return
	}

	err = s.PrepareDB(ctx)
	if err != nil {
		fmt.Printf("prepare db err  %v\n", err)
	}

	go func() {
		err1 := s.CommitBalanceInfo()
		if err1 != nil {
			cancel()
			fmt.Printf("err : %v\n", err1)
		}
	}()

	go func() {
		err1 := s.GetBalanceOfAddrAsync()
		if err1 != nil {
			cancel()
			fmt.Printf("err : %v\n", err1)

		}
	}()
	go func() {
		err1 := s.GetAddressAsync()
		if err1 != nil {
			cancel()
			fmt.Printf("err : %v\n", err1)

		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-s.FinishCh:

		for {
			if !(len(s.BalanceInfoCh) == 0 && len(s.AddressCh) == 0) {
				time.Sleep(2 * time.Second)
				continue
			}
			fmt.Println("Do the final processing....")
			end, _ := s.RetrieveBestAddressId(ctx)
			for {
				latest, _ := s.RetrieveBestBalanceIndex(ctx)
				if latest != end {
					time.Sleep(5 * time.Second)
					continue
				}
				fmt.Println("Finished !!!!!!!!!!!!!!!!!!!!!!!!!!")
				break
			}

			close(s.BalanceInfoCh)
			close(s.AddressCh)
			close(s.FinishCh)

			cancel()
			fmt.Println("Bye !")
			return

		}
	}

}
