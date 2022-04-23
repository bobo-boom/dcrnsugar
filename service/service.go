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

	cache            *cache.CacheAddress
	balanceInfoCache *cache.BalanceInfoCache

	AddressCh     chan []*AddressAndId
	BalanceInfoCh chan []*dbtypes.BalanceInfo
	FinishCh      chan bool

	ctx context.Context
}

func NewService(config *config.Config, ctx context.Context) (*Service, error) {

	// new  cache
	addrCache := cache.NewCacheAddress()
	balanceInfoCache := cache.NewBalanceInfoCache()

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
	//new channel
	addrCh := make(chan []*AddressAndId, 10000)
	balanceInfoCh := make(chan []*dbtypes.BalanceInfo, 10000)
	//new service
	s := &Service{
		cdb:              db,
		config:           config,
		cache:            addrCache,
		client:           c,
		balanceInfoCache: balanceInfoCache,
		AddressCh:        addrCh,
		BalanceInfoCh:    balanceInfoCh,
		ctx:              ctx,
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
	if err != nil {
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

func (s *Service) StoreBalanceBatch(balance []*dbtypes.BalanceInfo) error {

	err := s.cdb.InsertAddsBalances(balance)
	if err != nil {
		fmt.Printf("store balance err %v\n", err)
		return err
	}
	return err
}
func (s *Service) GetAddressAsync() error {

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	start, err := s.cdb.RetrieveBestBalanceIndex(ctx)
	if err != nil {
		fmt.Printf("get start err %v\n", err)
		return err
	}
	end, err := s.cdb.RetrieveBestAddressId(ctx)
	if err != nil {
		fmt.Printf("get end err %v\n", err)
		return err
	}
	//获取地址
	fmt.Printf("get address  from  %d  to %d\n", start, end)
	var step int64
	step = 10000
	for i := start; i <= end; {
		if end-i < int64(step) {
			step = end - i
		}
		select {
		case <-s.ctx.Done():
			fmt.Println("GetAddressAsync  exit!!!")
			return nil
		default:
		}
		//address, err := s.cdb.RetrieveAddress(ctx, i)

		addresses, ids, err := s.cdb.RetrieveAddresses(ctx, i, i+step)

		if err != nil {
			fmt.Printf("get  id %d   to %d", ids[0], ids[len(ids)]-1)
			return err
		}
		addressIds := make([]*AddressAndId, 0)
		for index, address := range addresses {
			ad := &AddressAndId{address: address, id: ids[index]}
			addressIds = append(addressIds, ad)
		}
		s.AddressCh <- addressIds
		i += step
	}
	s.FinishCh <- true
	return nil

}
func (s *Service) GetGetBalanceOfAddrs(addrs []*AddressAndId) ([]*dbtypes.BalanceInfo, error) {
	infos := make([]*dbtypes.BalanceInfo, 0)
	for _, addrID := range addrs {
		balance, err := s.GetBalanceOfAddr(addrID.address)
		if err != nil {
			fmt.Printf("get  %s balance err: %v\n", addrID.address, balance)
			return nil, err
		}
		info := &dbtypes.BalanceInfo{Balance: balance, Index: addrID.id, Flag: false, Address: addrID.address}
		infos = append(infos, info)
	}
	fmt.Printf("get balacne %d to %d  success \n", addrs[0].id, addrs[len(addrs)-1].id)

	return infos, nil
}

// 获取余额
func (s *Service) GetBalanceOfAddrAsync() error {
	// get address form ch
	var adds []*AddressAndId
	for {
		select {
		case adds = <-s.AddressCh:
			addsInfo, err := s.GetGetBalanceOfAddrs(adds)
			if err != nil {
				fmt.Printf("get balacne %d to %d  success \n", adds[0].id, adds[len(adds)-1].id)

			}
			s.BalanceInfoCh <- addsInfo
		case <-s.ctx.Done():
			fmt.Println("GetBalanceOfAddrAsync exit!")
			return nil
		}

	}
}

func (s *Service) CommitBalanceInfo() error {
	var bs []*dbtypes.BalanceInfo
	for {
		select {
		case bs = <-s.BalanceInfoCh:
			err := s.StoreBalanceBatch(bs)
			if err != nil {
				fmt.Printf("store balanceInfo from %d  to %d err : %v", bs[0].Index, bs[len(bs)-1].Index, err)
				return err
			}
			balanceIndex := &dbtypes.BalanceIndex{Index: bs[len(bs)-1].Index}
			err = s.cdb.InsertBalanceIndex(balanceIndex)
			if err != nil {
				fmt.Printf("insert balance index err %v\n", err)
				return err
			}
			fmt.Printf("handle id from  %d   to %d success ......\n", bs[0].Index, bs[len(bs)-1].Index)

		case <-s.ctx.Done():
			fmt.Println("CommitBalanceInfo exit!")
			return nil

		}

	}
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
		if s.cache.GetAddressStatus(address) {
			fmt.Printf("%s  already  finished ......\n", address)
			return nil
		}

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
	s.cache.WriteCache(address, true)
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
