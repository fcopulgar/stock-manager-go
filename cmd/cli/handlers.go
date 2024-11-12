// cli/handlers.go
package cli

import (
	"fmt"
	"github.com/fcopulgar/stock-manager-go/api"
	"github.com/fcopulgar/stock-manager-go/models"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (cli *CLI) viewPortfolios() {
	portfolios, err := cli.portfolioService.GetAllPortfolios()
	if err != nil {
		fmt.Fprintf(cli.writer, "Error retrieving portfolios: %v\n", err)
		return
	}

	if len(portfolios) == 0 {
		fmt.Fprintln(cli.writer, "No portfolios available.")
		return
	}

	for _, p := range portfolios {
		fmt.Fprintf(cli.writer, "ID: %d, Name: %s\n", p.ID, p.Name)
		fmt.Fprintln(cli.writer, "Stocks:")
		for _, stock := range p.Stocks {
			fmt.Fprintf(cli.writer, "- %s: %d shares bought on %s\n", stock.Symbol, stock.Quantity, stock.BuyDate.Format("2006-01-02"))

			// Obtener el precio de cierre de la acción en la fecha de compra
			price, err := cli.portfolioService.GetPriceClose(stock.Symbol, stock.BuyDate)
			if err != nil {
				fmt.Fprintf(cli.writer, "  Price on %s: Could not retrieve price (%v)\n", stock.BuyDate.Format("2006-01-02"), err)
			} else {
				fmt.Fprintf(cli.writer, "  Price on %s: $%.2f\n", stock.BuyDate.Format("2006-01-02"), price)
			}
		}

		// Calcular APR usando la fecha de compra más temprana y hoy
		startDate := earliestBuyDate(p)
		endDate := time.Now()
		apr, err := cli.portfolioService.CalculateAPR(&p, startDate, endDate)
		if err != nil {
			fmt.Fprintf(cli.writer, "APR: Could not calculate APR (%v)\n", err)
		} else {
			fmt.Fprintf(cli.writer, "APR: %.2f%%\n", apr*100)
		}
		fmt.Println()
	}

	fmt.Fprint(cli.writer, "Enter the ID of the portfolio to edit/delete (or press Enter to return): ")
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(cli.writer, "Error reading input: %v\n", err)
		return
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Fprintln(cli.writer, "Invalid ID.")
		return
	}

	portfolio, err := cli.portfolioService.GetPortfolioByID(id)
	if err != nil {
		fmt.Fprintf(cli.writer, "Error retrieving portfolio: %v\n", err)
		return
	}

	cli.managePortfolio(portfolio)
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

func (cli *CLI) managePortfolio(portfolio *models.Portfolio) {
	fmt.Fprintf(cli.writer, "Portfolio: %s\n", portfolio.Name)
	fmt.Fprintln(cli.writer, "Select an option:")
	fmt.Fprintln(cli.writer, "1. Edit portfolio")
	fmt.Fprintln(cli.writer, "2. Delete portfolio")
	fmt.Fprintln(cli.writer, "3. Return")

	fmt.Fprint(cli.writer, "Option: ")
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(cli.writer, "Error reading input: %v\n", err)
		return
	}
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		cli.editPortfolio(portfolio)
	case "2":
		err := cli.portfolioService.DeletePortfolio(portfolio.ID)
		if err != nil {
			fmt.Fprintf(cli.writer, "Error deleting portfolio: %v\n", err)
		} else {
			fmt.Fprintln(cli.writer, "Portfolio deleted successfully.")
		}
	case "3":
		return
	default:
		fmt.Fprintln(cli.writer, "Invalid option.")
	}
}

func (cli *CLI) editPortfolio(portfolio *models.Portfolio) {
	// TODO: Implementar la lógica de edición de cartera
	fmt.Fprintln(cli.writer, "Edit portfolio functionality is not implemented yet.")
}

func (cli *CLI) createPortfolioManual() {
	fmt.Fprint(cli.writer, "Enter the name of the portfolio: ")
	name, err := cli.reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(cli.writer, "Error reading portfolio name: %v\n", err)
		return
	}
	name = strings.TrimSpace(name)

	symbols, err := api.GetSP500Symbols(&api.DefaultHTTPClient{})
	if err != nil {
		fmt.Fprintf(cli.writer, "Error retrieving S&P 500 symbols: %v\n", err)
		return
	}

	var stocks []models.Stock

	for {
		fmt.Fprintln(cli.writer, "\nSelect a stock from the S&P 500 or press Enter to finish:")
		for i, symbol := range symbols {
			fmt.Fprintf(cli.writer, "%d. %s\n", i+1, symbol)
		}

		fmt.Fprint(cli.writer, "Stock number: ")
		input, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(cli.writer, "Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)

		if input == "" {
			break
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(symbols) {
			fmt.Fprintln(cli.writer, "Invalid input.")
			continue
		}

		symbol := symbols[index-1]

		fmt.Fprintf(cli.writer, "Enter the quantity of shares for %s: ", symbol)
		qtyInput, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(cli.writer, "Error reading quantity: %v\n", err)
			continue
		}
		qtyInput = strings.TrimSpace(qtyInput)
		quantity, err := strconv.Atoi(qtyInput)
		if err != nil || quantity <= 0 {
			fmt.Fprintln(cli.writer, "Invalid quantity.")
			continue
		}

		fmt.Fprintf(cli.writer, "Enter the purchase date (YYYY-MM-DD) for %s: ", symbol)
		dateInput, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(cli.writer, "Error reading date: %v\n", err)
			continue
		}
		dateInput = strings.TrimSpace(dateInput)
		buyDate, err := time.Parse("2006-01-02", dateInput)
		if err != nil {
			fmt.Fprintln(cli.writer, "Invalid date.")
			continue
		}

		buyPrice, err := cli.portfolioService.GetPriceClose(symbol, buyDate)
		if err != nil {
			fmt.Fprintf(cli.writer, "Error getting price for %s on %s: %v... continuing with the next...\n", symbol, buyDate.Format("2006-01-02"), err)
			continue
		}

		stock := models.Stock{
			Symbol:   symbol,
			Quantity: quantity,
			BuyDate:  buyDate,
			BuyPrice: buyPrice,
		}

		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		fmt.Fprintln(cli.writer, "No stocks added to the portfolio.")
		return
	}

	portfolio := &models.Portfolio{
		Name:   name,
		Stocks: stocks,
	}

	err = cli.portfolioService.CreatePortfolioManual(portfolio)
	if err != nil {
		fmt.Fprintf(cli.writer, "Error saving portfolio: %v\n", err)
		return
	}

	fmt.Fprintln(cli.writer, "Portfolio created successfully.")
}

func (cli *CLI) createPortfolioRandom() {
	symbols, err := api.GetSP500Symbols(&api.DefaultHTTPClient{})
	if err != nil {
		fmt.Fprintf(cli.writer, "Error retrieving S&P 500 symbols: %v\n", err)
		return
	}

	var stocks []models.Stock

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		index := rand.Intn(len(symbols))
		symbol := symbols[index]

		quantity := rand.Intn(100) + 1

		start := time.Now().AddDate(-3, 0, 0).Unix()
		end := time.Now().Unix()
		timestamp := rand.Int63n(end-start) + start
		buyDate := time.Unix(timestamp, 0)

		buyPrice, err := cli.portfolioService.GetPriceClose(symbol, buyDate)
		if err != nil {
			fmt.Fprintf(cli.writer, "Error getting price for %s on %s: %v... continuing with the next...\n", symbol, buyDate.Format("2006-01-02"), err)
			continue
		}

		stock := models.Stock{
			Symbol:   symbol,
			Quantity: quantity,
			BuyDate:  buyDate,
			BuyPrice: buyPrice,
		}

		stocks = append(stocks, stock)
	}

	portfolio := &models.Portfolio{
		Name:   fmt.Sprintf("Random Portfolio %d", rand.Intn(1000)),
		Stocks: stocks,
	}

	err = cli.portfolioService.CreatePortfolioManual(portfolio)
	if err != nil {
		fmt.Fprintf(cli.writer, "Error saving portfolio: %v\n", err)
		return
	}

	// Calcular APR
	startDate := earliestBuyDate(*portfolio)
	endDate := time.Now()
	apr, err := cli.portfolioService.CalculateAPR(portfolio, startDate, endDate)
	if err != nil {
		fmt.Fprintf(cli.writer, "Error calculating APR: %v\n", err)
		return
	}

	fmt.Fprintf(cli.writer, "Random portfolio created successfully. APR: %.2f%%\n", apr*100)
}
