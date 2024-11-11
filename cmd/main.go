package main

import (
	"github.com/fcopulgar/stock-manager-go/cmd/cli"
	"github.com/fcopulgar/stock-manager-go/config"
	"github.com/fcopulgar/stock-manager-go/repositories"
	"github.com/fcopulgar/stock-manager-go/services"
)

func main() {
	config.LoadConfig()

	// Initialize the repository and services
	repo := repositories.NewSQLitePortfolioRepository("portfolios.db")
	apiKey := config.GetEnv("ALPHAVANTAGE_API_KEY")
	stockService := services.NewAlphaVantageService(apiKey)
	portfolioService := services.NewPortfolioService(repo, stockService)

	// Start the CLI
	cli := cli.NewCLI(portfolioService)
	cli.Run()
}
