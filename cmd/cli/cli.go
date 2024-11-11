package cli

import (
	"bufio"
	"fmt"
	"github.com/fcopulgar/stock-manager-go/services"
	"os"
	"strings"
)

type CLI struct {
	portfolioService *services.PortfolioService
	reader           *bufio.Reader
}

func NewCLI(portfolioService *services.PortfolioService) *CLI {
	return &CLI{
		portfolioService: portfolioService,
		reader:           bufio.NewReader(os.Stdin),
	}
}

func (cli *CLI) Run() {
	for {
		fmt.Println("Select an option:")
		fmt.Println("1. View portfolios")
		fmt.Println("2. Create portfolio manually")
		fmt.Println("3. Create random portfolio")
		fmt.Println("4. Exit")

		fmt.Print("Option: ")
		input, _ := cli.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			cli.viewPortfolios()
		case "2":
			cli.createPortfolioManual()
		case "3":
			cli.createPortfolioRandom()
		case "4":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid option.")
		}
	}
}
