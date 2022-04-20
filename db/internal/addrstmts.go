package internal

const (
	CreateBalanceTable = `CREATE TABLE IF NOT EXISTS balance (
		address TEXT PRIMARY KEY,
		id INT8 ,
		value INT8,
	);`

	InsertBalanceRow = `INSERT INTO balance (address,balance,index,flag) VALUES ($1, $2, $3 ,$4 );`

	UpsertBalanceRow = InsertBalanceRow + `ON CONFLICT (address)
		DO UPDATE SET balance = $5;`

	SelectAddressRow = `SELECT id, address FROM addresses 
			Where address=$1;`

	SelectAddressRows = `SELECT id, address FROM addresses 
			Where id >= $1 AND id < $2;`
)
