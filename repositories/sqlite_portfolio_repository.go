package repositories

import (
	"database/sql"
	models "github.com/fcopulgar/stock-manager-go/models"
	"log"
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
	// TODO
	return nil, nil
}

func (repo *SQLitePortfolioRepository) GetByID(id int) (*models.Portfolio, error) {
	// TODO
	return nil, nil
}

func (repo *SQLitePortfolioRepository) Save(portfolio *models.Portfolio) error {
	// TODO
	return nil
}

func (repo *SQLitePortfolioRepository) Update(portfolio *models.Portfolio) error {
	// TODO
	return nil
}

func (repo *SQLitePortfolioRepository) Delete(id int) error {
	// TODO
	return nil
}
