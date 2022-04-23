package service

// InsightAddressInfo models basic information about an address.
type InsightAddressInfo struct {
	Address                  string   `json:"addrStr,omitempty"`
	Balance                  float64  `json:"balance"`
	BalanceSat               int64    `json:"balanceSat"`
	TotalReceived            float64  `json:"totalReceived"`
	TotalReceivedSat         int64    `json:"totalReceivedSat"`
	TotalSent                float64  `json:"totalSent"`
	TotalSentSat             int64    `json:"totalSentSat"`
	UnconfirmedBalance       float64  `json:"unconfirmedBalance"`
	UnconfirmedBalanceSat    int64    `json:"unconfirmedBalanceSat"`
	UnconfirmedTxAppearances int64    `json:"unconfirmedTxApperances"` // [sic]
	TxAppearances            int64    `json:"txApperances"`            // [sic]
	TransactionsID           []string `json:"transactions,omitempty"`
}

type AddressAndId struct {
	id      int64
	address string
}
