package cli

import (
	"fmt"
	"github.com/fcopulgar/stock-manager-go/api"
	"github.com/fcopulgar/stock-manager-go/models"
	"strconv"
	"strings"
	"time"
)

func (cli *CLI) viewPortfolios() {
	portfolios, err := cli.portfolioService.Repo.GetAll()
	if err != nil {
		fmt.Printf("Error retrieving portfolios: %v\n", err)
		return
	}

	if len(portfolios) == 0 {
		fmt.Println("No portfolios available.")
		return
	}

	for _, p := range portfolios {
		fmt.Printf("ID: %d, Name: %s\n", p.ID, p.Name)
		fmt.Println("Stocks:")
		for _, stock := range p.Stocks {
			fmt.Printf("- %s: %d shares bought on %s\n", stock.Symbol, stock.Quantity, stock.BuyDate.Format("2006-01-02"))
		}

		// Calculate APR
		startDate := earliestBuyDate(p)
		endDate := time.Now()
		apr, err := cli.portfolioService.CalculateAPR(&p, startDate, endDate)
		if err != nil {
			fmt.Printf("Error calculating APR: %v\n", err)
		} else {
			fmt.Printf("APR: %.2f%%\n", apr*100)
		}
		fmt.Println()
	}

	fmt.Print("Enter the ID of the portfolio to edit/delete (or press Enter to return): ")
	input, _ := cli.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}

	portfolio, err := cli.portfolioService.Repo.GetByID(id)
	if err != nil {
		fmt.Printf("Error retrieving portfolio: %v\n", err)
		return
	}

	cli.managePortfolio(portfolio)
}

func (cli *CLI) managePortfolio(portfolio *models.Portfolio) {
	fmt.Printf("Portfolio: %s\n", portfolio.Name)
	fmt.Println("Select an option:")
	fmt.Println("1. Edit portfolio")
	fmt.Println("2. Delete portfolio")
	fmt.Println("3. Return")

	fmt.Print("Option: ")
	input, _ := cli.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		cli.editPortfolio(portfolio)
	case "2":
		err := cli.portfolioService.Repo.Delete(portfolio.ID)
		if err != nil {
			fmt.Printf("Error deleting portfolio: %v\n", err)
		} else {
			fmt.Println("Portfolio deleted successfully.")
		}
	case "3":
		return
	default:
		fmt.Println("Invalid option.")
	}
}

func (cli *CLI) editPortfolio(portfolio *models.Portfolio) {
	// TODO
}

func (cli *CLI) createPortfolioManual() {
	fmt.Print("Enter the name of the portfolio: ")
	name, _ := cli.reader.ReadString('\n')
	name = strings.TrimSpace(name)

	symbols, err := api.GetSP500Symbols()
	if err != nil {
		fmt.Printf("Error retrieving S&P 500 symbols: %v\n", err)
		return
	}

	var stocks []models.Stock

	for {
		fmt.Println("-Select a stock from the S&P 500 or press Enter to finish:")
		for i, symbol := range symbols {
			fmt.Printf("%d. %s\n", i+1, symbol)
		}

		fmt.Print("Stock number: ")
		input, _ := cli.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			break
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(symbols) {
			fmt.Println("Invalid input.")
			continue
		}

		symbol := symbols[index-1]

		fmt.Printf("Enter the quantity of shares for %s: ", symbol)
		qtyInput, _ := cli.reader.ReadString('\n')
		qtyInput = strings.TrimSpace(qtyInput)
		quantity, err := strconv.Atoi(qtyInput)
		if err != nil || quantity <= 0 {
			fmt.Println("Invalid quantity.")
			continue
		}

		fmt.Printf("Enter the purchase date (YYYY-MM-DD) for %s: ", symbol)
		dateInput, _ := cli.reader.ReadString('\n')
		dateInput = strings.TrimSpace(dateInput)
		buyDate, err := time.Parse("2006-01-02", dateInput)
		if err != nil {
			fmt.Println("Invalid date.")
			continue
		}

		stock := models.Stock{
			Symbol:   symbol,
			Quantity: quantity,
			BuyDate:  buyDate,
		}

		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		fmt.Println("No stocks added to the portfolio.")
		return
	}

	portfolio := &models.Portfolio{
		Name:   name,
		Stocks: stocks,
	}

	err = cli.portfolioService.Repo.Save(portfolio)
	if err != nil {
		fmt.Printf("Error saving portfolio: %v\n", err)
		return
	}

	fmt.Println("Portfolio created successfully.")
}

func (cli *CLI) createPortfolioRandom() {
	// TODO
}

func earliestBuyDate(portfolio models.Portfolio) time.Time {
	if len(portfolio.Stocks) == 0 {
		return time.Now()
	}

	earliest := portfolio.Stocks[0].BuyDate
	for _, stock := range portfolio.Stocks {
		if stock.BuyDate.Before(earliest) {
			earliest = stock.BuyDate
		}
	}
	return earliest
}
