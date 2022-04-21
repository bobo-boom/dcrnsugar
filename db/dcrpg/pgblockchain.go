package dcrpg

import (
	"context"
	"database/sql"
	"github.com/bobo-boom/dcrnsugar/cache"
	"github.com/bobo-boom/dcrnsugar/config"
)

type ChainDB struct {
	ctx          context.Context
	db           *sql.DB
	addressCache *cache.CacheAddress
}

func NewChainDB(config *config.Config, addrCache *cache.CacheAddress, shutdown func()) (*ChainDB, error) {
	//connect to the PostgreSql  daemon and return the *sql.DB
	db, err := Connect(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		return nil, err
	}

	return &ChainDB{
		db:           db,
		addressCache: addrCache,
	}, nil

}
