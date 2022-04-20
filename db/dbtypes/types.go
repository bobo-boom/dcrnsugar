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
	Balance int64
	Index   int64
	Flag    bool
	Address string
}
