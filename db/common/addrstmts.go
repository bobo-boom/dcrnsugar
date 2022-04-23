package common

const (
	BalanceTableName = "balance"

	BalanceIndexTableName = "balanceindex"

	AddressTableName = "addresses"

	CreateBalanceTable = `CREATE TABLE IF NOT EXISTS balance (
		id SERIAL8 PRIMARY KEY,
		address TEXT ,
		balance INT8 ,
		index   INT8 ,
		flag    BOOLEAN

	);`
	insertBalanceRow = `INSERT INTO balance (address,balance,index,flag) VALUES ($1, $2, $3 ,$4 )`

	InsertBalanceRow = insertBalanceRow + `RETURNING id;`

	UpsertBalanceRow = insertBalanceRow + ` ON CONFLICT (address) DO UPDATE  SET balance = $5 RETURNING id;`

	IndexOfBalanceTableOnAddress = "address_index"

	IndexBalanceTableOnAddress = `CREATE UNIQUE INDEX IF NOT EXISTS ` + IndexOfBalanceTableOnAddress +
		` ON balance(address);`

	//BalanceIdex

	CreateBalanceIndexTable = `CREATE TABLE IF NOT EXISTS balanceindex (
		id SERIAL8 PRIMARY KEY,
        index  INT8
	);`
	InitBalanceIndexRow = `INSERT INTO balanceindex (index) VALUES (1)`

	BalanceIndexCountRow = `SELECT count(index) from balanceindex`

	insertBalanceIndexRow = `INSERT INTO balanceindex (index) VALUES ($1 )`

	InsertBalanceIndexRow = insertBalanceIndexRow + `RETURNING id;`

	SelectBalanceIndexBestRow = `SELECT  index FROM balanceindex 
			ORDER BY id desc  limit 1;`

	SelectAddressRow = `SELECT  address FROM ` + AddressTableName +
		` Where id=$1;`

	SelectAddressRows = `SELECT id, address FROM ` + AddressTableName +
		` Where id >= $1 AND id < $2;`

	SelectBestAddressIdRow = `SELECT  id FROM ` + AddressTableName +
		` ORDER BY id desc limit 1;`
)
