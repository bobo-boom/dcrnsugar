package dcrpg

import (
	"context"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/cache"
	"github.com/bobo-boom/dcrnsugar/config"
	"github.com/bobo-boom/dcrnsugar/db/dbtypes"
	"testing"
	"time"
)

func TestChainDB_CreateBalanceTable(t *testing.T) {

	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	err = db.CreateBalanceTable()
	if err != nil {
		fmt.Printf("err %v\n", err)
		return

	}
}

func TestChainDB_CreateAddressIndexOfBalanceTable(t *testing.T) {

	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	err = db.CreateAddressIndexOfBalanceTable()
	if err != nil {
		fmt.Printf("err %v\n", err)
		return

	}
}

func TestChainDB_InsertAddsBalance(t *testing.T) {

	cof := &config.DefaultConfig
	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	balanceInfo := &dbtypes.BalanceInfo{
		Balance: 1000000,
		Index:   1000000,
		Flag:    false,
		Address: "dsaihfadf",
	}
	err = db.InsertAddsBalance(balanceInfo)
	if err != nil {
		fmt.Printf("err  : %v\n", err)
		return
	}
}

func TestChainDB_CreateBalanceIndexTable(t *testing.T) {

	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	err = db.CreateBalanceIndexTable()
	if err != nil {
		fmt.Printf("err %v\n", err)
		return

	}
}

func TestChainDB_InsertBalanceIndex(t *testing.T) {

	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	bi := &dbtypes.BalanceIndex{
		Index: 100000,
	}
	err = db.InsertBalanceIndex(bi)
	if err != nil {
		fmt.Printf("err %v\n", err)

	}
}
func TestRetrieveBestBalanceIndex(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}

	index, err := db.RetrieveBestBalanceIndex(ctx)
	err = err
	if err != nil {
		fmt.Printf("err %v\n", err)

	}
	fmt.Println("index ", index)
}
func TestRetrieveAddresses(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cof := &config.DefaultConfig

	db, err := NewChainDB(cof, cache.NewCacheAddress(), nil)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}
	id := 2
	address, err := db.RetrieveAddress(ctx, int64(id))
	if err != nil {
		fmt.Printf("err %v\n", err)
	}
	fmt.Println(address)
}
