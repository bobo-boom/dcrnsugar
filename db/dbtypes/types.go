package dbtypes

import "context"

var (
	// PGCancelError is the error string PostgreSQL returns when a query fails
	// to complete due to user requested cancellation.
	PGCancelError       = "pq: canceling statement due to user request"
	CtxDeadlineExceeded = context.DeadlineExceeded.Error()
	TimeoutPrefix       = "TIMEOUT of PostgreSQL query"
)

// BalanceInfo  Balance table info
type BalanceInfo struct {
	Balance int64  `json:"balance"`
	Index   int64  `json:"index"`
	Flag    bool   `json:"flag"`
	Address string `json:"address"`
}

// BalanceIndex
type BalanceIndex struct {
	Index int64 `json:"index"`
}
