package common

const (
	BalanceTableName = "balanceasync"

	BalanceIndexTableName = "balanceindexasync"

	AddressTableName = "addresses"

	CreateBalanceTable = `CREATE TABLE IF NOT EXISTS ` + BalanceTableName + ` (
		id SERIAL8 PRIMARY KEY,
		address TEXT ,
		balance INT8 ,
		index   INT8 ,
		flag    BOOLEAN

	);`
	insertBalanceRow = `INSERT INTO ` + BalanceTableName + ` (address,balance,index,flag) VALUES ($1, $2, $3 ,$4 )`

	InsertBalanceRow = insertBalanceRow + `RETURNING id;`

	UpsertBalanceRow = insertBalanceRow + ` ON CONFLICT (address) DO UPDATE  SET balance = $5 RETURNING id;`

	IndexOfBalanceTableOnAddress = "address_index_async"

	IndexBalanceTableOnAddress = `CREATE UNIQUE INDEX IF NOT EXISTS ` + IndexOfBalanceTableOnAddress +
		` ON ` + BalanceTableName + `(address);`

	//BalanceIdex

	CreateBalanceIndexTable = `CREATE TABLE IF NOT EXISTS ` + BalanceIndexTableName + ` (
		id SERIAL8 PRIMARY KEY,
        index  INT8
	);`
	InitBalanceIndexRow = `INSERT INTO ` + BalanceIndexTableName + `  (index) VALUES (1)`

	BalanceIndexCountRow = `SELECT count(index) from ` + BalanceIndexTableName

	insertBalanceIndexRow = `INSERT INTO ` + BalanceIndexTableName + `  (index) VALUES ($1 )`

	InsertBalanceIndexRow = insertBalanceIndexRow + `RETURNING id;`

	SelectBalanceIndexBestRow = `SELECT  index FROM ` + BalanceIndexTableName + `  
			ORDER BY id desc  limit 1;`

	SelectAddressRow = `SELECT  address FROM ` + AddressTableName +
		` Where id=$1;`

	SelectAddressRows = `SELECT id, address FROM ` + AddressTableName +
		` Where id >= $1 AND id < $2;`

	SelectBestAddressIdRow = `SELECT  id FROM ` + AddressTableName +
		` ORDER BY id desc limit 1;`
)
