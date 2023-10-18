package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"blockchain_trade/internal/models"
)

const (
	insertSQL = `
INSERT INTO tickers (
	symbol, price_24h, volume_24h, last_trade_price
) VALUES (
	?, ?, ?, ?
)
`

	schemaSQL = `
CREATE TABLE IF NOT EXISTS tickers (
    symbol VARCHAR(8),
    price_24h FLOAT,
    volume_24h FLOAT,
	last_trade_price FLOAT
);

CREATE INDEX IF NOT EXISTS tickers_symbol ON tickers(symbol);
`

	selectSQL = `
SELECT symbol, price_24h, volume_24h, last_trade_price from tickers
`
)

type TickersDBModel struct {
	DB *sql.DB
}

func (tickersDBModel *TickersDBModel) Init() error {
	if _, err := tickersDBModel.DB.Exec(schemaSQL); err != nil {
		return err
	}

	tickers := map[string]models.Ticker{
		"BTC-USD": {"BTC-USD", 0, 0, 0},
		"ETH-USD": {"ETH-USD", 0, 0, 0},
		"BTC-TRY": {"BTC-TRY", 0, 0, 0},
		"ETH-TRY": {"ETH-TRY", 0, 0, 0},
		"BTC-GBP": {"BTC-GBP", 0, 0, 0},
		"ETH-GBP": {"ETH-GBP", 0, 0, 0},
		"BTC-EUR": {"BTC-EUR", 0, 0, 0},
		"ETH-EUR": {"ETH-EUR", 0, 0, 0},
		"BTC-USDT": {"BTC-USDT", 0, 0, 0},
		"ETH-USDT": {"ETH-USDT", 0, 0, 0},
	}

	if err := tickersDBModel.Flush(tickers); err != nil {
		fmt.Println("ERROR: insert - ", err)
		return err
	}

	return nil
}

func (tickersDBModel *TickersDBModel) Flush(tickers map[string]models.Ticker) error {
	stmt, err := tickersDBModel.DB.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := tickersDBModel.DB.Begin()
	if err != nil {
		return err
	}

	for _, ticker := range tickers {
		_, err = tx.Stmt(stmt).Exec(
			ticker.Symbol, ticker.Price24h, ticker.Volume24h, ticker.LastTradePrice)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return stmt.Close()
}

func (tickersDBModel *TickersDBModel) GetAll() (map[string]models.Ticker, error) {
	var tickers map[string]models.Ticker
	tickers = make(map[string]models.Ticker)

	rows, err := tickersDBModel.DB.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ticker := models.Ticker{}

		err := rows.Scan(
			&ticker.Symbol, 
			&ticker.Price24h, 
			&ticker.Volume24h, 
			&ticker.LastTradePrice)
		if err != nil {
			return nil, err
		}

		tickers[ticker.Symbol] = ticker
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickers, rows.Close()
}

func (tickersDBModel *TickersDBModel) Close() error {
	tickersDBModel.DB.Close()

	return nil
}