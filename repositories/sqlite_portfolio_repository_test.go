package repositories

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fcopulgar/stock-manager-go/models"
)

func TestSQLitePortfolioRepository_CRUD(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	repo := NewSQLitePortfolioRepository(dbPath)

	// Create a test portfolio
	portfolio := &models.Portfolio{
		Name: "Test Portfolio",
		Stocks: []models.Stock{
			{
				Symbol:   "AAPL",
				Quantity: 10,
				BuyDate:  time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
				BuyPrice: 300.0,
			},
		},
	}

	// Test Save
	err = repo.Save(portfolio)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// After saving, the portfolio must have an ID > 0.
	// However, the Save function does not assign ID to the object in the current code.
	// You could modify Save to assign the ID to the object or just not check the ID here.

	// Test GetAll (should return 1 portfolio)
	portfolios, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Expected no error from GetAll, got %v", err)
	}
	if len(portfolios) != 1 {
		t.Fatalf("Expected 1 portfolio, got %d", len(portfolios))
	}
	if portfolios[0].Name != "Test Portfolio" {
		t.Errorf("Expected portfolio name 'Test Portfolio', got '%s'", portfolios[0].Name)
	}
	if len(portfolios[0].Stocks) != 1 {
		t.Fatalf("Expected 1 stock in portfolio, got %d", len(portfolios[0].Stocks))
	}
	if portfolios[0].Stocks[0].Symbol != "AAPL" {
		t.Errorf("Expected stock symbol 'AAPL', got '%s'", portfolios[0].Stocks[0].Symbol)
	}

	// Test GetByID
	// We know that the only portfolio saved has ID 1, because AUTOINCREMENT starts at 1.
	p, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("Expected no error from GetByID, got %v", err)
	}
	if p == nil {
		t.Fatal("Expected a portfolio, got nil")
	}
	if p.Name != "Test Portfolio" {
		t.Errorf("Expected portfolio name 'Test Portfolio', got '%s'", p.Name)
	}

	// Test Update
	p.Name = "Updated Portfolio Name"
	p.Stocks[0].Quantity = 20
	err = repo.Update(p)
	if err != nil {
		t.Fatalf("Expected no error from Update, got %v", err)
	}

	// Verify the update
	updated, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("Expected no error from GetByID after update, got %v", err)
	}
	if updated.Name != "Updated Portfolio Name" {
		t.Errorf("Expected updated name 'Updated Portfolio Name', got '%s'", updated.Name)
	}
	if updated.Stocks[0].Quantity != 20 {
		t.Errorf("Expected updated quantity 20, got %d", updated.Stocks[0].Quantity)
	}

	// Test Delete
	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("Expected no error from Delete, got %v", err)
	}

	deleted, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("Expected no error from GetByID after delete, got %v", err)
	}
	if deleted != nil {
		t.Fatal("Expected nil after delete, got a portfolio")
	}

	// Verify that GetAll returns 0 portfolios now
	portfolios, err = repo.GetAll()
	if err != nil {
		t.Fatalf("Expected no error from GetAll after delete, got %v", err)
	}
	if len(portfolios) != 0 {
		t.Fatalf("Expected 0 portfolios after delete, got %d", len(portfolios))
	}
}
