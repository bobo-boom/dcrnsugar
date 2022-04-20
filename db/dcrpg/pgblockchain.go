package dcrpg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/cache"
	"github.com/bobo-boom/dcrnsugar/config"
	"github.com/bobo-boom/dcrnsugar/db/dbtypes"
	"time"
)

type ChainDB struct {
	ctx          context.Context
	queryTimeout time.Duration
	db           *sql.DB
	addressCache *cache.CacheAddress
}

func (c *ChainDB) timeoutError() string {
	return fmt.Sprintf("%s after %v", dbtypes.TimeoutPrefix, c.queryTimeout)

}

func NewChainDB(ctx context.Context, config *config.Config, addrCache *cache.CacheAddress, shutdown func()) (*ChainDB, error) {
	//connect to the PostgreSql  daemon and return the *sql.DB
	db, err := Connect(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		return nil, err
	}
	queryTimeout := time.Hour

	return &ChainDB{
		ctx:          ctx,
		db:           db,
		queryTimeout: queryTimeout,
		addressCache: addrCache,
	}, nil

}
