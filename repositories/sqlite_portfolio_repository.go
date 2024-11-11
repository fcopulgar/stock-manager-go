package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/fcopulgar/stock-manager-go/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLitePortfolioRepository struct {
	DB *sql.DB
}

func NewSQLitePortfolioRepository(dbPath string) *SQLitePortfolioRepository {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}

	repo := &SQLitePortfolioRepository{DB: db}
	repo.createTables()
	return repo
}

func (repo *SQLitePortfolioRepository) createTables() {
	portfolioTable := `CREATE TABLE IF NOT EXISTS portfolios (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`

	stockTable := `CREATE TABLE IF NOT EXISTS stocks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        portfolio_id INTEGER,
        symbol TEXT NOT NULL,
        quantity INTEGER NOT NULL,
        buy_date TEXT NOT NULL,
        FOREIGN KEY(portfolio_id) REFERENCES portfolios(id)
    );`

	_, err := repo.DB.Exec(portfolioTable)
	if err != nil {
		log.Fatalf("Error creating the portfolios table: %v", err)
	}

	_, err = repo.DB.Exec(stockTable)
	if err != nil {
		log.Fatalf("Error when creating the stocks table: %v", err)
	}
}

func (repo *SQLitePortfolioRepository) GetAll() ([]models.Portfolio, error) {
	portfolios := []models.Portfolio{}

	rows, err := repo.DB.Query("SELECT id, name FROM portfolios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var portfolio models.Portfolio
		err := rows.Scan(&portfolio.ID, &portfolio.Name)
		if err != nil {
			return nil, err
		}

		stocks, err := repo.getStocksByPortfolioID(portfolio.ID)
		if err != nil {
			return nil, err
		}
		portfolio.Stocks = stocks

		portfolios = append(portfolios, portfolio)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (repo *SQLitePortfolioRepository) GetByID(id int) (*models.Portfolio, error) {
	var portfolio models.Portfolio

	err := repo.DB.QueryRow("SELECT id, name FROM portfolios WHERE id = ?", id).Scan(&portfolio.ID, &portfolio.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Portfolio not found
		}
		return nil, err
	}

	stocks, err := repo.getStocksByPortfolioID(portfolio.ID)
	if err != nil {
		return nil, err
	}
	portfolio.Stocks = stocks

	return &portfolio, nil
}

func (repo *SQLitePortfolioRepository) Save(portfolio *models.Portfolio) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO portfolios (name) VALUES (?)", portfolio.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	portfolioID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	portfolio.ID = int(portfolioID)

	for _, stock := range portfolio.Stocks {
		_, err = tx.Exec(
			"INSERT INTO stocks (portfolio_id, symbol, quantity, buy_date) VALUES (?, ?, ?, ?)",
			portfolioID, stock.Symbol, stock.Quantity, stock.BuyDate.Format("2006-01-02"),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *SQLitePortfolioRepository) Update(portfolio *models.Portfolio) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE portfolios SET name = ? WHERE id = ?", portfolio.Name, portfolio.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM stocks WHERE portfolio_id = ?", portfolio.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, stock := range portfolio.Stocks {
		_, err = tx.Exec(
			"INSERT INTO stocks (portfolio_id, symbol, quantity, buy_date) VALUES (?, ?, ?, ?)",
			portfolio.ID, stock.Symbol, stock.Quantity, stock.BuyDate.Format("2006-01-02"),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *SQLitePortfolioRepository) Delete(id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM stocks WHERE portfolio_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM portfolios WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *SQLitePortfolioRepository) getStocksByPortfolioID(portfolioID int) ([]models.Stock, error) {
	stocks := []models.Stock{}

	rows, err := repo.DB.Query(
		"SELECT symbol, quantity, buy_date FROM stocks WHERE portfolio_id = ?",
		portfolioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		var buyDateStr string

		err := rows.Scan(&stock.Symbol, &stock.Quantity, &buyDateStr)
		if err != nil {
			return nil, err
		}

		stock.BuyDate, err = time.Parse("2006-01-02", buyDateStr)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}
