package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/cache"
	"github.com/bobo-boom/dcrnsugar/client"
	"github.com/bobo-boom/dcrnsugar/config"
	"github.com/bobo-boom/dcrnsugar/db/dbtypes"
	"github.com/bobo-boom/dcrnsugar/db/dcrpg"
	"net/http"
	"time"
)

// Service
type Service struct {
	cdb    *dcrpg.ChainDB
	client *client.Client
	config *config.Config
	cache  *cache.CacheAddress
}

func NewService(config *config.Config) (*Service, error) {

	// new  cache
	addrCache := cache.NewCacheAddress()

	// new chainDB
	db, err := dcrpg.NewChainDB(config, addrCache, nil)
	if err != nil {
		fmt.Printf("new chainDB err : %v", err)
		return nil, err
	}
	//new client
	c, err := client.New(config)
	if err != nil {
		fmt.Printf("new client err : %v\n", err)
		return nil, err
	}
	// new service
	s := &Service{
		cdb:    db,
		config: config,
		cache:  addrCache,
		client: c,
	}
	return s, nil
}
func (s *Service) PrepareDB(ctx context.Context) error {

	// create address table

	err := s.cdb.CreateBalanceTable()
	if err != nil {
		fmt.Printf("PrepareDB : create balance table err : %v\n", err)
		return err
	}
	err = s.cdb.CreateAddressIndexOfBalanceTable()
	if err!=nil{
		fmt.Printf("PrepareDB :CreateAddressIndexOfBalanceTable err : %v\n", err)
		return err
	}

	// create addressindex table
	err = s.cdb.CreateBalanceIndexTable()
	if err != nil {
		fmt.Printf("PrepareDB : create balanceindex table err : %v\n", err)
		return err
	}

	indexCount, err := s.cdb.RetrieveBalanceIndexCount(ctx)

	if err != nil {
		fmt.Printf("retrieve balance index count err%v\n", err)
	}

	if indexCount == 0 {
		err := s.cdb.InitBalanceIndexTable()
		if err != nil {
			fmt.Printf("init balance index table err %v \n", err)
		}
		fmt.Printf("init balance index table ......\n")

	}
	return nil

}
func (s *Service) GetBalanceOfAddr(address string) (int64, error) {

	url := s.client.GenerateGetBalanceUrl(address)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("get balance err :  %v\n", err)
		return 0, err
	}
	details := &client.RequestDetails{
		HttpRequest: request,
	}
	repbytes, err := s.client.SendRequest(details)
	if err != nil {
		fmt.Printf("get balance err :  %v\n", err)
		return 0, err
	}
	addrInfo := &InsightAddressInfo{}

	err = json.Unmarshal(repbytes, addrInfo)
	if err != nil {
		fmt.Printf("get balance err : %v\n", err)
		return 0, err
	}
	balance := addrInfo.BalanceSat
	return balance, nil

}

func (s *Service) StoreBalance(balance *dbtypes.BalanceInfo) error {

	err := s.cdb.InsertAddsBalance(balance)
	if err != nil {
		fmt.Printf("store balance err %v\n", err)
		return err
	}
	return err
}
func (s *Service) HandleAddress(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TimeOut)*time.Second)
	defer cancel()

	//1、获取地址
	address, err := s.cdb.RetrieveAddress(ctx, id)

	if err != nil {
		fmt.Printf("get  id: %d  address: %s\n", id, address)
		return err
	}
	//2、查看缓存
	if s.cache.IsExist(address) {
		fmt.Printf("%s  already  finished ......\n", address)
		return nil
	}
	//3、获取余额
	balance, err := s.GetBalanceOfAddr(address)
	if err != nil {
		fmt.Printf("get  %s balance err: %v\n", address, balance)
		return err
	}
	//4、入库
	info := &dbtypes.BalanceInfo{
		Balance: balance,
		Index:   id,
		Flag:    false,
		Address: address,
	}
	err = s.StoreBalance(info)
	if err != nil {
		fmt.Printf("store balanceInfo of %s  err : %v", address, err)
		return err
	}
	//5、记录缓存
	s.cache.WriteCache(address)
	//6、记录进度
	balanceInfo := &dbtypes.BalanceIndex{Index: id}
	err = s.cdb.InsertBalanceIndex(balanceInfo)
	if err != nil {
		fmt.Printf("insert balance index err %v\n", err)
		return err
	}
	fmt.Printf("handle id : %d    %s success ......\n", id, address)

	return nil
}
func (s *Service) HandleAddresses(start, end int64) error {
	for i := start; i < end; i++ {
		err := s.HandleAddress(i)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TimeOut)*time.Second)
	defer cancel()

	err := s.PrepareDB(ctx)
	if err != nil {
		fmt.Printf("prepare db err %v\n", err)
		return
	}
	start, err := s.cdb.RetrieveBestBalanceIndex(ctx)
	if err != nil {
		fmt.Printf("get start err %v\n", err)
		return
	}
	end, err := s.cdb.RetrieveBestAddressId(ctx)
	if err != nil {
		fmt.Printf("get end err %v\n", err)
		return
	}
	fmt.Printf("start handling from %d to %d addresses .....\n", start, end)
	err = s.HandleAddresses(start, end)
	if err != nil {
		fmt.Printf("handle addresses err %v\n", err)
		return
	}
	fmt.Printf("work finished !!!!!!!")
	return
}
