package main

import (
	"bufio"
	"fmt"
	"github.com/fcopulgar/stock-manager-go/config"
	"github.com/fcopulgar/stock-manager-go/repositories"
	"github.com/fcopulgar/stock-manager-go/services"
	"os"
	"strconv"
	"strings"
)

func main() {
	config.LoadConfig()

	repo := repositories.NewSQLitePortfolioRepository("portfolios.db")
	apiKey := config.GetEnv("ALPHAVANTAGE_API_KEY")
	stockService := services.NewAlphaVantageService(apiKey)
	portfolioService := services.NewPortfolioService(repo, stockService)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Select an option:")
		fmt.Println("1. Show portfolios")
		fmt.Println("2. Create portfolio manually")
		fmt.Println("3. Create random portfolio")
		fmt.Println("4. Quit")

		fmt.Print("Option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			showPortfolios(portfolioService, reader)
		case "2":
			createPortfolio(portfolioService, reader)
		case "3":
			crearteRandomPortfolio(portfolioService)
		case "4":
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func showPortfolios(ps *services.PortfolioService, reader *bufio.Reader) {
	portfolios, err := ps.Repo.GetAll()
	if err != nil {
		fmt.Printf("Error in obtaining portfolios: %v\n", err)
		return
	}

	if len(portfolios) == 0 {
		fmt.Println("No portfolios registered.")
		return
	}

	for _, p := range portfolios {
		fmt.Printf("ID: %d, Name: %s\n", p.ID, p.Name)
	}

	fmt.Print("Enter the ID of the portfolio you wish to view (or press Enter to return): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("ID invalid.")
		return
	}

	portfolio, err := ps.Repo.GetByID(id)
	if err != nil {
		fmt.Printf("Error in obtaining the portfolio: %v\n", err)
		return
	}

	showPortfolioDetails(ps, portfolio, reader)
}
